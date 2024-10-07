package ftp

import (
	"FTP_Server/internal/auth"
	"fmt"
	"log"
	"net"
)

func StartFTPServer() {
	listener, err := net.Listen("tcp", ":21")
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor FTP: %v", err)
	}
	defer listener.Close()

	fmt.Println("Servidor FTP rodando na porta 21...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Fprintln(conn, "220 Bem-vindo ao servidor FTP")

	username, password := "", ""
	fmt.Fprint(conn, "331 Nome de usuário OK, aguarde a senha.\n")
	fmt.Fscanln(conn, &username)
	fmt.Fprint(conn, "331 Senha necessária.\n")
	fmt.Fscanln(conn, &password)

	if auth.Authenticate(username, password) {
		fmt.Fprint(conn, "230 Login OK.\n")
	} else {
		fmt.Fprint(conn, "530 Login incorreto.\n")
	}
}
