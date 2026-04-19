package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotaFiscal_AdicionarItem(t *testing.T) {
	tests := []struct {
		name        string
		notaAberta  bool
		expectedErr string
	}{
		{
			name:        "Sucesso - Adicionar item em Nota Fiscal Aberta",
			notaAberta:  true,
			expectedErr: "",
		},
		{
			name:        "Falha - Tentar adicionar item em Nota Fiscal Fechada",
			notaAberta:  false,
			expectedErr: "não é possível alterar uma nota fiscal fechada",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nota, _ := NewNotaFiscal(1)

			if !tt.notaAberta {
				nota.Fechar()
			}

			item, _ := NewItemNotaFiscal("PROD-001", 5)
			err := nota.AdicionarItem(*item)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
				assert.Len(t, nota.Itens(), 0)
			} else {
				assert.NoError(t, err)
				assert.Len(t, nota.Itens(), 1)
			}
		})
	}
}

func TestNotaFiscal_Fechar(t *testing.T) {
	nota, _ := NewNotaFiscal(1)

	assert.Equal(t, Aberta, nota.Status())

	nota.Fechar()

	assert.Equal(t, Fechada, nota.Status())

	nota.Fechar()
	assert.Equal(t, Fechada, nota.Status())
}

func TestReconstituirNotaFiscal(t *testing.T) {
	item, _ := NewItemNotaFiscal("PROD-001", 5)
	itens := []ItemNotaFiscal{*item}
	dataCriacao := time.Now()

	nota, err := ReconstituirNotaFiscal(10, Fechada, dataCriacao, itens)

	assert.NoError(t, err)
	assert.Equal(t, 10, nota.NumeroSequencial())
	assert.Equal(t, Fechada, nota.Status())
	assert.Equal(t, dataCriacao, nota.DataCriacao())
	assert.Len(t, nota.Itens(), 1)
}
