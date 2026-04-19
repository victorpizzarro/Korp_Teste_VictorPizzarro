package main

import (
	"log"
	"os"

	_ "Korp_Teste_VictorPizzarro/service-estoque/docs"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/handler"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/infrastructure"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/repository"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           Service Estoque API
// @version         1.0
// @description     Microsserviço de controle de estoque para o teste Korp.
// @host            localhost:8081
// @BasePath        /api/v1
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Lendo variáveis do sistema.")
	}

	db, err := infrastructure.ConectarBanco()
	if err != nil {
		log.Fatalf("Falha ao conectar no banco de dados: %v", err)
	}

	err = db.AutoMigrate(&repository.ProdutoDB{})
	if err != nil {
		log.Fatalf("Falha ao rodar migrations: %v", err)
	}

	repo := repository.NewProdutoRepository(db)

	cadastrarService := service.NewCadastrarProdutoService(repo)
	debitarService := service.NewDebitarEstoqueService(repo)

	produtoHandler := handler.NewProdutoHandler(cadastrarService, debitarService, repo)

	router := gin.Default()
	ConfigurarRotas(router, produtoHandler)

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8081"
	}

	log.Printf("Serviço de Estoque (API) rodando na porta %s...", apiPort)
	if err := router.Run(":" + apiPort); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
