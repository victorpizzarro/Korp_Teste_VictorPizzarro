package service

import (
	"Korp_Teste_VictorPizzarro/service-estoque/internal/repository"
	"errors"
)

type DebitarEstoqueService struct {
	repo repository.ProdutoRepository
}

func NewDebitarEstoqueService(repo repository.ProdutoRepository) *DebitarEstoqueService {
	return &DebitarEstoqueService{repo: repo}
}

func (service *DebitarEstoqueService) Executar(codigo string, quantidade int) error {

	produto, err := service.repo.BuscarProdutoPorCodigo(codigo)
	if err != nil {
		return err
	}

	if produto == nil {
		return errors.New("produto não encontrado")
	}

	err = produto.DiminuirEstoque(quantidade)
	if err != nil {
		return err
	}

	err = service.repo.AtualizarSaldo(produto)
	if err != nil {
		return err
	}

	return nil
}
