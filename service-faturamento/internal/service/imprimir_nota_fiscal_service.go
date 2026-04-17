package service

import (
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/domain"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/repository"
	"errors"
	"fmt"
)

type ImprimirNotaFiscalService struct {
	repo          repository.NotaFiscalRepository
	estoqueClient EstoqueClient
}

func NewImprimirNotaFiscalService(repo repository.NotaFiscalRepository, estoque EstoqueClient) *ImprimirNotaFiscalService {
	return &ImprimirNotaFiscalService{
		repo:          repo,
		estoqueClient: estoque,
	}
}

func (service *ImprimirNotaFiscalService) Executar(numero int) (*domain.NotaFiscal, error) {
	nota, err := service.repo.BuscarNotaFiscalPorNumero(numero)
	if err != nil {
		return nil, err
	}
	if nota == nil {
		return nil, errors.New("nota fiscal não encontrada")
	}

	if nota.Status() == domain.Fechada {
		return nil, errors.New("esta nota fiscal já foi impressa/fechada")
	}

	err = service.estoqueClient.AbaterEstoqueLote(nota.Itens())
	if err != nil {
		return nil, fmt.Errorf("falha ao integrar com estoque: %w", err)
	}

	nota.Fechar()

	err = service.repo.AtualizarStatus(nota)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar status da nota: %w", err)
	}

	return nota, nil
}
