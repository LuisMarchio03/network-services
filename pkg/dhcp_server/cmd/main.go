package main

import (
	"context"
	"fmt"
	"network-services-server-dhcp/internal/dhcp"
	"network-services-server-dhcp/internal/mongodb"
	"network-services-server-dhcp/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// connectToMongoDB estabelece uma conexão com o servidor MongoDB utilizando o URI fornecido.
// Recebe um contexto (ctx) para controle de tempo e um possível erro se a conexão falhar.
// Retorna uma instância de *mongodb.MongoDB, que representa a conexão com o MongoDB, e um erro, se houver.
//
// Exemplo de uso:
//
//	client, err := connectToMongoDB()
//	if err != nil {
//	    logger.Error("Erro ao conectar ao MongoDB:", err)
//	    return nil, err
//	}
//	defer client.Close(context.Background())
//
//	- Agora você pode usar 'client' para realizar operações no servidor MongoDB, como inserção, consulta, etc.
func connectToMongoDB() (*mongodb.MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configuração do URI para conexão com o MongoDB
	uri := "mongodb://root:example@localhost:27017"
	dbName := "dhcp"
	collectionName := "ip_addresses"

	// Conectar ao servidor MongoDB utilizando a função ConnectMongoDB do pacote mongodb
	client, err := mongodb.ConnectMongoDB(ctx, uri, dbName, collectionName)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao MongoDB: %v", err)
	}

	fmt.Println("Conexão com o MongoDB estabelecida")
	return client, nil
}

// waitForShutdown aguarda por sinais de interrupção (SIGINT) ou término (SIGTERM) para encerrar o servidor DHCP de forma segura.
//
// Exemplo de uso:
//
//	waitForShutdown()
//
//	- Essa função é utilizada para manter o servidor DHCP em execução até que um sinal de interrupção seja recebido.
func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("Servidor DHCP encerrado")
}

// main é a função principal do programa.
// Inicializa o logger, conecta-se ao MongoDB, cria uma instância do servidor DHCP e inicia o servidor em uma goroutine.
// Aguarda por um sinal para encerrar o servidor DHCP de forma segura.
//
// Exemplo de uso:
//
//	main()
//
//	- Esta função é chamada para iniciar o servidor DHCP. Ela gerencia todas as etapas necessárias para o funcionamento do servidor.
func main() {
	// Inicializar o logger
	logger.InitLogger()

	// Conectar ao MongoDB
	db, err := connectToMongoDB()
	if err != nil {
		logger.Error("Erro ao conectar ao MongoDB:", err)
		return
	}
	defer db.Close(context.Background())

	// Criar uma nova instância do servidor DHCP
	server, err := dhcp.NewServer(context.Background(), db)
	if err != nil {
		logger.Error("Erro ao criar servidor DHCP:", err)
		return
	}

	// Iniciar o servidor DHCP em uma goroutine
	go func() {
		err := server.Start()
		if err != nil {
			logger.Error("Erro ao iniciar servidor DHCP:", err)
		}
	}()

	// Aguardar por um sinal para encerrar o servidor DHCP
	waitForShutdown()
}
