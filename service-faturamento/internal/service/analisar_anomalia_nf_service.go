package service

import (
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/repository"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type AnalisarAnomaliaNFService struct {
	repo repository.NotaFiscalRepository
}

func NewAnalisarAnomaliaNFService(repo repository.NotaFiscalRepository) *AnalisarAnomaliaNFService {
	return &AnalisarAnomaliaNFService{repo: repo}
}

type AnomaliaResponse struct {
	TemAnomalia bool   `json:"tem_anomalia"`
	Mensagem    string `json:"mensagem"`
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (s *AnalisarAnomaliaNFService) Executar(numero int) (*AnomaliaResponse, error) {
	nota, err := s.repo.BuscarNotaFiscalPorNumero(numero)
	if err != nil {
		return nil, err
	}
	if nota == nil {
		return nil, errors.New("nota fiscal não encontrada")
	}

	prompt := "Analise esta lista de itens de uma nota fiscal. Existe alguma quantidade absurdamente alta que pareça um erro humano de digitação para um varejo/atacado padrão (ex: acima de 50 ou 100 dependendo)? Responda apenas com um JSON obrigatoriamente sem formatacao markdown, no formato estrito: {\"tem_anomalia\": bool, \"mensagem\": \"justificativa curta\"}.\n\nItens:\n"
	for _, item := range nota.Itens() {
		prompt += fmt.Sprintf("- Produto %s: %d unidades\n", item.CodigoProduto(), item.Quantidade())
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return &AnomaliaResponse{TemAnomalia: false, Mensagem: "Chave do Gemini ausente."}, nil
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + apiKey

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini retornou status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gResp geminiResponse
	if err := json.Unmarshal(bodyBytes, &gResp); err != nil {
		return nil, fmt.Errorf("erro parser gemini: %w", err)
	}

	if len(gResp.Candidates) == 0 || len(gResp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("resposta vazia do gemini")
	}

	respText := gResp.Candidates[0].Content.Parts[0].Text
	respText = strings.TrimSpace(respText)
	respText = strings.ReplaceAll(respText, "```json", "")
	respText = strings.ReplaceAll(respText, "```", "")

	var result AnomaliaResponse
	if err := json.Unmarshal([]byte(respText), &result); err != nil {
		return nil, fmt.Errorf("erro lendo JSON final do gemini: %w, payload recebido: %s", err, respText)
	}

	return &result, nil
}
