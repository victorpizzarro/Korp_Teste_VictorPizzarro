package main

import (
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigurarRotas(router *gin.Engine, nfHandler *handler.NotaFiscalHandler) {

	router.GET("/healthy", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"status":  "Faturamento Service está vivo!",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		api.POST("/notas-fiscais", nfHandler.Cadastrar)
		api.POST("/notas-fiscais/:numero/imprimir", nfHandler.Imprimir)
	}
}
