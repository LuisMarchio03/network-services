package main

import (
	"fmt"
	"net"
	"strings"
)

// Função para enviar uma requisição de resolução DNS para o servidor
func sendDNSRequest(domain string, serverAddr string) {
	// Explicação sobre o que é um cliente DNS
	fmt.Println("Cliente DNS ativo. O cliente enviará uma requisição para resolver o domínio:", domain)
	fmt.Println("Ele se conectará ao servidor DNS para receber o endereço IP correspondente.")

	// Criar um endereço UDP do servidor
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Erro ao resolver endereço do servidor:", err)
		return
	}

	// Criar uma conexão UDP
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Erro ao conectar com o servidor DNS:", err)
		return
	}
	defer conn.Close()

	// Enviar a requisição (nome do domínio)
	_, err = conn.Write([]byte(strings.TrimSpace(domain)))
	if err != nil {
		fmt.Println("Erro ao enviar requisição:", err)
		return
	}

	// Receber a resposta (endereço IP)
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao ler resposta do servidor:", err)
		return
	}

	response := strings.TrimSpace(string(buffer[:n]))
	fmt.Println("Resposta recebida do servidor DNS: O IP de", domain, "é", response)
}

// Função principal do cliente DNS
func main() {
	fmt.Println("Bem-vindo ao simulador de cliente DNS!")
	fmt.Println("Você pode simular o envio de uma requisição para resolver um domínio.")
	fmt.Println("O cliente enviará uma solicitação para o servidor DNS e obterá o endereço IP.")

	var domain, serverAddr string
	fmt.Print("Digite o domínio que deseja resolver: ")
	fmt.Scan(&domain)

	fmt.Print("Digite o endereço do servidor DNS (ex: localhost:53): ")
	fmt.Scan(&serverAddr)

	// Enviar a requisição DNS
	sendDNSRequest(domain, serverAddr)
}
