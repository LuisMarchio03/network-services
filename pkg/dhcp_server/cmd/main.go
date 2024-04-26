package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"network-services-server-dhcp/internal/dhcp"
	"network-services-server-dhcp/internal/mongodb"
	"network-services-server-dhcp/logger"
)

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

	// Aguardar por um sinal para encerrar o servidor
	waitForShutdown()
}

func connectToMongoDB() (*mongodb.MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := "mongodb://root:example@localhost:27017"
	dbName := "dhcp"
	collectionName := "ip_addresses"

	// Conectar ao servidor MongoDB
	client, err := mongodb.ConnectMongoDB(ctx, uri, dbName, collectionName)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao MongoDB: %v", err)
	}

	fmt.Println("Conexão com o MongoDB estabelecida")
	return client, nil
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("Servidor DHCP encerrado")
}
