package service

import "Korp_Teste_VictorPizzarro/service-faturamento/internal/domain"

type EstoqueClient interface {
	AbaterEstoqueLote(itens []domain.ItemNotaFiscal) error
}
