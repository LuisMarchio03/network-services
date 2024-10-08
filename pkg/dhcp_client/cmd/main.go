package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	// Estabelecer conexão UDP com o servidor DHCP na porta 67
	conn, err := net.Dial("udp", "127.0.0.1:67")
	if err != nil {
		fmt.Printf("Erro ao conectar ao servidor DHCP: %v\n", err)
		return
	}
	defer conn.Close()

	// Construir e enviar uma mensagem DHCP Discover
	dhcpDiscoverMsg, err := buildDHCPDiscover()
	if err != nil {
		fmt.Printf("Erro ao construir mensagem DHCP Discover: %v\n", err)
		return
	}

	_, err = conn.Write(dhcpDiscoverMsg)
	if err != nil {
		fmt.Printf("Erro ao enviar mensagem DHCP Discover: %v\n", err)
		return
	}

	fmt.Println("Mensagem DHCP Discover enviada com sucesso")

	// Aguardar e processar a resposta do servidor DHCP
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Erro ao receber resposta DHCP: %v\n", err)
		return
	}
	_ = n

	// Imprimir o buffer recebido
	fmt.Printf("Buffer Recebido:\n%v\n", buffer)

	// Extrair e imprimir os dados relevantes da resposta DHCP
	offeredIP := net.IPv4(buffer[16], buffer[17], buffer[18], buffer[19])
	fmt.Printf("Endereço IP oferecido pelo servidor DHCP: %s\n", offeredIP.String())

	// Gateway padrão (bytes 4 a 7)
	gatewayIP := net.IPv4(buffer[4], buffer[5], buffer[6], buffer[7])
	fmt.Printf("Gateway Padrão: %s\n", gatewayIP.String())

	// Máscara de Sub-rede (bytes 1 a 3)
	subnetMask := net.IPv4(buffer[1], buffer[2], buffer[3], 0)
	fmt.Printf("Máscara de Sub-rede: %s\n", subnetMask.String())

	// Servidores DNS (bytes 20 a 27)
	dnsServer1 := net.IPv4(buffer[20], buffer[21], buffer[22], buffer[23])
	dnsServer2 := net.IPv4(buffer[24], buffer[25], buffer[26], buffer[27])
	fmt.Printf("Servidores DNS: %s, %s\n", dnsServer1.String(), dnsServer2.String())

	// Tempo de Locação (lease time) (bytes 244 a 247)
	leaseTime := binary.BigEndian.Uint32(buffer[244:248])
	fmt.Printf("Tempo de Locação (Lease Time): %d segundos\n", leaseTime)

	// Construir e enviar uma mensagem DHCP Request com o endereço IP oferecido
	dhcpRequestMsg, err := buildDHCPRequest(offeredIP)
	if err != nil {
		fmt.Printf("Erro ao construir mensagem DHCP Request: %v\n", err)
		return
	}

	_, err = conn.Write(dhcpRequestMsg)
	if err != nil {
		fmt.Printf("Erro ao enviar mensagem DHCP Request: %v\n", err)
		return
	}

	fmt.Println("Mensagem DHCP Request enviada com sucesso")

	// Aguardar e processar a resposta do servidor DHCP para a mensagem DHCP Request
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Printf("Erro ao receber resposta DHCP: %v\n", err)
		return
	}

	fmt.Printf("Resposta DHCP para DHCP Request (bytes recebidos: %d):\n", n)

	// Verifique se a resposta DHCP tem o tamanho esperado (300 bytes)
	if n != 300 {
		fmt.Printf("Erro: resposta DHCP com tamanho inválido (%d bytes)\n", n)
		return
	}

	// Imprimir o buffer recebido
	fmt.Printf("Buffer Recebido:\n%v\n", buffer)

	// Extrair e interpretar os campos da resposta DHCP
	offeredIP = net.IPv4(buffer[20], buffer[21], buffer[22], buffer[23])
	gatewayIP = net.IPv4(buffer[32], buffer[33], buffer[34], buffer[35])
	subnetMask = net.IPv4(buffer[4], buffer[5], buffer[6], buffer[7])
	dnsServer1 = net.IPv4(buffer[36], buffer[37], buffer[38], buffer[39])
	dnsServer2 = net.IPv4(buffer[24], buffer[25], buffer[26], buffer[27])
	leaseTime = binary.BigEndian.Uint32(buffer[244:248])

	// Imprimir os dados extraídos da resposta DHCP
	fmt.Printf("Endereço IP oferecido pelo servidor DHCP: %s\n", offeredIP.String())
	fmt.Printf("Gateway Padrão: %s\n", gatewayIP.String())
	fmt.Printf("Máscara de Sub-rede: %s\n", subnetMask.String())
	fmt.Printf("Servidores DNS: %s, %s\n", dnsServer1.String(), dnsServer2.String())
	fmt.Printf("Tempo de Locação (Lease Time): %d segundos\n", leaseTime)

	// Exemplo de uso da func buildDHCPRelease():

	//leaseIP := net.ParseIP("192.168.1.100")
	//clientMAC := "01:02:03:04:05:06"
	//
	//// Construir mensagem DHCP Release
	//dhcpReleaseMsg, err := buildDHCPRelease(clientMAC, leaseIP)
	//if err != nil {
	//	fmt.Printf("Erro ao construir mensagem DHCP Release: %v\n", err)
	//	return
	//}
	//
	//// Enviar mensagem DHCP Release para o servidor DHCP
	//_, err = conn.Write(dhcpReleaseMsg)
	//if err != nil {
	//	fmt.Printf("Erro ao enviar mensagem DHCP Release: %v\n", err)
	//	return
	//}
	//
	//fmt.Println("Mensagem DHCP Release enviada com sucesso")

}

// buildDHCPDiscover constrói uma mensagem DHCP Discover.
func buildDHCPDiscover() ([]byte, error) {
	var buf bytes.Buffer

	// Tipo de mensagem DHCP Discover (1)
	buf.WriteByte(1)
	// Tipo de hardware (Ethernet) e tamanho do endereço MAC (6 bytes)
	buf.WriteByte(1)
	buf.WriteByte(6)
	// Número de transação DHCP (XID)
	xid := uint32(123456789)
	binary.Write(&buf, binary.BigEndian, xid)
	// Segundos desde o boot (0)
	secs := uint16(0)
	binary.Write(&buf, binary.BigEndian, secs)
	// Flags (0)
	flags := uint16(0)
	binary.Write(&buf, binary.BigEndian, flags)
	// Endereço IP do cliente (0.0.0.0)
	clientIP := net.IPv4(0, 0, 0, 0).To4()
	binary.Write(&buf, binary.BigEndian, clientIP)
	// Endereço IP oferecido pelo servidor (0.0.0.0)
	yourIP := net.IPv4(0, 0, 0, 0).To4()
	binary.Write(&buf, binary.BigEndian, yourIP)
	// Endereço IP do servidor DHCP (0.0.0.0)
	serverIP := net.IPv4(0, 0, 0, 0).To4()
	binary.Write(&buf, binary.BigEndian, serverIP)
	// Endereço IP do gateway (0.0.0.0)
	gatewayIP := net.IPv4(0, 0, 0, 0).To4()
	binary.Write(&buf, binary.BigEndian, gatewayIP)
	// Endereço MAC do cliente (01:02:03:04:05:06)
	clientHardwareAddr, err := net.ParseMAC("01:02:03:04:05:06")
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear MAC: %v", err)
	}
	buf.Write(clientHardwareAddr)

	// Preencher o restante do buffer com zeros para completar até 300 bytes
	paddingLength := 300 - buf.Len()
	if paddingLength > 0 {
		buf.Write(make([]byte, paddingLength))
	}

	return buf.Bytes(), nil
}

// buildDHCPRequest constrói uma mensagem DHCP Request com base no endereço IP oferecido.
func buildDHCPRequest(offeredIP net.IP) ([]byte, error) {
	var buf bytes.Buffer

	// Tipo de mensagem DHCP Request (3)
	buf.WriteByte(3)

	// Tipo de hardware (Ethernet) e tamanho do endereço MAC (6 bytes)
	buf.WriteByte(1)
	buf.WriteByte(6)

	// Número de transação DHCP (XID)
	xid := uint32(123456789)
	binary.Write(&buf, binary.BigEndian, xid)

	// Segundos desde o boot (0)
	secs := uint16(0)
	binary.Write(&buf, binary.BigEndian, secs)

	// Flags (0)
	flags := uint16(0)
	binary.Write(&buf, binary.BigEndian, flags)

	// Endereço IP do cliente (0.0.0.0)
	clientIP := net.IPv4(0, 0, 0, 0).To4()
	binary.Write(&buf, binary.BigEndian, clientIP)

	// Endereço IP oferecido pelo servidor DHCP
	binary.Write(&buf, binary.BigEndian, offeredIP)

	// Endereço IP do servidor DHCP (0.0.0.0)
	serverIP := net.IPv4(0, 0, 0, 0).To4()
	binary.Write(&buf, binary.BigEndian, serverIP)

	// Endereço IP do gateway (0.0.0.0)
	gatewayIP := net.IPv4(0, 0, 0, 0).To4()
	binary.Write(&buf, binary.BigEndian, gatewayIP)

	// Endereço MAC do cliente (01:02:03:04:05:06)
	clientHardwareAddr, err := net.ParseMAC("01:02:03:04:05:06")
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear MAC: %v", err)
	}
	buf.Write(clientHardwareAddr)

	// Preencher o restante do buffer com zeros para completar até 300 bytes
	paddingLength := 300 - buf.Len()
	if paddingLength > 0 {
		buf.Write(make([]byte, paddingLength))
	}

	return buf.Bytes(), nil
}

// buildDHCPRelease constrói uma mensagem DHCP Release para liberar um endereço IP.
func buildDHCPRelease(clientMAC string, leaseIP net.IP) ([]byte, error) {
	var buf bytes.Buffer

	// Tipo de mensagem DHCP Release (7)
	buf.WriteByte(7)

	// Tipo de hardware (Ethernet) e tamanho do endereço MAC (6 bytes)
	buf.WriteByte(1)
	buf.WriteByte(6)

	// Número de transação DHCP (XID)
	xid := uint32(123456789)
	binary.Write(&buf, binary.BigEndian, xid)

	// Preencher com zeros para o restante do cabeçalho DHCP
	buf.Write(make([]byte, 202))

	// Endereço IP do cliente
	clientIP := net.IPv4(0, 0, 0, 0).To4()
	buf.Write(clientIP)

	// Endereço IP liberado (lease IP)
	buf.Write(leaseIP)

	// Endereço MAC do cliente
	clientHardwareAddr, err := net.ParseMAC(clientMAC)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear MAC: %v", err)
	}
	buf.Write(clientHardwareAddr)

	// Preencher o restante do buffer com zeros para completar até 300 bytes
	paddingLength := 300 - buf.Len()
	if paddingLength > 0 {
		buf.Write(make([]byte, paddingLength))
	}

	return buf.Bytes(), nil
}
