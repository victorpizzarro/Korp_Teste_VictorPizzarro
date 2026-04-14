package service

import (
	"Korp_Teste_VictorPizzarro/service-estoque/internal/domain"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/repository"
)

type CadastrarProdutoService struct {
	repo repository.ProdutoRepository
}

func NewCadastrarProdutoService(repo repository.ProdutoRepository) *CadastrarProdutoService {
	return &CadastrarProdutoService{repo: repo}
}

func (service *CadastrarProdutoService) Cadastrar(codigo, descricao string, saldo int) (*domain.Produto, bool, error) {

	produtoExistente, err := service.repo.BuscarProdutoPorCodigo(codigo)
	if err != nil {
		return nil, false, err
	}

	if produtoExistente != nil {
		return produtoExistente, false, nil
	}

	novoProduto, err := domain.NewProduto(codigo, descricao, saldo)
	if err != nil {
		return nil, false, err
	}

	err = service.repo.SalvarProduto(novoProduto)
	if err != nil {
		return nil, false, err
	}

	return novoProduto, true, nil
}
