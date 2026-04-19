package service

import (
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/domain"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotaFiscalRepository struct {
	mock.Mock
}

func (m *MockNotaFiscalRepository) SalvarNotaFiscal(nota *domain.NotaFiscal) error {
	args := m.Called(nota)
	return args.Error(0)
}

func (m *MockNotaFiscalRepository) BuscarNotaFiscalPorNumero(numero int) (*domain.NotaFiscal, error) {
	args := m.Called(numero)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotaFiscal), args.Error(1)
}

func (m *MockNotaFiscalRepository) AtualizarStatus(nota *domain.NotaFiscal) error {
	args := m.Called(nota)
	return args.Error(0)
}

func (m *MockNotaFiscalRepository) ListarTodas() ([]*domain.NotaFiscal, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.NotaFiscal), args.Error(1)
}

type MockEstoqueClient struct {
	mock.Mock
}

func (m *MockEstoqueClient) AbaterEstoqueLote(itens []domain.ItemNotaFiscal) error {
	args := m.Called(itens)
	return args.Error(0)
}

func TestImprimirNotaFiscalService_Executar(t *testing.T) {
	tests := []struct {
		name                 string
		numeroNota           int
		setupNota            func() *domain.NotaFiscal
		setupMockRepo        func(repo *MockNotaFiscalRepository, nota *domain.NotaFiscal)
		setupMockEstoque     func(estoque *MockEstoqueClient, nota *domain.NotaFiscal)
		expectedErrSubstring string
		expectedStatus       domain.StatusNotaFiscal
	}{
		{
			name:       "Sucesso - Emissão da NF (Happy Path)",
			numeroNota: 1,
			setupNota: func() *domain.NotaFiscal {
				nota, _ := domain.NewNotaFiscal(1)
				item, _ := domain.NewItemNotaFiscal("PROD-001", 5)
				nota.AdicionarItem(*item)
				return nota
			},
			setupMockRepo: func(repo *MockNotaFiscalRepository, nota *domain.NotaFiscal) {
				repo.On("BuscarNotaFiscalPorNumero", 1).Return(nota, nil)
				repo.On("AtualizarStatus", mock.AnythingOfType("*domain.NotaFiscal")).Return(nil)
			},
			setupMockEstoque: func(estoque *MockEstoqueClient, nota *domain.NotaFiscal) {
				estoque.On("AbaterEstoqueLote", nota.Itens()).Return(nil)
			},
			expectedErrSubstring: "",
			expectedStatus:       domain.Fechada,
		},
		{
			name:       "Falha - Tentar fechar NF já fechada",
			numeroNota: 2,
			setupNota: func() *domain.NotaFiscal {
				item, _ := domain.NewItemNotaFiscal("PROD-001", 5)
				nota, _ := domain.ReconstituirNotaFiscal(2, domain.Fechada, time.Now(), []domain.ItemNotaFiscal{*item})
				return nota
			},
			setupMockRepo: func(repo *MockNotaFiscalRepository, nota *domain.NotaFiscal) {
				repo.On("BuscarNotaFiscalPorNumero", 2).Return(nota, nil)
			},
			setupMockEstoque: func(estoque *MockEstoqueClient, nota *domain.NotaFiscal) {

			},
			expectedErrSubstring: "esta nota fiscal já foi impressa/fechada",
			expectedStatus:       domain.Fechada,
		},
		{
			name:       "Falha - Estoque recusa o Batch (Falha na requisição HTTTP)",
			numeroNota: 3,
			setupNota: func() *domain.NotaFiscal {
				nota, _ := domain.NewNotaFiscal(3)
				item, _ := domain.NewItemNotaFiscal("PROD-999", 50)
				nota.AdicionarItem(*item)
				return nota
			},
			setupMockRepo: func(repo *MockNotaFiscalRepository, nota *domain.NotaFiscal) {
				repo.On("BuscarNotaFiscalPorNumero", 3).Return(nota, nil)

			},
			setupMockEstoque: func(estoque *MockEstoqueClient, nota *domain.NotaFiscal) {
				estoque.On("AbaterEstoqueLote", nota.Itens()).Return(errors.New("sistema indisponível"))
			},
			expectedErrSubstring: "falha ao integrar com estoque: sistema indisponível",
			expectedStatus:       domain.Aberta,
		},
		{
			name:       "Falha - Nota fiscal não encontrada",
			numeroNota: 99,
			setupNota: func() *domain.NotaFiscal {
				return nil
			},
			setupMockRepo: func(repo *MockNotaFiscalRepository, nota *domain.NotaFiscal) {
				repo.On("BuscarNotaFiscalPorNumero", 99).Return(nil, nil)
			},
			setupMockEstoque: func(estoque *MockEstoqueClient, nota *domain.NotaFiscal) {
			},
			expectedErrSubstring: "nota fiscal não encontrada",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockNotaFiscalRepository)
			mockEstoque := new(MockEstoqueClient)

			notaMockada := tt.setupNota()
			tt.setupMockRepo(mockRepo, notaMockada)
			tt.setupMockEstoque(mockEstoque, notaMockada)

			service := NewImprimirNotaFiscalService(mockRepo, mockEstoque)
			notaResult, err := service.Executar(tt.numeroNota)

			if tt.expectedErrSubstring != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrSubstring)
				if notaMockada != nil {
					assert.Equal(t, tt.expectedStatus, notaMockada.Status())
				}
				assert.Nil(t, notaResult)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, notaResult)
				assert.Equal(t, tt.expectedStatus, notaResult.Status())
			}

			mockRepo.AssertExpectations(t)
			mockEstoque.AssertExpectations(t)
		})
	}
}
