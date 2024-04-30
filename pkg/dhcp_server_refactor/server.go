package dhcp_server_refactor

import (
	"context"
	"fmt"
	"net"

	"network-services-server-dhcp/mongodb"
)

// net.PacketConn:
// - net.PacketConn é uma interface definida no pacote net do Go.
// - Essa interface representa uma conexão de pacotes genérica.
// - PacketConn é implementada por vários tipos de conexões de rede, incluindo conexões UDP (net.UDPConn) e conexões IP (net.IPConn).
// - Oferece métodos para ler e escrever pacotes genéricos, independentemente do protocolo de rede subjacente.

// net.UDPConn:
// - net.UDPConn é um tipo de conexão específico para o protocolo UDP (User Datagram Protocol).
// - Essa estrutura implementa a interface net.PacketConn e fornece métodos adicionais específicos do UDP, como ReadFromUDP e WriteToUDP para leitura e escrita de dados UDP.
// - Permite a criação de sockets UDP para comunicação de rede.

type Server struct {
	conn *net.UDPConn
	db   *mongodb.MongoDB
	ctx  context.Context
}

// NewServer... cria um novo servidor DHCP
// O Server DHCP vai rodar na porta 67 do tipo UDP

// Uma conexão UDP (User Datagram Protocol) é uma forma de comunicação de rede que opera no modelo de transferência de dados não confiável.
// - O UDP é um protocolo de camada de transporte que permite que aplicativos enviem datagramas (pacotes) sem conexão pela rede.
// - Ele é mais leve e mais rápido do que o TCP, mas não garante a entrega ou a ordem dos pacotes.
// - Uma conexão UDP é identificada por um par de endereços IP e números de porta (local e remoto).

// Exemplo de uso:
// - Para criar uma conexão UDP e enviar dados:
//     conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234})
//     if err != nil {
//         fmt.Println("Erro ao abrir conexão UDP:", err)
//         return
//     }
//     defer conn.Close()
//     message := []byte("Olá, mundo!")
//     _, err = conn.Write(message)
//     if err != nil {
//         fmt.Println("Erro ao enviar mensagem UDP:", err)
//         return
//     }

// - Para receber dados em uma conexão UDP:
//     buffer := make([]byte, 1024)
//     n, addr, err := conn.ReadFromUDP(buffer)
//     if err != nil {
//         fmt.Println("Erro ao receber mensagem UDP:", err)
//         return
//     }
//     fmt.Printf("Recebido %d bytes de %s: %s\n", n, addr.String(), string(buffer[:n]))

// A Func NewServer vai receber como parâmetros:
//	-- ctx: Contexto de execução para controle de tempo e cancelamento do MongoDB
//	-- db: Instância do MongoDB para interação com o banco de dados.

// Por fim, o NewServer vai retornar:
// -- Um ponteiro para a estrutura Server e um possível erro, caso houver...

func NewServer(ctx context.Context, db *mongodb.MongoDB) (*Server, error) {
	// 1. Primeiramente, precisamos inicializar o Socket UDP na porta 67

	// Um socket UDP é uma interface de comunicação para enviar e receber datagramas UDP (User Datagram Protocol).
	// - O UDP é um protocolo de camada de transporte que opera sem conexão, o que significa que os datagramas são enviados de forma independente.
	// - Os sockets UDP permitem que um programa envie e receba dados em formato de datagrama pela rede.
	// - Cada socket UDP é identificado por um endereço IP e um número de porta local.

	// Exemplo de criação de um socket UDP:
	//     conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 1234})
	//     if err != nil {
	//         fmt.Println("Erro ao criar socket UDP:", err)
	//         return
	//     }
	//     defer conn.Close()

	// No exemplo acima:
	// - net.ListenUDP cria um socket UDP que escuta conexões na porta especificada (1234 neste caso).
	// - Depois de usar o socket, é importante fechá-lo para liberar os recursos.
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir socket UDP: %v", err) // Printa uma msg e espera receber tbm o erro que foi detectado.
	}

	// Criar uma nova instância do servidor DHCP
	server := &Server{
		conn: conn,
		db:   db,
		ctx:  ctx,
	}

	fmt.Print("Servidor DHCP iniciado")

	return server, nil
}
