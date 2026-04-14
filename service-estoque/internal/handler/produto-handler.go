package handler

import (
	"Korp_Teste_VictorPizzarro/service-estoque/internal/domain"
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

func NewProdutoResponse(produto *domain.Produto) ProdutoResponse {
	return ProdutoResponse{
		Codigo:    produto.Codigo(),
		Descricao: produto.Descricao(),
		Saldo:     produto.Saldo(),
	}
}

type ProdutoHandler struct {
	service *service.CadastrarProdutoService
}

func NewProdutoHandler(service *service.CadastrarProdutoService) *ProdutoHandler {
	return &ProdutoHandler{service: service}
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
