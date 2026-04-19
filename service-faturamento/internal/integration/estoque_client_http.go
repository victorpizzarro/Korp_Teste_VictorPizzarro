package integration

import (
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/domain"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/service"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type estoqueClientHttp struct {
	baseURL    string
	httpClient *http.Client
}

func NewEstoqueClientHttp(baseURL string) service.EstoqueClient {
	return &estoqueClientHttp{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type ItemDebitoPayload struct {
	CodigoProduto string `json:"codigoProduto"`
	Quantidade    int    `json:"quantidade"`
}

type LoteDebitoPayload struct {
	Itens []ItemDebitoPayload `json:"itens"`
}

func (client *estoqueClientHttp) AbaterEstoqueLote(itens []domain.ItemNotaFiscal) error {
	var payload LoteDebitoPayload

	for _, item := range itens {
		payload.Itens = append(payload.Itens, ItemDebitoPayload{
			CodigoProduto: item.CodigoProduto(),
			Quantidade:    item.Quantidade(),
		})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/produtos/debitar-lote", client.baseURL)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return errors.New("falha de comunicação com o serviço de estoque: sistema indisponível")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		var errResp struct {
			Erro string `json:"erro"`
		}
		if err := json.NewDecoder(response.Body).Decode(&errResp); err == nil && errResp.Erro != "" {
			return errors.New(errResp.Erro)
		}
		return fmt.Errorf("erro no serviço de estoque (Status: %d)", response.StatusCode)
	}

	return nil
}
