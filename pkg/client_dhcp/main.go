package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("udp", "localhost:67")
	if err != nil {
		fmt.Printf("Erro ao conectar ao servidor DHCP: %v\n", err)
		return
	}
	defer conn.Close()

	// Construir uma mensagem DHCP Discover
	messageType := byte(1)                    // Tipo 1 indica DHCP Discover
	hardwareType := byte(1)                   // Tipo 1 indica endereço MAC Ethernet
	hardwareLength := byte(6)                 // Tamanho do endereço MAC
	xid := uint32(123456789)                  // Número de transação DHCP
	secs := uint16(0)                         // Segundos desde que o processo de boot começou
	flags := uint16(0)                        // Flags (reservado, não usado)
	clientIP := net.IPv4(0, 0, 0, 0).To4()     // Endereço IP do cliente (0.0.0.0)
	yourIP := net.IPv4(0, 0, 0, 0).To4()       // Endereço IP oferecido pelo servidor (0.0.0.0)
	serverIP := net.IPv4(0, 0, 0, 0).To4()     // Endereço IP do servidor DHCP (0.0.0.0)
	gatewayIP := net.IPv4(0, 0, 0, 0).To4()    // Endereço IP do gateway (0.0.0.0)
	clientHardwareAddr, _ := net.ParseMAC("01:02:03:04:05:06") // Endereço MAC do cliente

	// Construir a mensagem DHCP Discover como um buffer de bytes
	var buf bytes.Buffer
	buf.WriteByte(messageType)
	buf.WriteByte(hardwareType)
	buf.WriteByte(hardwareLength)
	buf.WriteByte(0)                          // Hops (não usado)
	binary.Write(&buf, binary.BigEndian, xid)
	binary.Write(&buf, binary.BigEndian, secs)
	binary.Write(&buf, binary.BigEndian, flags)
	binary.Write(&buf, binary.BigEndian, clientIP)
	binary.Write(&buf, binary.BigEndian, yourIP)
	binary.Write(&buf, binary.BigEndian, serverIP)
	binary.Write(&buf, binary.BigEndian, gatewayIP)
	binary.Write(&buf, binary.BigEndian, clientHardwareAddr)
	buf.Write(make([]byte, 202-10-6)) // Preencher com zeros para completar até 300 bytes (tamanho mínimo)

	// Enviar mensagem DHCP para o servidor
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Erro ao enviar mensagem DHCP: %v\n", err)
		return
	}

	fmt.Println("Mensagem DHCP Discover enviada com sucesso")

	// Aguardar resposta do servidor
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Erro ao receber resposta DHCP: %v\n", err)
		return
	}

	_ = n

	// Extrair o endereço IP "yourIP" da resposta DHCP
	yourIP = net.IPv4(buffer[16], buffer[17], buffer[18], buffer[19])

	fmt.Printf("Endereço IP yourIP: %s\n", yourIP.String())
}
