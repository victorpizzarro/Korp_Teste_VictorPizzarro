package domain

import (
	"errors"
	"time"
)

type NotaFiscal struct {
	numeroSequencial NumeroSequencial
	status           StatusNotaFiscal
	itens            []ItemNotaFiscal
	dataCriacao      time.Time
}

func NewNotaFiscal(numero int) (*NotaFiscal, error) {

	numeroVO, err := NewNumeroSequencial(numero)
	if err != nil {
		return nil, err
	}

	return &NotaFiscal{
		numeroSequencial: numeroVO,
		status:           Aberta,
		itens:            make([]ItemNotaFiscal, 0),
		dataCriacao:      time.Now(),
	}, nil
}

func (notaFiscal *NotaFiscal) AdicionarItem(item ItemNotaFiscal) error {
	if notaFiscal.status == Fechada {
		return errors.New("não é possível alterar uma nota fiscal fechada")
	}

	notaFiscal.itens = append(notaFiscal.itens, item)
	return nil
}

func (notaFiscal *NotaFiscal) Fechar() {
	if notaFiscal.status == Fechada {
		return
	}
	notaFiscal.status = Fechada
}

func ReconstituirNotaFiscal(numero int, status StatusNotaFiscal, dataCriacao time.Time, itens []ItemNotaFiscal) (*NotaFiscal, error) {
	numeroVO, err := NewNumeroSequencial(numero)
	if err != nil {
		return nil, err
	}

	return &NotaFiscal{
		numeroSequencial: numeroVO,
		status:           status,
		itens:            itens,
		dataCriacao:      dataCriacao,
	}, nil
}

func (notaFiscal *NotaFiscal) NumeroSequencial() int    { return notaFiscal.numeroSequencial.Valor() }
func (notaFiscal *NotaFiscal) Status() StatusNotaFiscal { return notaFiscal.status }
func (notaFiscal *NotaFiscal) Itens() []ItemNotaFiscal  { return notaFiscal.itens }
func (notaFiscal *NotaFiscal) DataCriacao() time.Time   { return notaFiscal.dataCriacao }
