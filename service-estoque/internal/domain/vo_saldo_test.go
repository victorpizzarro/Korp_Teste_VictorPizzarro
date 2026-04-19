package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSaldo(t *testing.T) {
	tests := []struct {
		name        string
		quantidade  int
		expectedErr string
	}{
		{
			name:        "Sucesso - Saldo zerado",
			quantidade:  0,
			expectedErr: "",
		},
		{
			name:        "Sucesso - Saldo positivo",
			quantidade:  10,
			expectedErr: "",
		},
		{
			name:        "Falha - Evitar saldo inicial negativo (Proteção no VO)",
			quantidade:  -5,
			expectedErr: "o saldo inicial não pode ser negativo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saldo, err := NewSaldo(tt.quantidade)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
				assert.Equal(t, 0, saldo.Valor()) // zero value behavior
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.quantidade, saldo.Valor())
			}
		})
	}
}

func TestSaldo_Subtrair(t *testing.T) {
	tests := []struct {
		name                 string
		saldoInicial         int
		quantidadeSubtraida  int
		expectedSaldo        int
		expectedErr          string
	}{
		{
			name:                "Sucesso - Subtração válida",
			saldoInicial:        10,
			quantidadeSubtraida: 5,
			expectedSaldo:       5,
			expectedErr:         "",
		},
		{
			name:                "Falha - Saldo insuficiente",
			saldoInicial:        5,
			quantidadeSubtraida: 10,
			expectedSaldo:       0,
			expectedErr:         "saldo insuficiente para esta operação",
		},
		{
			name:                "Falha - Subtrair quantidade negativa",
			saldoInicial:        10,
			quantidadeSubtraida: -2,
			expectedSaldo:       0,
			expectedErr:         "a quantidade a ser subtraída deve ser maior que zero",
		},
		{
			name:                "Falha - Subtrair quantidade zerada",
			saldoInicial:        10,
			quantidadeSubtraida: 0,
			expectedSaldo:       0,
			expectedErr:         "a quantidade a ser subtraída deve ser maior que zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saldo, _ := NewSaldo(tt.saldoInicial)
			novoSaldo, err := saldo.Subtrair(tt.quantidadeSubtraida)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSaldo, novoSaldo.Valor())
			}
		})
	}
}
