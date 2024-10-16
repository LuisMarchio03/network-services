package main

import (
	"fmt"
	"net"
	"time"
)

// Função para exibir uma explicação e aguardar ação do usuário
func promptUser(message string) {
	fmt.Println(message)
	fmt.Print("Pressione Enter para continuar...")
	fmt.Scanln()
}

func main() {
	// Introdução ao cliente DHCP
	fmt.Println("=== Cliente DHCP Simples ===")
	fmt.Println("O cliente DHCP precisa de um endereço IP para se conectar à rede.")
	fmt.Println("Ele enviará uma solicitação para o servidor DHCP na rede, pedindo um IP.")
	fmt.Println("O servidor irá responder com um IP disponível para o cliente utilizar.\n")

	// Configurando o MAC Address do cliente
	fmt.Println("Passo 1: Configurando o endereço MAC.")
	fmt.Println("O MAC Address é um identificador único para dispositivos na rede.")
	fmt.Println("Ele é usado pelo servidor DHCP para identificar dispositivos e atribuir IPs específicos.\n")

	// Configuração do MAC Address
	var macAddress string
	fmt.Print("Informe o MAC Address do cliente (ex: 00:1A:2B:3C:4D:5E): ")
	fmt.Scanln(&macAddress)

	fmt.Println("\nCliente configurado. Vamos enviar uma solicitação ao servidor DHCP.\n")

	// Explicação sobre o processo de descoberta DHCP
	promptUser("Etapa 1: Enviando mensagem DHCP Discover para encontrar servidores DHCP na rede...")

	// Configurando o endereço do servidor DHCP
	serverAddr := net.UDPAddr{
		Port: 67,                       // Porta padrão do servidor DHCP
		IP:   net.ParseIP("127.0.0.1"), // Endereço do servidor
	}

	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		fmt.Println("Erro ao se conectar ao servidor DHCP:", err)
		return
	}
	defer conn.Close()

	// Simulando envio de solicitação DHCP Discover
	time.Sleep(2 * time.Second)

	// Enviando solicitação DHCP Discover para o servidor
	message := fmt.Sprintf("DHCP Discover: Cliente MAC: %s", macAddress)
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Erro ao enviar mensagem DHCP Discover:", err)
		return
	}

	fmt.Println("Mensagem DHCP Discover enviada. Aguardando resposta...\n")

	// Recebendo a oferta de IP do servidor
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao receber oferta de IP do servidor:", err)
		return
	}

	// Exibindo a resposta do servidor DHCP
	fmt.Printf("Resposta recebida do servidor: %s\n", string(buffer[:n]))

	// Explicação sobre a aceitação do IP
	promptUser("Etapa 2: Enviando DHCP Request para confirmar o IP...\n")

	// Enviando mensagem DHCP Request
	requestMessage := fmt.Sprintf("DHCP Request: Cliente MAC: %s solicita o IP oferecido", macAddress)
	_, err = conn.Write([]byte(requestMessage))
	if err != nil {
		fmt.Println("Erro ao enviar mensagem DHCP Request:", err)
		return
	}

	fmt.Println("Mensagem DHCP Request enviada. Aguardando confirmação...\n")

	// Recebendo confirmação DHCP Acknowledge do servidor
	n, _, err = conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao receber confirmação do servidor:", err)
		return
	}

	// Exibindo a resposta do servidor DHCP
	fmt.Printf("Resposta recebida do servidor: %s\n", string(buffer[:n]))

	// Finalização com sucesso
	fmt.Println("\nProcesso DHCP finalizado com sucesso. O IP foi atribuído ao cliente.")
}
