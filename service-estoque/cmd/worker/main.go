package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Korp_Teste_VictorPizzarro/service-estoque/internal/infrastructure"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/repository"
	"Korp_Teste_VictorPizzarro/service-estoque/internal/service"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Lendo variáveis do sistema.")
	}

	db, err := infrastructure.ConectarBanco()
	if err != nil {
		log.Fatalf("Falha ao conectar no banco de dados: %v", err)
	}

	repo := repository.NewProdutoRepository(db)
	debitarService := service.NewDebitarEstoqueService(repo)

	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Falha ao conectar no RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Falha ao abrir o canal: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"diminuir_estoque_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Falha ao declarar a fila: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Falha ao registrar o consumidor: %v", err)
	}

	log.Println("Worker de Estoque iniciado. Conectado ao PostgreSQL e RabbitMQ...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		for d := range msgs {
			var data struct {
				Codigo     string `json:"codigo"`
				Quantidade int    `json:"quantidade"`
			}

			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("ERRO: Falha ao decodificar JSON: %v", err)
				continue
			}

			log.Printf("Processando débito: %d unidades do produto %s", data.Quantidade, data.Codigo)

			err = debitarService.Executar(data.Codigo, data.Quantidade)

			if err != nil {
				log.Printf("ERRO ao processar produto %s: %v", data.Codigo, err)
				continue
			}

			log.Printf("SUCESSO: Saldo do produto %s atualizado no banco!", data.Codigo)
		}
	}()

	<-stop

	log.Println("\nSinal recebido. Desligando o worker de Estoque...")
}
