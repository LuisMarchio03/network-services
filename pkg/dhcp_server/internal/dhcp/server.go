package dhcp

import (
	"context"
	"fmt"
	"net"

	"network-services-server-dhcp/internal/mongodb"
	"network-services-server-dhcp/logger"
)

// StartServer inicia o servidor DHCP
func StartServer() {
	// Inicializar o socket UDP na porta 67
	conn, err := net.ListenPacket("udp", ":67")
	if err != nil {
		logger.Error("Erro ao abrir socket UDP:", err)
		return
	}
	defer conn.Close()

	logger.Info("Servidor DHCP iniciado")

	// Buffer para receber mensagens DHCP
	buffer := make([]byte, 1024)

	// Loop infinito para lidar com as mensagens DHCP
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			logger.Error("Erro ao ler mensagem DHCP:", err)
			continue
		}
		_ = n //! REMOVER ISSO

		logger.Info(fmt.Sprintf("Mensagem de informação: %s", addr.String()))

		// Processar mensagem DHCP e enviar resposta
		// TODO: Implementar a lógica para processar e responder mensagens DHCP
	}
}

// buildDHCPOffer constrói uma resposta DHCP Offer em resposta a uma solicitação DHCP Discover.
// Esta função é responsável por identificar um endereço IP disponível para oferecer ao cliente,
// construir a estrutura do pacote DHCP Offer e enviá-lo de volta ao cliente.
// Parâmetros:
//   - request: O pacote DHCP Discover enviado pelo cliente.
//
// Retorna:
//   - Um pacote DHCP Offer pronto para ser enviado ao cliente.
//
// TODO: Implementar a lógica para construir uma resposta DHCP Offer
func buildDHCPOffer(request []byte, db *mongodb.MongoDB, ctx context.Context) []byte {
	// Implemente a lógica para construir uma resposta DHCP Offer aqui
	// O exemplo abaixo simplesmente copia o pacote de solicitação como resposta

	offer := make([]byte, len(request))
	copy(offer, request)

	// Modifica o tipo de mensagem para DHCP Offer
	offer[0] = 2

	// TODO: Identificar um endereço IP disponível e incluí-lo na oferta
	// 1. Chamar FindAvailableIP do mongoDB e pegar o IP address
	availableIP, err := db.FindAvailableIP(ctx)
	if err != nil {
		logger.Error("Erro ao encontrar um endereço IP disponível:", err)
		return nil
	}

	// 2. Incluir o endereço IP na oferta DHCP
	copy(offer[16:20], availableIP) // Os bytes 16-19 contêm o endereço IP oferecido

	// TODO: Enviar a resposta DHCP Offer para o endereço IP do cliente
	// 1. Implementar a lógica de envio de resposta para o cliente aqui

	return offer
}

// buildDHCPAcknowledge constrói uma resposta DHCP Acknowledge
func buildDHCPAcknowledge(request []byte) []byte {
	// TODO: Implementar a lógica para construir uma resposta DHCP Acknowledge
	// O exemplo abaixo simplesmente copia o pacote de solicitação como resposta
	acknowledge := make([]byte, len(request))
	copy(acknowledge, request)
	acknowledge[0] = 5 // Modifica o tipo de mensagem para DHCP Acknowledge
	return acknowledge
}
