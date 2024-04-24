// Desenvolvimento DHCP server

// Passo 1: Entendimento do Protocolo DHCP

// - Estude o protocolo DHCP para entender como ele funciona, incluindo os diferentes tipos de mensagens DHCP e seu formato.
// - Familiarize-se com as etapas do processo DHCP, como Discover, Offer, Request e Acknowledge.

// Passo 2: Configuração do Ambiente de Desenvolvimento

// - Configure seu ambiente de desenvolvimento Go, garantindo que você tenha todas as ferramentas necessárias instaladas, como Go Compiler.
// - Crie um novo projeto Go para o seu servidor DHCP.

// Passo 3: Implementação do Socket UDP

// - Crie um socket UDP para escutar na porta 67, que é a porta padrão do servidor DHCP.
// - Implemente a lógica para lidar com mensagens DHCP recebidas e enviar respostas DHCP adequadas.

// Passo 4: Manipulação das Mensagens DHCP

// - Implemente a lógica para interpretar as diferentes mensagens DHCP, como Discover, Request, etc.
// - Construa respostas DHCP apropriadas com base nas mensagens recebidas.

// Passo 5: Lógica de Atribuição de Endereços IP

// - Crie uma estrutura de dados para gerenciar os endereços IP disponíveis e atribuídos.
// - Implemente a lógica para atribuir endereços IP aos clientes DHCP conforme necessário.

// Passo 6: Testes e Depuração

// - Teste seu servidor DHCP em um ambiente de teste.
// - Depure quaisquer problemas encontrados durante os testes e ajuste sua implementação conforme necessário.

// Passo 7: Melhorias e Otimizações

// - Considere adicionar recursos adicionais ao seu servidor DHCP, como suporte a pools de endereços IP, configurações adicionais de rede, etc.
// - Otimizar o desempenho e a eficiência do seu servidor DHCP conforme necessário.

package main

import (
	"context"
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

	// Iniciar o servidor DHCP
	go dhcp.StartServer()

	// Aguardar por um sinal para encerrar o servidor
	waitForShutdown()
}

func connectToMongoDB() (*mongodb.MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := "mongodb://localhost:27017"
	dbName := "dhcp"
	collectionName := "ip_addresses"

	db, err := mongodb.ConnectMongoDB(ctx, uri, dbName, collectionName)
	if err != nil {
		return nil, err
	}

	logger.Info("Conexão com o MongoDB estabelecida")
	return db, nil
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("Servidor DHCP encerrado")
}
