package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"ntp/ntptime"
)

func main() {
	server := flag.String("server", "ntp-server:123", "endereço do servidor NTP")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *server)
	if err != nil {
		fmt.Println("Erro ao resolver servidor:", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Erro ao conectar:", err)
		os.Exit(1)
	}
	defer conn.Close()

	req := make([]byte, 48)
	req[0] = 0x23 // LI=0, VN=4, Mode=3 (client)
	if _, err := conn.Write(req); err != nil {
		fmt.Println("Erro ao enviar:", err)
		os.Exit(1)
	}
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	resp := make([]byte, 48)
	if _, err := conn.Read(resp); err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		os.Exit(1)
	}
	secs := binary.BigEndian.Uint32(resp[40:44])
	frac := binary.BigEndian.Uint32(resp[44:48])
	serverTime := ntptime.FromNTP(secs, frac)
	fmt.Printf("NTP: hora do servidor = %s (offset ~ %s)\n",
		serverTime.Format(time.RFC3339Nano), time.Since(serverTime).Round(time.Millisecond))
}
