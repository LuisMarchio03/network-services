// main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/olivere/elastic/v7" // Biblioteca para trabalhar com Elasticsearch
	"github.com/streadway/amqp"     // Biblioteca para trabalhar com RabbitMQ
)

// LogEntry representa uma entrada de log com informações relevantes.
type LogEntry struct {
	Timestamp        string                 `json:"timestamp"`          // Hora em que o log foi gerado
	Service          string                 `json:"service"`            // Nome do serviço que gerou o log
	Level            string                 `json:"level"`              // Nível do log (ex: INFO, ERROR)
	Message          string                 `json:"message"`            // Mensagem do log
	Context          map[string]interface{} `json:"context"`            // Informações adicionais no contexto do log
	ComputerClientIP string                 `json:"computer_client_ip"` // IP do cliente que gerou o log
}

func main() {
	// Conectar ao RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Conexão com o servidor RabbitMQ
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err) // Log de erro caso a conexão falhe
	}
	defer conn.Close() // Fecha a conexão quando o main finalizar

	// Cria um canal para comunicação com o RabbitMQ
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err) // Log de erro caso a abertura do canal falhe
	}
	defer ch.Close() // Fecha o canal quando o main finalizar

	// Consumir a fila de logs
	msgs, err := ch.Consume(
		"logs_queue", // Nome da fila
		"",           // Consumer (vazio para usar o padrão)
		true,         // Auto-ack: true se a mensagem deve ser confirmada automaticamente
		false,        // Exclusiva: true se o canal deve ser exclusivo
		false,        // No-wait: true se não quer esperar pela confirmação do servidor
		false,        // Args: argumentos adicionais (neste caso, nenhum)
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err) // Log de erro caso o consumidor falhe
	}

	// Conectar ao Elasticsearch
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false)) // Criação do cliente Elasticsearch
	if err != nil {
		log.Fatalf("Error creating the client: %s", err) // Log de erro caso a criação do cliente falhe
	}

	// Processar mensagens da fila
	forever := make(chan bool) // Canal para manter o programa em execução

	go func() {
		for d := range msgs { // Loop para processar mensagens recebidas
			logEntry := LogEntry{} // Inicializa uma nova entrada de log
			if err := json.Unmarshal(d.Body, &logEntry); err != nil {
				log.Printf("Error unmarshaling log entry: %s", err) // Log de erro caso a deserialização falhe
				continue                                            // Continua para a próxima mensagem em caso de erro
			}

			// Enviar o log para o Elasticsearch
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Define um contexto com timeout
			defer cancel()                                                          // Cancela o contexto ao finalizar

			_, err := client.Index().
				Index("logs").      // Nome do índice onde os logs serão armazenados
				BodyJson(logEntry). // Corpo da requisição com a entrada de log
				Do(ctx)             // Executa a operação
			if err != nil {
				log.Printf("Error sending log to Elasticsearch: %s", err) // Log de erro caso o envio falhe
				continue                                                  // Continua para a próxima mensagem em caso de erro
			}

			log.Printf("Log sent to Elasticsearch: %v", logEntry) // Log de sucesso
		}
	}()

	log.Println("Waiting for logs. To exit press CTRL+C") // Mensagem indicando que o programa está aguardando logs
	<-forever                                             // Mantém o programa em execução até receber um sinal de interrupção
}
