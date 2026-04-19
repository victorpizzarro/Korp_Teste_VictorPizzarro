package handler

import (
	"Korp_Teste_VictorPizzarro/service-estoque/internal/domain"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/repository"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProdutoRequest struct {
	Codigo    string `json:"codigo" binding:"required"`
	Descricao string `json:"descricao" binding:"required"`
	Saldo     int    `json:"saldo"`
}

type ProdutoResponse struct {
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
	Mensagem  string `json:"mensagem,omitempty"`
}

type DebitarRequest struct {
	Quantidade int `json:"quantidade" binding:"required,gt=0"`
}

type LoteDebitoRequest struct {
	Itens []struct {
		CodigoProduto string `json:"codigoProduto" binding:"required"`
		Quantidade    int    `json:"quantidade" binding:"required,gt=0"`
	} `json:"itens" binding:"required,dive"`
}

func NewProdutoResponse(produto *domain.Produto) ProdutoResponse {
	return ProdutoResponse{
		Codigo:    produto.Codigo(),
		Descricao: produto.Descricao(),
		Saldo:     produto.Saldo(),
	}
}

type ProdutoHandler struct {
	service        *service.CadastrarProdutoService
	debitarService *service.DebitarEstoqueService
	repo           repository.ProdutoRepository
}

func NewProdutoHandler(service *service.CadastrarProdutoService, debitarService *service.DebitarEstoqueService, repo repository.ProdutoRepository) *ProdutoHandler {
	return &ProdutoHandler{
		service:        service,
		debitarService: debitarService,
		repo:           repo,
	}
}

// @Summary      Listar Produtos
// @Description  Retorna todos os produtos cadastrados com seus saldos
// @Tags         produtos
// @Produce      json
// @Success      200  {array}   handler.ProdutoResponse
// @Failure      500  {object}  map[string]string
// @Router       /produtos [get]
func (h *ProdutoHandler) Listar(c *gin.Context) {
	produtos, err := h.repo.ListarTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao listar produtos"})
		return
	}

	var response []ProdutoResponse
	for _, p := range produtos {
		response = append(response, NewProdutoResponse(p))
	}

	if response == nil {
		response = []ProdutoResponse{}
	}

	c.JSON(http.StatusOK, response)
}

// @Description  Cria um produto no banco de dados com código, descrição e saldo inicial
// @Tags         produtos
// @Accept       json
// @Produce      json
// @Param        produto  body      handler.ProdutoRequest  true  "Dados do Produto"
// @Success      201      {object}  handler.ProdutoResponse
// @Success      200      {object}  handler.ProdutoResponse
// @Failure      400      {object}  map[string]string
// @Failure      422      {object}  map[string]string
// @Router       /produtos [post]
func (h *ProdutoHandler) Cadastrar(c *gin.Context) {
	var request ProdutoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Campos obrigatórios ausentes ou em formato inválido"})
		return
	}

	produto, criado, err := h.service.Cadastrar(
		request.Codigo,
		request.Descricao,
		request.Saldo,
	)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": err.Error()})
		return
	}

	response := NewProdutoResponse(produto)

	if criado {
		response.Mensagem = "Produto cadastrado com sucesso"
		c.JSON(http.StatusCreated, response)
		return
	}

	response.Mensagem = "Produto já está cadastrado"
	c.JSON(http.StatusOK, response)
}

// @Description  Debita a quantidade de um produto no estoque individualmente
// @Tags         produtos
// @Accept       json
// @Produce      json
// @Param        codigo      path      string                  true  "Código do Produto"
// @Param        quantidade  body      handler.DebitarRequest  true  "Quantidade a debitar"
// @Success      200         {object}  map[string]string
// @Failure      400         {object}  map[string]string
// @Failure      422         {object}  map[string]string
// @Router       /produtos/{codigo}/debitar [post]
func (h *ProdutoHandler) Debitar(c *gin.Context) {
	codigo := c.Param("codigo")

	var request DebitarRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Quantidade inválida"})
		return
	}

	err := h.debitarService.Executar(codigo, request.Quantidade)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Estoque debitado com sucesso"})
}

// @Summary Debitar estoque em lote
// @Tags Produtos
// @Accept json
// @Produce json
// @Param request body LoteDebitoRequest true "Itens para débito"
// @Success 204 "Débito realizado com sucesso"
// @Failure 400,422,500 {object} map[string]string
// @Router /produtos/debitar-lote [post]
func (h *ProdutoHandler) DebitarLote(c *gin.Context) {
	var req LoteDebitoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "payload inválido"})
		return
	}

	var itensParaDebito []repository.ItemDebito

	for _, item := range req.Itens {
		itensParaDebito = append(itensParaDebito, repository.ItemDebito{
			Codigo:     item.CodigoProduto,
			Quantidade: item.Quantidade,
		})
	}

	err := h.debitarService.ExecutarLote(itensParaDebito)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
