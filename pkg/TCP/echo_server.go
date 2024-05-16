package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	PORT = ":8888"
)

// handleConnection lida com a conexão de um cliente TCP.
// - A função lê mensagens do cliente, ecoa de volta para o cliente e imprime a mensagem recebida.
// - A conexão é fechada ao final da função ou em caso de erro.
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Lê a mensagem do cliente até encontrar um caractere de nova linha ('\n').
		message, err := reader.ReadString('\n')
		if err != nil {
			// Se houver um erro durante a leitura, considera que o cliente se desconectou.
			fmt.Println("Cliente desconectado.")
			return
		}
		// Imprime a mensagem recebida do cliente no console do servidor.
		fmt.Printf("Mensagem recebida do cliente: %s", message)
		// Envia a mensagem de volta ao cliente.
		conn.Write([]byte(message))
	}
}

// main é o ponto de entrada do programa.
// - Cria um listener TCP que aguarda conexões de clientes na porta especificada.
// - Aceita conexões de clientes e cria uma goroutine para tratar cada conexão simultaneamente.
func main() {
	// Cria um listener TCP na porta especificada.
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		// Se houver um erro ao iniciar o servidor, imprime a mensagem de erro e sai do programa.
		fmt.Println("Erro ao iniciar o servidor:", err)
		os.Exit(1)
	}
	defer listener.Close() // Garante que o listener seja fechado ao final da execução.

	fmt.Printf("Servidor de Eco TCP esperando por conexões na porta %s...\n", PORT)

	for {
		// Aceita uma conexão de cliente.
		conn, err := listener.Accept()
		if err != nil {
			// Se houver um erro ao aceitar a conexão, imprime a mensagem de erro e continua aguardando novas conexões.
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		// Cria uma goroutine para tratar a conexão do cliente de forma concorrente.
		go handleConnection(conn)
	}
}
