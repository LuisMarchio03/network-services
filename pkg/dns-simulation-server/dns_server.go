package main

import (
	"fmt"
	"net"
	"strings"
)

// Função para iniciar o servidor DNS
func startDNSServer(port string, domainIPMap map[string]string) {
	// Explicação sobre o que é o DNS e como ele resolve nomes de domínio
	fmt.Println("DNS Server iniciado...")
	fmt.Println("O DNS (Domain Name System) é responsável por traduzir nomes de domínio em endereços IP.")
	fmt.Println("O servidor DNS receberá pedidos para resolução de nomes e retornará o endereço IP correspondente.")

	// Iniciar um listener UDP
	addr := net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor DNS:", err)
		return
	}
	defer conn.Close()

	// Loop para receber requisições de DNS
	buffer := make([]byte, 1024)
	for {
		// Receber uma requisição de um cliente
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao ler do buffer:", err)
			continue
		}

		request := strings.TrimSpace(string(buffer[:n]))
		fmt.Println("Requisição recebida para resolver o domínio:", request)

		// Procurar o IP associado ao domínio solicitado
		ip, found := domainIPMap[request]
		if found {
			fmt.Println("Domínio encontrado:", request, "->", ip)
		} else {
			ip = "0.0.0.0" // Caso o domínio não seja encontrado
			fmt.Println("Domínio não encontrado:", request)
		}

		// Enviar a resposta de volta ao cliente
		_, err = conn.WriteToUDP([]byte(ip), clientAddr)
		if err != nil {
			fmt.Println("Erro ao enviar resposta para o cliente:", err)
		}
	}
}

// Função principal do servidor DNS
func main() {
	domainIPMap := make(map[string]string)

	// Explicação sobre a configuração inicial
	fmt.Println("Bem-vindo ao simulador de servidor DNS!")
	fmt.Println("Você configurará domínios e seus IPs correspondentes.")
	fmt.Println("Isso permitirá que o servidor resolva nomes quando solicitado por um cliente.")

	// Configurar domínios e IPs
	for {
		var domain, ip string
		fmt.Print("Digite o domínio (ou 'sair' para finalizar): ")
		fmt.Scan(&domain)
		if domain == "sair" {
			break
		}

		fmt.Print("Digite o endereço IP para o domínio '" + domain + "': ")
		fmt.Scan(&ip)
		domainIPMap[domain] = ip
	}

	// Iniciar o servidor DNS
	fmt.Println("Iniciando o servidor DNS na porta 53...")
	startDNSServer("53", domainIPMap)
}
