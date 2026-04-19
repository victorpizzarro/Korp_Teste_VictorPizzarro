package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduto_DiminuirEstoque(t *testing.T) {
	tests := []struct {
		name          string
		saldoInicial  int
		qtdASubtrair  int
		expectedSaldo int
		expectedErr   string
	}{
		{
			name:          "Sucesso - Produto diminui estoque corretamente",
			saldoInicial:  20,
			qtdASubtrair:  5,
			expectedSaldo: 15,
			expectedErr:   "",
		},
		{
			name:          "Falha - Saldo do Produto insuficiente",
			saldoInicial:  5,
			qtdASubtrair:  10,
			expectedSaldo: 5,
			expectedErr:   "saldo insuficiente para esta operação",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			produto, _ := NewProduto("PROD-001", "Notebook", tt.saldoInicial)

			err := produto.DiminuirEstoque(tt.qtdASubtrair)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
				assert.Equal(t, tt.expectedSaldo, produto.Saldo()) // Saldo não deve alterar
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSaldo, produto.Saldo())
			}
		})
	}
}
