package service

import (
	"Korp_Teste_VictorPizzarro/service-estoque/internal/repository"
)

type DebitarEstoqueService struct {
	repo repository.ProdutoRepository
}

func NewDebitarEstoqueService(repo repository.ProdutoRepository) *DebitarEstoqueService {
	return &DebitarEstoqueService{repo: repo}
}

func (service *DebitarEstoqueService) Executar(codigo string, quantidade int) error {
	return service.repo.DebitarSaldo(codigo, quantidade)
}

func (service *DebitarEstoqueService) ExecutarLote(itens []repository.ItemDebito) error {
	return service.repo.DebitarSaldoLote(itens)
}
