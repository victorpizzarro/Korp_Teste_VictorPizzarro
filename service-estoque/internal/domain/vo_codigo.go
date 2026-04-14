package domain

import (
	"errors"
	"strings"
)

type Codigo struct {
	valor string
}

func NewCodigo(v string) (Codigo, error) {
	v = strings.TrimSpace(v)

	if v == "" {
		return Codigo{}, errors.New("código não pode ser vazio")
	}
	return Codigo{valor: v}, nil
}

func (c Codigo) Valor() string {
	return c.valor
}
