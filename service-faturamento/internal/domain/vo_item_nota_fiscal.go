package domain

import "errors"

type ItemNotaFiscal struct {
	codigoProduto string
	quantidade    int
}

func NewItemNotaFiscal(codigoProduto string, quantidade int) (*ItemNotaFiscal, error) {
	if codigoProduto == "" {
		return nil, errors.New("o código do produto não pode ser vazio")
	}

	if quantidade <= 0 {
		return nil, errors.New("a quantidade do item deve ser maior que zero")
	}

	return &ItemNotaFiscal{
		codigoProduto: codigoProduto,
		quantidade:    quantidade,
	}, nil
}

func (itemNotaFiscal *ItemNotaFiscal) CodigoProduto() string { return itemNotaFiscal.codigoProduto }
func (itemNotaFiscal *ItemNotaFiscal) Quantidade() int       { return itemNotaFiscal.quantidade }
