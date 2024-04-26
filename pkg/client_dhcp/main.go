package main

import (
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

	// Simular uma mensagem DHCP Discover (tipo 1)
	messageType := byte(1)
	message := []byte{messageType}

	// Enviar mensagem DHCP para o servidor
	_, err = conn.Write(message)
	if err != nil {
		fmt.Printf("Erro ao enviar mensagem DHCP: %v\n", err)
		return
	}

	fmt.Println("Mensagem DHCP Discover enviada com sucesso")

	// Aguardar resposta do servidor
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Printf("Erro ao receber resposta DHCP: %v\n", err)
		return
	}

	fmt.Println("Resposta DHCP recebida:", buffer)
}
