package main

import (
	"fmt"
	"net"
	"time"
)

// Função para solicitar uma entrada do usuário com explicação
func promptUser(message string) string {
	var input string
	fmt.Println(message)
	fmt.Print("Digite o valor e pressione Enter: ")
	fmt.Scanln(&input)
	return input
}

func main() {
	// Introdução e explicação sobre o que o servidor DHCP faz
	fmt.Println("=== Bem-vindo ao simulador explicativo do Servidor DHCP ===")
	fmt.Println("O DHCP (Dynamic Host Configuration Protocol) é um protocolo que automatiza a atribuição de endereços IP aos dispositivos conectados em uma rede.")
	fmt.Println("Agora você vai configurar o servidor DHCP.\n")

	// Configurando o intervalo de IPs e máscara de sub-rede
	ipStart := promptUser("Informe o endereço IP inicial do intervalo que será atribuído aos clientes (ex: 192.168.1.100):")
	subnetMask := promptUser("Informe a máscara de sub-rede (ex: 255.255.255.0):")

	// Definindo a porta onde o servidor vai escutar
	port := promptUser("Informe a porta em que o servidor deve escutar (ex: 8080 ou deixe em branco para usar a porta 67):")
	if port == "" {
		port = "67"
	}

	addr := net.UDPAddr{
		Port: 67, // Porta padrão do DHCP ou a definida pelo usuário
		IP:   net.ParseIP("0.0.0.0"),
	}

	// Iniciando o servidor para escutar na porta UDP
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor DHCP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("\nConfiguração concluída. O servidor DHCP está pronto para fornecer endereços IP.\n")
	buffer := make([]byte, 1024)

	// Loop para escutar mensagens de clientes
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao ler a solicitação do cliente:", err)
			continue
		}

		// Exibindo a solicitação recebida
		fmt.Printf("Solicitação recebida de %s: %s\n", clientAddr, string(buffer[:n]))

		// Simulando oferta de IP
		offeredIP := ipStart
		response := fmt.Sprintf("DHCP Offer: Oferecendo IP: %s com máscara: %s", offeredIP, subnetMask)
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			fmt.Println("Erro ao enviar a oferta de IP:", err)
			continue
		}

		fmt.Printf("Oferta de IP %s enviada para o cliente %s.\n", offeredIP, clientAddr)

		// Receber DHCP Request
		n, clientAddr, err = conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao receber DHCP Request do cliente:", err)
			continue
		}
		fmt.Printf("Solicitação DHCP Request recebida de %s: %s\n", clientAddr, string(buffer[:n]))

		// Confirmar com DHCP Acknowledge
		ackMessage := fmt.Sprintf("DHCP ACK: IP %s confirmado para o cliente", offeredIP)
		_, err = conn.WriteToUDP([]byte(ackMessage), clientAddr)
		if err != nil {
			fmt.Println("Erro ao enviar DHCP Acknowledge:", err)
			continue
		}

		fmt.Printf("DHCP Acknowledge enviado para o cliente %s. IP %s atribuído com sucesso.\n", clientAddr, offeredIP)
		time.Sleep(5 * time.Second) // Simulação de espera para a próxima solicitação
	}
}
