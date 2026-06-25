package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

// loadZones carrega o mapa domínio->IP de um arquivo JSON.
func loadZones(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	zones := make(map[string]string)
	if err := json.Unmarshal(data, &zones); err != nil {
		return nil, err
	}
	return zones, nil
}

func startDNSServer(domainIPMap map[string]string) {
	addr := net.UDPAddr{Port: 53, IP: net.ParseIP("0.0.0.0")}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor DNS:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("DNS Server (didático) escutando em :53 com %d domínios\n", len(domainIPMap))

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao ler do buffer:", err)
			continue
		}
		request := strings.TrimSpace(string(buffer[:n]))
		ip, found := domainIPMap[request]
		if !found {
			ip = "0.0.0.0"
			fmt.Println("Domínio não encontrado:", request)
		} else {
			fmt.Println("Resolvido:", request, "->", ip)
		}
		if _, err := conn.WriteToUDP([]byte(ip), clientAddr); err != nil {
			fmt.Println("Erro ao responder:", err)
		}
	}
}

func main() {
	zonesFile := os.Getenv("DNS_ZONES_FILE")
	if zonesFile == "" {
		zonesFile = "zones.json"
	}
	zones, err := loadZones(zonesFile)
	if err != nil {
		fmt.Printf("Erro ao carregar zonas de %s: %v\n", zonesFile, err)
		os.Exit(1)
	}
	startDNSServer(zones)
}
