package service

import (
	"Korp_Teste_VictorPizzarro/service-estoque/internal/domain"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProdutoRepository struct {
	mock.Mock
}

func (mock *MockProdutoRepository) SalvarProduto(produto *domain.Produto) error {
	args := mock.Called(produto)
	return args.Error(0)
}

func (mock *MockProdutoRepository) BuscarProdutoPorCodigo(codigo string) (*domain.Produto, error) {
	args := mock.Called(codigo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Produto), args.Error(1)
}

func (mock *MockProdutoRepository) AtualizarSaldo(produto *domain.Produto) error {
	args := mock.Called(produto)
	return args.Error(0)
}

func (mock *MockProdutoRepository) DebitarSaldo(codigo string, quantidade int) error {
	args := mock.Called(codigo, quantidade)
	return args.Error(0)
}

func (mock *MockProdutoRepository) DebitarSaldoLote(itens []repository.ItemDebito) error {
	args := mock.Called(itens)
	return args.Error(0)
}

func (mock *MockProdutoRepository) ListarTodos() ([]*domain.Produto, error) {
	args := mock.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Produto), args.Error(1)
}

func TestDebitarEstoqueService_ExecutarLote(t *testing.T) {
	tests := []struct {
		name          string
		itensInput    []repository.ItemDebito
		mockReturnErr error
		expectedErr   string
	}{
		{
			name: "Sucesso - Debitar Saldo em Lote",
			itensInput: []repository.ItemDebito{
				{Codigo: "PROD-001", Quantidade: 5},
				{Codigo: "PROD-002", Quantidade: 2},
			},
			mockReturnErr: nil,
			expectedErr:   "",
		},
		{
			name: "Falha - Saldo Insuficiente em Lote",
			itensInput: []repository.ItemDebito{
				{Codigo: "PROD-001", Quantidade: 50},
			},
			mockReturnErr: errors.New("saldo insuficiente para o produto: PROD-001"),
			expectedErr:   "saldo insuficiente para o produto: PROD-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProdutoRepository)
			if tt.itensInput != nil {
				mockRepo.On("DebitarSaldoLote", tt.itensInput).Return(tt.mockReturnErr)
			}

			service := NewDebitarEstoqueService(mockRepo)
			err := service.ExecutarLote(tt.itensInput)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
