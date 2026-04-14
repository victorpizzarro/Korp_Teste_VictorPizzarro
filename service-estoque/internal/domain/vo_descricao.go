package domain

import (
	"errors"
	"strings"
)

type Descricao struct {
	valor string
}

func NewDescricao(v string) (Descricao, error) {

	v = strings.TrimSpace(v)

	if v == "" {
		return Descricao{}, errors.New("descrição não pode ser vazia")
	}
	return Descricao{valor: v}, nil
}

func (d Descricao) Valor() string {
	return d.valor
}
