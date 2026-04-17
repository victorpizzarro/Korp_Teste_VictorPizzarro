package domain

import "errors"

var err = errors.New("não é possível alterar uma fatura fechada")

type StatusNotaFiscal string

const (
	Aberta  StatusNotaFiscal = "Aberta"
	Fechada StatusNotaFiscal = "Fechada"
)
