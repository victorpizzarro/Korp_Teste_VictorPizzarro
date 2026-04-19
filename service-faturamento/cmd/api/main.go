package main

import (
	"log"
	"os"

	_ "Korp_Teste_VictorPizzarro/service-faturamento/docs"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/handler"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/infrastructure"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/integration"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/repository"
	"Korp_Teste_VictorPizzarro/service-faturamento/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           Service Faturamento API
// @version         1.0
// @description     Microsserviço de faturamento e emissão de notas para o teste Korp.
// @host            localhost:8082
// @BasePath        /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Lendo variáveis do sistema.")
	}

	db, err := infrastructure.ConectarBanco()
	if err != nil {
		log.Fatalf("Falha ao conectar no banco de dados: %v", err)
	}

	err = db.AutoMigrate(&repository.NotaFiscalDB{}, &repository.ItemNotaFiscalDB{})
	if err != nil {
		log.Fatalf("Falha ao rodar migrations: %v", err)
	}

	estoqueURL := os.Getenv("ESTOQUE_URL")
	if estoqueURL == "" {
		estoqueURL = "http://localhost:8081/api/v1"
	}
	estoqueClient := integration.NewEstoqueClientHttp(estoqueURL)

	repo := repository.NewNotaFiscalRepository(db)

	cadastrarService := service.NewCadastrarNotaFiscalService(repo)
	imprimirService := service.NewImprimirNotaFiscalService(repo, estoqueClient)
	analisarService := service.NewAnalisarAnomaliaNFService(repo)

	nfHandler := handler.NewNotaFiscalHandler(cadastrarService, imprimirService, analisarService, repo)

	router := gin.Default()
	ConfigurarRotas(router, nfHandler)

	portaApi := os.Getenv("API_PORT")
	if portaApi == "" {
		portaApi = "8082"
	}

	log.Printf("Serviço de Faturamento rodando na porta %s...", portaApi)
	if err := router.Run(":" + portaApi); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
