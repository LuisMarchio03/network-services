package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

var perguntas = []struct {
	Texto        string
	Alternativas []string
	Resposta     string
}{
	{"Qual a porta padrão do protocolo HTTP?", []string{"A) 80", "B) 21", "C) 25", "D) 53"}, "A"},
	{"Qual protocolo é usado para envio de emails?", []string{"A) HTTP", "B) SMTP", "C) FTP", "D) SSH"}, "B"},
	{"Qual camada do modelo OSI o protocolo IP pertence?", []string{"A) Física", "B) Rede", "C) Transporte", "D) Sessão"}, "B"},
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Servidor rodando na porta 8080...")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	go func() {
		<-shutdown
		fmt.Println("Servidor encerrando...")
		ln.Close()
		os.Exit(0)
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	totalAcertos := 0

	for i, pergunta := range perguntas {
		conn.Write([]byte(pergunta.Texto + "\n"))
		for _, alt := range pergunta.Alternativas {
			conn.Write([]byte(alt + "\n"))
		}
		conn.Write([]byte("Digite sua resposta (A, B, C, D): "))
		resposta, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler resposta:", err)
			return
		}
		resposta = strings.TrimSpace(resposta)

		infoCliente := strings.Split(resposta, "|")
		if len(infoCliente) != 5 {
			conn.Write([]byte("Formato de resposta inválido. Por favor, use o formato: MAC|IP Conectado|IP Cliente|Número da Questão|Resposta\n"))
			return
		}

		macAddress := infoCliente[0]
		ipConectado := infoCliente[1]
		ipCliente := infoCliente[2]
		numeroQuestao := infoCliente[3]
		respostaAluno := infoCliente[4]

		if strings.ToUpper(respostaAluno) == perguntas[i].Resposta {
			conn.Write([]byte("Resposta correta!\n"))
			totalAcertos++
		} else {
			conn.Write([]byte("Resposta incorreta!\n"))
		}

		log.Printf("[%s] Resposta do cliente [MAC: %s, IP conectado: %s, IP cliente: %s] Questão %s: %s\n", time.Now().Format(time.RFC3339), macAddress, ipConectado, ipCliente, numeroQuestao, respostaAluno)
	}

	media := float64(totalAcertos) / float64(len(perguntas)) * 10
	conn.Write([]byte(fmt.Sprintf("Sua média final é: %.2f\n", media)))
}
