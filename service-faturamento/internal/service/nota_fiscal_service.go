package service

import (
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/domain"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/repository"
)

type CriarNotaFiscalDTO struct {
	NumeroSequencial int
	Itens            []ItemDTO
}

type ItemDTO struct {
	CodigoProduto string
	Quantidade    int
}

type CadastrarNotaFiscalService struct {
	repo repository.NotaFiscalRepository
}

func NewCadastrarNotaFiscalService(repo repository.NotaFiscalRepository) *CadastrarNotaFiscalService {
	return &CadastrarNotaFiscalService{repo: repo}
}

func (service *CadastrarNotaFiscalService) Cadastrar(input CriarNotaFiscalDTO) (*domain.NotaFiscal, error) {

	nota, err := domain.NewNotaFiscal(input.NumeroSequencial)
	if err != nil {
		return nil, err
	}

	for _, itemInput := range input.Itens {
		itemVO, err := domain.NewItemNotaFiscal(itemInput.CodigoProduto, itemInput.Quantidade)
		if err != nil {
			return nil, err
		}

		err = nota.AdicionarItem(*itemVO)
		if err != nil {
			return nil, err
		}
	}

	err = service.repo.SalvarNotaFiscal(nota)
	if err != nil {
		return nil, err
	}

	return nota, nil
}
