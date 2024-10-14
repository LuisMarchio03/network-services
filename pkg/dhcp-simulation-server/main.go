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
	fmt.Println("Quando um cliente (como um computador ou celular) se conecta à rede, ele usa o DHCP para pedir um IP, e o servidor DHCP fornece um IP disponível.")
	fmt.Println("Agora você vai configurar o servidor DHCP.\n")

	// Explicação sobre o intervalo de IPs e a máscara de sub-rede
	fmt.Println("Passo 1: Definindo o intervalo de IPs.")
	fmt.Println("O intervalo de IPs define quais endereços podem ser atribuídos aos clientes.")
	fmt.Println("Por exemplo, você pode definir que os IPs devem começar em 192.168.1.100 e ir até 192.168.1.150.\n")

	// Configurando o intervalo de IPs
	ipStart := promptUser("Informe o endereço IP inicial do intervalo que será atribuído aos clientes (ex: 192.168.1.100):")
	subnetMask := promptUser("Informe a máscara de sub-rede (ex: 255.255.255.0):")

	fmt.Println("\nPasso 2: Definindo a porta UDP.")
	fmt.Println("O DHCP usa o protocolo UDP para comunicação. O servidor precisa escutar as solicitações dos clientes em uma porta específica.")
	fmt.Println("Por padrão, o DHCP usa a porta 67 para o servidor e a 68 para o cliente.\n")

	// Definindo a porta onde o servidor vai escutar
	port := promptUser("Informe a porta em que o servidor deve escutar (ex: 8080 ou deixe em branco para usar a porta 67):")
	if port == "" {
		port = "67"
	}

	// Configuração do endereço de escuta do servidor
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
	fmt.Println("Agora, o servidor está aguardando solicitações de clientes DHCP.\n")
	fmt.Println("Etapas do DHCP:")
	fmt.Println("1. **Discover**: O cliente envia uma solicitação na rede para descobrir servidores DHCP.")
	fmt.Println("2. **Offer**: O servidor responde oferecendo um IP.")
	fmt.Println("3. **Request**: O cliente solicita formalmente o IP oferecido.")
	fmt.Println("4. **Acknowledge**: O servidor confirma a atribuição do IP e o cliente pode usá-lo.\n")

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

		// Enviando a resposta para o cliente
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			fmt.Println("Erro ao enviar a oferta de IP:", err)
			continue
		}

		fmt.Printf("Oferta de IP %s enviada para o cliente %s.\n", offeredIP, clientAddr)
		time.Sleep(5 * time.Second) // Simulação de espera para a próxima solicitação
	}
}
