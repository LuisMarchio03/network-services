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
// Ele inicializa um socket UDP na porta 67 para escutar mensagens DHCP.
// Parâmetros:
//   - ctx: Contexto de execução para controle de tempo e cancelamento.
//   - db: Instância do MongoDB para interação com o banco de dados.
//
// Retorna:
//   - Um ponteiro para a estrutura Server e um possível erro, se houver.
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
// Ele fica em um loop infinito para ler e processar as mensagens DHCP recebidas.
// Fecha o socket ao encerrar a execução.
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
//   - addr: O endereço do cliente que enviou a mensagem DHCP.
func (s *Server) processDHCPMessage(request []byte, addr net.Addr) {
	// Interpretar o tipo de mensagem DHCP
	messageType := request[0]

	// Processar a mensagem DHCP com base no tipo
	switch messageType {
	case 1: // DHCP Discover
		// Construir e enviar uma oferta DHCP em resposta ao DHCP Discover
		offer := s.buildDHCPOffer(request)
		if offer != nil {
			s.conn.WriteTo(offer, addr)
		}
	case 3: // DHCP Request
		// Extrair o endereço IP solicitado da mensagem DHCP Request
		requestedIP := net.IP(request[12:16]).String()

		// Verificar se o endereço IP solicitado está disponível
		if s.isIPAvailable(requestedIP) {
			// Construir e enviar uma mensagem DHCP ACK (Acknowledge)
			ack := s.buildDHCPAck(request, requestedIP)
			s.conn.WriteTo(ack, addr)
		} else {
			// Construir e enviar uma mensagem DHCP NAK (Not Acknowledge)
			nak := s.buildDHCPNak(request)
			s.conn.WriteTo(nak, addr)
		}
	}
}

// isIPAvailable verifica se o endereço IP especificado está disponível para concessão.
func (s *Server) isIPAvailable(ip string) bool {
	// Buscar ctx
	ctx := s.ctx

	// Consultar o MongoDB para verificar se o endereço IP está disponível
	_, err := s.db.FindIP(ctx, ip)
	if err != nil {
		// Endereço IP não encontrado no banco de dados, considerado disponível
		return true
	}
	// Endereço IP encontrado no banco de dados, considerado indisponível
	return false
}

// buildDHCPOffer constrói uma resposta DHCP Offer em resposta a uma solicitação DHCP Discover.
// Parâmetros:
//   - request: O pacote DHCP Discover enviado pelo cliente.
//
// Retorna:
//   - Um pacote DHCP Offer pronto para ser enviado ao cliente.
func (s *Server) buildDHCPOffer(request []byte) []byte {
	// Verifica a validade da solicitação DHCP
	if len(request) < 20 {
		logger.Info("Solicitação DHCP inválida: tamanho insuficiente")
		return nil
	}

	// Constrói uma oferta DHCP com um endereço IP disponível
	offer := make([]byte, 20) // Por exemplo, oferece espaço para 20 bytes
	offer[0] = 2              // Define o tipo de mensagem como DHCP Offer

	// Encontra um endereço IP disponível usando o MongoDB
	availableIP, err := s.db.FindAvailableIP(s.ctx)
	if err != nil {
		logger.Error("Erro ao encontrar um endereço IP disponível:", err)
		return nil
	}

	// Converte o endereço IP para o formato correto e o adiciona à oferta DHCP
	ip := net.ParseIP(availableIP).To4()
	if len(ip) != 4 {
		logger.Info("Endereço IP inválido")
		return nil
	}
	copy(offer[16:20], ip)

	err = s.db.UpdateIPAssignment(s.ctx, availableIP, true)
	if err != nil {
		logger.Error("Erro ao atualizar status do IP:", err)
		return nil
	}

	return offer
}

// buildDHCPAck constrói uma mensagem DHCP Acknowledge (ACK) em resposta a uma solicitação DHCP Request.
func (s *Server) buildDHCPAck(request []byte, ip string) []byte {
	// Constrói a mensagem DHCP ACK
	ack := make([]byte, 20) // Tamanho mínimo para uma mensagem DHCP ACK
	ack[0] = 5              // Define o tipo de mensagem como DHCP ACK
	copy(ack[16:20], net.ParseIP(ip).To4())

	return ack
}

// buildDHCPNak constrói uma mensagem DHCP Not Acknowledge (NAK) em resposta a uma solicitação DHCP Request.
func (s *Server) buildDHCPNak(request []byte) []byte {
	// Constrói a mensagem DHCP NAK
	nak := make([]byte, 20) // Tamanho mínimo para uma mensagem DHCP NAK
	nak[0] = 6              // Define o tipo de mensagem como DHCP NAK

	return nak
}
