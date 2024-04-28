package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	// Configurar uma conexão de leitura UDP para receber respostas do servidor DHCP
	readConn, err := net.ListenPacket("udp", ":68")
	if err != nil {
		fmt.Printf("Erro ao configurar conexão de leitura UDP: %v\n", err)
		return
	}
	defer readConn.Close()

	// Definir um prazo de tempo para receber respostas do servidor DHCP
	readConn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Construir e enviar uma mensagem DHCP Discover
	discoverPacket := buildDHCPDiscover()
	_, err = readConn.WriteTo(discoverPacket, &net.UDPAddr{IP: net.IPv4bcast, Port: 67})
	if err != nil {
		fmt.Printf("Erro ao enviar mensagem DHCP Discover: %v\n", err)
		return
	}

	fmt.Println("Mensagem DHCP Discover enviada")

	// Aguardar uma resposta (DHCP Offer) do servidor DHCP
	offerPacket := make([]byte, 1024)
	n, _, err := readConn.ReadFrom(offerPacket)
	if err != nil {
		fmt.Printf("Erro ao receber oferta DHCP: %v\n", err)
		return
	}
	_ = n

	fmt.Println("Oferta DHCP recebida")

	// Extrair o endereço IP oferecido pelo servidor DHCP (yourIP)
	yourIP := net.IPv4(offerPacket[16], offerPacket[17], offerPacket[18], offerPacket[19])

	fmt.Println("offerPacket ", offerPacket)
	fmt.Println("yourIP ", yourIP)

	// Construir e enviar uma mensagem DHCP Request
	requestPacket := buildDHCPRequest(yourIP)
	_, err = readConn.WriteTo(requestPacket, &net.UDPAddr{IP: net.IPv4bcast, Port: 67})
	if err != nil {
		fmt.Printf("Erro ao enviar mensagem DHCP Request: %v\n", err)
		return
	}

	fmt.Println("Mensagem DHCP Request enviada")

	// Aguardar uma resposta (DHCP ACK) do servidor DHCP
	ackPacket := make([]byte, 1024)
	n, _, err = readConn.ReadFrom(ackPacket)
	if err != nil {
		fmt.Printf("Erro ao receber ACK DHCP: %v\n", err)
		return
	}

	fmt.Println("ACK DHCP recebido", ackPacket)

	// Extrair o endereço IP concedido pelo servidor DHCP
	var yourIPAck net.IP
	if len(ackPacket) >= 20 { // Verificar se há pelo menos 20 bytes no pacote
		yourIPAck = net.IP(ackPacket[20:24]) // Extrair os bytes do 21 ao 24 para o endereço IP
	}

	// Verificar se o endereço IP está corretamente definido
	if yourIPAck.To4() == nil {
		fmt.Println("Endereço IP concedido inválido pelo servidor DHCP (ACK)")
	} else {
		//fmt.Printf("Endereço IP concedido pelo servidor DHCP (ACK): %s\n", yourIPAck.String())
		fmt.Println("Endereço IP concedido pelo servidor DHCP (ACK)")
	}
}

// buildDHCPDiscover constrói uma mensagem DHCP Discover.
func buildDHCPDiscover() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(1))                                  // MessageType: DHCP Discover
	buf.WriteByte(byte(1))                                  // HardwareType: Ethernet
	buf.WriteByte(byte(6))                                  // HardwareAddrLength: 6 bytes (MAC address)
	buf.WriteByte(byte(0))                                  // Hops
	binary.Write(&buf, binary.BigEndian, uint32(123456789)) // Transaction ID
	binary.Write(&buf, binary.BigEndian, uint16(0))         // Seconds elapsed
	binary.Write(&buf, binary.BigEndian, uint16(0))         // Flags
	buf.Write(net.IPv4(0, 0, 0, 0).To4())                   // Client IP address
	buf.Write(net.IPv4(0, 0, 0, 0).To4())                   // Your IP address
	buf.Write(net.IPv4(0, 0, 0, 0).To4())                   // Server IP address
	buf.Write(net.IPv4(0, 0, 0, 0).To4())                   // Gateway IP address
	buf.Write(net.HardwareAddr{0, 1, 2, 3, 4, 5})           // Client hardware address
	buf.Write(make([]byte, 202-10-6))                       // Server host name (unused)
	buf.Write(make([]byte, 300-10-64))                      // Boot filename (unused)

	return buf.Bytes()
}

// buildDHCPRequest constrói uma mensagem DHCP Request com o endereço IP solicitado.
func buildDHCPRequest(requestedIP net.IP) []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(3))                                  // MessageType: DHCP Request
	buf.WriteByte(byte(1))                                  // HardwareType: Ethernet
	buf.WriteByte(byte(6))                                  // HardwareAddrLength: 6 bytes (MAC address)
	buf.WriteByte(byte(0))                                  // Hops
	binary.Write(&buf, binary.BigEndian, uint32(123456789)) // Transaction ID
	binary.Write(&buf, binary.BigEndian, uint16(0))         // Seconds elapsed
	binary.Write(&buf, binary.BigEndian, uint16(0))         // Flags
	buf.Write(net.IPv4(0, 0, 0, 0).To4())                   // Client IP address
	buf.Write(net.IPv4(0, 0, 0, 0).To4())                   // Your IP address
	buf.Write(net.IPv4(0, 0, 0, 0).To4())                   // Server IP address
	buf.Write(net.IPv4(0, 0, 0, 0).To4())                   // Gateway IP address
	buf.Write(net.HardwareAddr{0, 1, 2, 3, 4, 5})           // Client hardware address
	buf.Write(make([]byte, 202-10-6))                       // Server host name (unused)
	buf.Write(make([]byte, 300-10-64))                      // Boot filename (unused)
	buf.Write(requestedIP.To4())                            // Requested IP address

	return buf.Bytes()
}
