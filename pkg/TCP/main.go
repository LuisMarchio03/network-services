package main

import (
	"fmt"
	"net"
	"syscall"
)

const (
	SOCK_STREAM = 1   // Socket de fluxo (TCP)
	IPPROTO_TCP = 6   // Protocolo TCP
	SOMAXCONN   = 128 // Tamanho máximo da fila de conexões pendentes
)

func main() {
	// Criar um socket TCP/IP
	sockfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, IPPROTO_TCP)
	if err != nil {
		fmt.Printf("Erro ao criar socket: %v\n", err)
		return
	}
	defer syscall.Close(sockfd)

	// Configurar estrutura de endereço
	var addr syscall.SockaddrInet4
	addr.Port = 8888                 // Porta a ser escutada
	copy(addr.Addr[:], net.IPv4zero) // Endereço IP (0.0.0.0 para escutar em todas as interfaces)

	// Vincular o socket ao endereço
	err = syscall.Bind(sockfd, &addr)
	if err != nil {
		fmt.Printf("Erro ao vincular o socket: %v\n", err)
		return
	}

	// Iniciar a escuta do socket
	err = syscall.Listen(sockfd, SOMAXCONN)
	if err != nil {
		fmt.Printf("Erro ao iniciar a escuta: %v\n", err)
		return
	}

	fmt.Println("Servidor TCP iniciado e escutando na porta 8888")

	// Loop infinito para aguardar por conexões
	for {
		// Aceitar uma nova conexão
		clientfd, _, err := syscall.Accept(sockfd)
		if err != nil {
			fmt.Printf("Erro ao aceitar conexão: %v\n", err)
			continue
		}

		fmt.Println("Cliente conectado")

		// Tratar a conexão em uma goroutine
		go handleConnection(clientfd)
	}

	// O código abaixo nunca será alcançado no loop infinito acima
}

func handleConnection(clientfd int) {
	defer syscall.Close(clientfd)

	// Buffer para armazenar dados recebidos
	buffer := make([]byte, 1024)

	// Loop para receber e enviar dados de volta ao cliente
	for {
		// Recebe dados do cliente
		bytesReceived, err := syscall.Read(clientfd, buffer)
		if err != nil {
			fmt.Printf("Erro ao receber dados: %v\n", err)
			break
		}

		// Ecoa os dados de volta para o cliente
		_, err = syscall.Write(clientfd, buffer[:bytesReceived])
		if err != nil {
			fmt.Printf("Erro ao enviar dados de volta ao cliente: %v\n", err)
			break
		}
	}
}
