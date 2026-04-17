package domain

import "errors"

type NumeroSequencial struct {
	valor int
}

func NewNumeroSequencial(valor int) (NumeroSequencial, error) {
	if valor < 0 {
		return NumeroSequencial{}, errors.New("o número sequencial não pode ser negativo")
	}

	return NumeroSequencial{
		valor: valor,
	}, nil
}

func (numeroSequencial NumeroSequencial) Valor() int {
	return numeroSequencial.valor
}

func (notaFiscal *NotaFiscal) DefinirNumeroGerado(numero int) error {
	nSequencial, err := NewNumeroSequencial(numero)
	if err != nil {
		return err
	}
	notaFiscal.numeroSequencial = nSequencial
	return nil
}
