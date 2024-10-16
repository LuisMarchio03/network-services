package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	totalQuestoes := 3
	totalAcertos := 0

	fmt.Print("Digite seu MAC Address: ")
	macAddress, _ := reader.ReadString('\n')
	macAddress = strings.TrimSpace(macAddress)

	fmt.Print("Digite o IP conectado na rede: ")
	ipConectado, _ := reader.ReadString('\n')
	ipConectado = strings.TrimSpace(ipConectado)

	fmt.Print("Digite o IP do cliente: ")
	ipCliente, _ := reader.ReadString('\n')
	ipCliente = strings.TrimSpace(ipCliente)

	for i := 1; i <= totalQuestoes; i++ {
		// Receber pergunta do servidor
		pergunta, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(pergunta)

		// Ler a resposta do aluno
		fmt.Print("Digite sua resposta (A, B, C, D): ")
		respostaAluno, _ := reader.ReadString('\n')
		respostaAluno = strings.TrimSpace(respostaAluno)

		// Enviar resposta com informações do cliente
		enviar := fmt.Sprintf("%s|%s|%s|%d|%s", macAddress, ipConectado, ipCliente, i, respostaAluno)
		fmt.Fprintf(conn, enviar+"\n")

		// Receber feedback
		feedback, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(feedback)

		if strings.Contains(feedback, "correta") {
			totalAcertos++
		}
	}

	// Calcular a média final
	media := float64(totalAcertos) / float64(totalQuestoes) * 10
	fmt.Printf("Sua média final é: %.2f\n", media)
}
