package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	domain := flag.String("domain", "example.com", "domínio a resolver")
	server := flag.String("server", "dns-server:53", "endereço do servidor DNS")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *server)
	if err != nil {
		fmt.Println("Erro ao resolver endereço do servidor:", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Erro ao conectar com o servidor DNS:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(strings.TrimSpace(*domain))); err != nil {
		fmt.Println("Erro ao enviar requisição:", err)
		os.Exit(1)
	}
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		os.Exit(1)
	}
	fmt.Printf("DNS: %s -> %s\n", *domain, strings.TrimSpace(string(buffer[:n])))
}
