package domain

import "errors"

type Saldo struct {
	quantidade int
}

func NewSaldo(qtd int) (Saldo, error) {
	if qtd < 0 {
		return Saldo{}, errors.New("o saldo inicial não pode ser negativo")
	}

	return Saldo{quantidade: qtd}, nil
}

func (s Saldo) Subtrair(qtd int) (Saldo, error) {
	if qtd <= 0 {
		return Saldo{}, errors.New("a quantidade a ser subtraída deve ser maior que zero")
	}

	novoValor := s.quantidade - qtd

	if novoValor < 0 {
		return Saldo{}, errors.New("saldo insuficiente para esta operação")
	}

	return Saldo{quantidade: novoValor}, nil
}

func (s Saldo) Valor() int {
	return s.quantidade
}
