package handler

import (
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/domain"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/repository"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ItemDTO struct {
	CodigoProduto string `json:"codigoProduto" binding:"required"`
	Quantidade    int    `json:"quantidade" binding:"required,gt=0"`
}

type NotaFiscalDTO struct {
	Itens []ItemDTO `json:"itens" binding:"required,dive"`
}

type ItemResponse struct {
	CodigoProduto string `json:"codigoProduto"`
	Quantidade    int    `json:"quantidade"`
}

type NotaFiscalResponse struct {
	NumeroSequencial int            `json:"numeroSequencial"`
	Status           string         `json:"status"`
	DataCriacao      string         `json:"dataCriacao"`
	Itens            []ItemResponse `json:"itens"`
	Mensagem         string         `json:"mensagem,omitempty"`
}

func NewNotaFiscalResponse(nota *domain.NotaFiscal) NotaFiscalResponse {
	var itensResp []ItemResponse
	for _, item := range nota.Itens() {
		itensResp = append(itensResp, ItemResponse{
			CodigoProduto: item.CodigoProduto(),
			Quantidade:    item.Quantidade(),
		})
	}

	return NotaFiscalResponse{
		NumeroSequencial: nota.NumeroSequencial(),
		Status:           string(nota.Status()),
		DataCriacao:      nota.DataCriacao().Format("2006-01-02T15:04:05Z07:00"),
		Itens:            itensResp,
	}
}

type NotaFiscalHandler struct {
	cadastrarNotaFiscalService *service.CadastrarNotaFiscalService
	imprimirNotaFiscalService  *service.ImprimirNotaFiscalService
	repo                       repository.NotaFiscalRepository
}

func NewNotaFiscalHandler(cadastrar *service.CadastrarNotaFiscalService, imprimir *service.ImprimirNotaFiscalService, repo repository.NotaFiscalRepository) *NotaFiscalHandler {
	return &NotaFiscalHandler{
		cadastrarNotaFiscalService: cadastrar,
		imprimirNotaFiscalService:  imprimir,
		repo:                       repo,
	}
}

// @Summary      Listar Notas Fiscais
// @Description  Retorna todas as notas fiscais cadastradas com seus itens
// @Tags         notas-fiscais
// @Produce      json
// @Success      200  {array}   handler.NotaFiscalResponse
// @Failure      500  {object}  map[string]string
// @Router       /notas-fiscais [get]
func (handler *NotaFiscalHandler) Listar(context *gin.Context) {
	notas, err := handler.repo.ListarTodas()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao listar notas fiscais"})
		return
	}

	var response []NotaFiscalResponse
	for _, nota := range notas {
		response = append(response, NewNotaFiscalResponse(nota))
	}

	if response == nil {
		response = []NotaFiscalResponse{}
	}

	context.JSON(http.StatusOK, response)
}

// @Summary      Cadastrar Nota Fiscal
// @Description  Cria uma nota fiscal no banco de dados com numeração sequencial e itens
// @Tags         notas-fiscais
// @Accept       json
// @Produce      json
// @Param        nota  body      handler.NotaFiscalDTO  true  "Dados da Nota Fiscal"
// @Success      201   {object}  handler.NotaFiscalResponse
// @Failure      400   {object}  map[string]string
// @Failure      422   {object}  map[string]string
// @Router       /notas-fiscais [post]
func (handler *NotaFiscalHandler) Cadastrar(context *gin.Context) {
	var request NotaFiscalDTO
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	input := service.CriarNotaFiscalDTO{
		Itens: make([]service.ItemDTO, 0),
	}
	for _, item := range request.Itens {
		input.Itens = append(input.Itens, service.ItemDTO(item))
	}

	nota, err := handler.cadastrarNotaFiscalService.Cadastrar(input)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"erro": err.Error()})
		return
	}

	response := NewNotaFiscalResponse(nota)
	response.Mensagem = "Nota Fiscal criada com sucesso"
	context.JSON(http.StatusCreated, response)
}

// @Summary      Imprimir Nota Fiscal
// @Description  Fecha a nota fiscal e debita o saldo no serviço de estoque
// @Tags         notas-fiscais
// @Accept       json
// @Produce      json
// @Param        numero  path      int  true  "Número da Nota Fiscal"
// @Success      200     {object}  handler.NotaFiscalResponse
// @Failure      400     {object}  map[string]string
// @Failure      422     {object}  map[string]string
// @Router       /notas-fiscais/{numero}/imprimir [post]
func (handler *NotaFiscalHandler) Imprimir(context *gin.Context) {
	numeroStr := context.Param("numero")
	numero, err := strconv.Atoi(numeroStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"erro": "Número inválido"})
		return
	}

	nota, err := handler.imprimirNotaFiscalService.Executar(numero)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"erro": err.Error()})
		return
	}

	response := NewNotaFiscalResponse(nota)
	response.Mensagem = "Nota Fiscal impressa e estoque atualizado"
	context.JSON(http.StatusOK, response)
}
