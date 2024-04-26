package dhcp

import (
	"context"
	"fmt"
	"net"

	"network-services-server-dhcp/internal/mongodb"
	"network-services-server-dhcp/logger"
)

// Server é uma estrutura que representa um servidor DHCP.
type Server struct {
	conn net.PacketConn
	db   *mongodb.MongoDB
	ctx  context.Context
}

// NewServer cria e configura um novo servidor DHCP.
func NewServer(ctx context.Context, db *mongodb.MongoDB) (*Server, error) {
	// Inicializar o socket UDP na porta 67
	conn, err := net.ListenPacket("udp", ":67")
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir socket UDP: %v", err)
	}

	// Criar uma nova instância do servidor DHCP
	server := &Server{
		conn: conn,
		db:   db,
		ctx:  ctx,
	}

	logger.Info("Servidor DHCP iniciado")

	return server, nil
}

// Start inicia o servidor DHCP para processar mensagens DHCP.
func (s *Server) Start() error {
	defer s.conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, addr, err := s.conn.ReadFrom(buffer)
		if err != nil {
			logger.Error("Erro ao ler mensagem DHCP:", err)
			continue
		}

		logger.Info(fmt.Sprintf("Mensagem recebida de: %s", addr.String()))

		// Processar a mensagem DHCP recebida
		go s.processDHCPMessage(buffer[:n], addr)
	}
}

// processDHCPMessage processa uma mensagem DHCP recebida e encaminha para a lógica de construção de resposta apropriada.
// Parâmetros:
//   - request: Um slice de bytes contendo a mensagem DHCP recebida.
//   - conn: Uma conexão de pacote net.PacketConn para enviar respostas DHCP.
//   - addr: O endereço net.Addr do cliente que enviou a mensagem DHCP.
//   - db: Uma instância do MongoDB para interagir com o banco de dados.
//   - ctx: O contexto da execução para controlar o tempo e o cancelamento de operações.
//
// A função processDHCPMessage determina o tipo de mensagem DHCP com base no primeiro byte da mensagem (messageType) e executa a lógica correspondente:
//   - Para DHCP Discover (messageType = 1): Chama buildDHCPOffer para construir uma oferta DHCP e a envia de volta ao cliente DHCP.
//   - Outros tipos de mensagens DHCP podem ser adicionados com casos adicionais no switch.
//
// processDHCPMessage processa uma mensagem DHCP recebida.
func (s *Server) processDHCPMessage(request []byte, addr net.Addr) {
	// Implemente a lógica para interpretar e responder à mensagem DHCP recebida
	// Você pode chamar funções auxiliares dentro do contexto do servidor DHCP (s)
	messageType := request[0]

	// Avaliar o tipo de mensagem DHCP e executar a lógica apropriada
	switch messageType {
	case 1: // DHCP Discover
		// Construir uma oferta DHCP em resposta ao DHCP Discover
		offer := s.buildDHCPOffer(request)

		// Enviar a oferta DHCP de volta ao cliente DHCP
		if offer != nil {
			s.conn.WriteTo(offer, addr)
		}

		// Adicione outros casos conforme necessário para outros tipos de mensagens DHCP
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
func (s *Server) buildDHCPOffer(request []byte) []byte {
	// Cria uma cópia da solicitação DHCP como base para a oferta DHCP
	offer := make([]byte, len(request))
	copy(offer, request)

	// Modifica o tipo de mensagem para DHCP Offer (código 2)
	offer[0] = 2

	// Identifica um endereço IP disponível usando o MongoDB
	availableIP, err := s.db.FindAvailableIP(s.ctx)
	if err != nil {
		logger.Error("Erro ao encontrar um endereço IP disponível:", err)
		return nil
	}

	// Converte o endereço IP de string para o formato de bytes IPv4
ip := net.ParseIP(availableIP).To4()
if ip == nil {
    logger.Error("Erro ao converter o endereço IP para formato IPv4")
    return nil
}

// Inclui o endereço IP na oferta DHCP (mínimo entre o tamanho do IP e o tamanho disponível na oferta)
copy(offer[16:16+len(ip)], ip)


	return offer
}

// buildDHCPAcknowledge constrói uma resposta DHCP Acknowledge
// func buildDHCPAcknowledge(request []byte) []byte {
// 	// TODO: Implementar a lógica para construir uma resposta DHCP Acknowledge
// 	// O exemplo abaixo simplesmente copia o pacote de solicitação como resposta
// 	acknowledge := make([]byte, len(request))
// 	copy(acknowledge, request)
// 	acknowledge[0] = 5 // Modifica o tipo de mensagem para DHCP Acknowledge
// 	return acknowledge
// }
