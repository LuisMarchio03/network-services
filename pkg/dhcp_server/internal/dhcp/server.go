package dhcp

import (
	"context"
	"encoding/binary"
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

	case 3: // DHCP Request
		// Processa uma mensagem DHCP Request recebida e envia uma resposta DHCP Acknowledge
		ack := s.buildDHCPAcknowledge(request)
		if ack != nil {
			s.conn.WriteTo(ack, addr)
		}

	case 5: // DHCP Renew
		// Processa uma mensagem DHCP Renew recebida e renova o endereço IP do cliente, se aplicável
		ack := s.buildDHCPAcknowledge(request)
		if ack != nil {
			s.conn.WriteTo(ack, addr)
		}

	case 7: // DHCP Release
		// Processa uma mensagem DHCP Release recebida e libera o endereço IP associado ao cliente
		s.releaseDHCPAddress(request)

	default:
		// Se o tipo de mensagem DHCP não for reconhecido, registra um aviso
		logger.Info("Tipo de mensagem DHCP não suportado: ")
		fmt.Print(messageType)
	}
}

// buildDHCPOffer constrói uma resposta DHCP Offer em resposta a uma solicitação DHCP Discover.
// Esta função é responsável por identificar um endereço IP disponível para oferecer ao cliente,
// construir a estrutura do pacote DHCP Offer com parâmetros adicionais e enviá-lo de volta ao cliente.
// Parâmetros:
//   - request: O pacote DHCP Discover enviado pelo cliente.
//
// Retorna:
//   - Um pacote DHCP Offer pronto para ser enviado ao cliente.
func (s *Server) buildDHCPOffer(request []byte) []byte {
	// Verifica se a solicitação DHCP tem o tamanho mínimo necessário
	if len(request) < 20 {
		logger.Info("Solicitação DHCP inválida: tamanho insuficiente")
		return nil
	}

	// Cria uma oferta DHCP com tamanho suficiente para conter todos os dados necessários
	offer := make([]byte, 300) // Por exemplo, oferece espaço para 300 bytes

	// Copia o tipo de mensagem para DHCP Offer (código 2)
	offer[0] = 2

	// Identifica um endereço IP disponível usando o MongoDB
	//availableIP, err := s.db.FindAvailableIP(s.ctx)
	availableIP := "192.168.2.42"
	//if err != nil {
	//	logger.Error("Erro ao encontrar um endereço IP disponível:", err)
	//	return nil
	//}

	// Converte o endereço IP de string para o formato de bytes IPv4
	ip := net.ParseIP(availableIP).To4()

	// Verifica se o endereço IP convertido tem o tamanho esperado (IPv4 = 4 bytes)
	if len(ip) != 4 {
		logger.Info("Endereço IP inválido")
		return nil
	}

	// Copia o endereço IP para a oferta DHCP
	copy(offer[16:20], ip)

	// Adiciona o tempo de locação (lease time) para o endereço IP (por exemplo, 3600 segundos)
	leaseTime := uint32(3600)
	binary.BigEndian.PutUint32(offer[244:248], leaseTime)

	// Adiciona o gateway padrão (por exemplo, 192.168.1.1) à oferta DHCP
	gatewayIP := net.ParseIP("192.168.1.1").To4()
	copy(offer[4:8], gatewayIP)

	// Adiciona a máscara de sub-rede (por exemplo, 255.255.255.0) à oferta DHCP
	subnetMask := net.IPv4Mask(255, 255, 255, 0)
	copy(offer[1:4], subnetMask)

	// Adiciona os servidores DNS (por exemplo, 8.8.8.8 e 8.8.4.4) à oferta DHCP
	dnsServer1 := net.ParseIP("8.8.8.8").To4()
	copy(offer[20:24], dnsServer1)

	dnsServer2 := net.ParseIP("8.8.4.4").To4()
	copy(offer[24:28], dnsServer2)

	return offer
}

// buildDHCPAcknowledge constrói uma resposta DHCP Acknowledge em resposta a uma solicitação DHCP.
// Esta função examina a solicitação DHCP recebida e constrói uma resposta DHCP Acknowledge contendo os parâmetros apropriados,
// como o endereço IP oferecido ou renovado.
// Parâmetros:
//   - request: O pacote DHCP Request enviado pelo cliente DHCP.
//
// Retorna:
//   - Um pacote DHCP Acknowledge pronto para ser enviado ao cliente DHCP em resposta à solicitação.
func (s *Server) buildDHCPAcknowledge(request []byte) []byte {
	// Verifica se a solicitação DHCP tem o tamanho mínimo necessário
	if len(request) < 240 {
		logger.Info("Solicitação DHCP inválida: tamanho insuficiente")
		return nil
	}

	// Validar o tipo de mensagem DHCP (Request) - Código 3
	if request[0] != 3 {
		logger.Info("Tipo de mensagem DHCP inválido para Acknowledge")
		return nil
	}

	// Extrair o endereço IP solicitado pelo cliente da solicitação DHCP Request
	requestedIP := net.IP(request[12:16])

	// Cria uma resposta DHCP Acknowledge com espaço suficiente para conter todos os dados necessários
	ack := make([]byte, 300)

	// Copia o tipo de mensagem para DHCP Acknowledge (código 5)
	ack[0] = 5

	// Copia o endereço IP do servidor DHCP para a resposta DHCP Acknowledge
	serverIP := net.ParseIP("192.168.2.130").To4() // Substitua pelo endereço IP do servidor DHCP
	copy(ack[20:24], serverIP)

	// Copia o endereço IP solicitado pelo cliente para a resposta DHCP Acknowledge
	copy(ack[16:20], requestedIP)

	// Define o tempo de locação (lease time) para o endereço IP (em segundos)
	leaseTime := uint32(3600) // Exemplo: 3600 segundos (1 hora) para uma locação temporária
	binary.BigEndian.PutUint32(ack[244:248], leaseTime)

	// Adicionar gateway padrão, máscara de sub-rede, servidor DNS, etc., conforme necessário
	defaultGateway := net.ParseIP("192.168.1.254").To4() // Exemplo de gateway padrão
	copy(ack[32:36], defaultGateway)

	subnetMask := net.IPv4Mask(255, 255, 255, 0) // Exemplo de máscara de sub-rede
	copy(ack[4:8], subnetMask)

	dnsServer := net.ParseIP("8.8.8.8").To4() // Exemplo de servidor DNS
	copy(ack[36:40], dnsServer)

	// Retorne a resposta DHCP Acknowledge construída
	return ack
}

// releaseDHCPAddress libera o endereço IP associado a uma mensagem DHCP Release.
// Esta função é chamada quando um cliente DHCP envia uma mensagem DHCP Release para liberar seu endereço IP.
// A função identifica o endereço IP a ser liberado com base na mensagem DHCP Release recebida e realiza as operações necessárias
// para liberar o endereço IP no sistema de gerenciamento de endereços IP (por exemplo, marcando-o como disponível novamente).
// Parâmetros:
//   - request: O pacote DHCP Release enviado pelo cliente DHCP.
func (s *Server) releaseDHCPAddress(request []byte) {
	// Verifica se a solicitação DHCP tem o tamanho mínimo necessário
	if len(request) < 240 {
		logger.Info("Solicitação DHCP Release inválida: tamanho insuficiente")
		return
	}

	// Extrai o endereço IP a ser liberado da mensagem DHCP Release
	ipAddress := net.IP(request[4:8]) // Supondo que o endereço IP a ser liberado está na posição correta no pacote DHCP Release

	// Libera o endereço IP no sistema de gerenciamento de endereços IP (por exemplo, marca como disponível novamente no banco de dados)
	if err := s.db.ReleaseIPAddress(s.ctx, ipAddress.String()); err != nil {
		logger.Error("Erro ao liberar endereço IP:", err)
		return
	}

	fmt.Print("Endereço IP", ipAddress.String(), "foi liberado com sucesso")

	// Aqui, você pode adicionar outras operações necessárias após a liberação do endereço IP, se necessário
}
