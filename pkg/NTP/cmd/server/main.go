package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"ntp/ntptime"
)

func main() {
	addr := net.UDPAddr{Port: 123, IP: net.ParseIP("0.0.0.0")}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Erro ao iniciar servidor NTP:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("NTP Server (didático) escutando em :123")

	buf := make([]byte, 48)
	for {
		n, client, err := conn.ReadFromUDP(buf)
		if err != nil || n == 0 {
			continue
		}
		resp := make([]byte, 48)
		// LI=0, VN=4, Mode=4 (server) -> 0b00_100_100 = 0x24
		resp[0] = 0x24
		resp[1] = 1 // stratum
		secs, frac := ntptime.ToNTP(time.Now().UTC())
		// Transmit Timestamp (bytes 40..47)
		binary.BigEndian.PutUint32(resp[40:44], secs)
		binary.BigEndian.PutUint32(resp[44:48], frac)
		if _, err := conn.WriteToUDP(resp, client); err != nil {
			fmt.Println("Erro ao responder NTP:", err)
		}
	}
}
