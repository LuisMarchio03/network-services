package ftp

import (
	"FTP_Server/internal/auth"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const rootDir = "/data"

func StartFTPServer() {
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		log.Fatalf("Erro ao criar diretório FTP: %v", err)
	}
	listener, err := net.Listen("tcp", ":21")
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor FTP: %v", err)
	}
	defer listener.Close()
	fmt.Println("Servidor FTP (didático) rodando na porta 21, servindo", rootDir)

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
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
	reply := func(s string) { w.WriteString(s + "\r\n"); w.Flush() }

	reply("220 Bem-vindo ao servidor FTP didatico")

	var authed bool
	var pendingUser string
	var dataListener net.Listener
	defer func() {
		if dataListener != nil {
			dataListener.Close()
		}
	}()

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		parts := strings.SplitN(line, " ", 2)
		cmd := strings.ToUpper(parts[0])
		arg := ""
		if len(parts) > 1 {
			arg = parts[1]
		}

		switch cmd {
		case "USER":
			pendingUser = arg
			reply("331 Usuario OK, envie a senha")
		case "PASS":
			if auth.Authenticate(pendingUser, arg) {
				authed = true
				reply("230 Login OK")
			} else {
				reply("530 Login incorreto")
			}
		case "SYST":
			reply("215 UNIX Type: L8")
		case "FEAT":
			reply("211 End")
		case "TYPE":
			reply("200 Type set")
		case "PWD":
			reply(`257 "/" is current directory`)
		case "CWD":
			reply("250 Diretorio alterado")
		case "PASV":
			if !authed {
				reply("530 Faca login primeiro")
				continue
			}
			if dataListener != nil {
				dataListener.Close()
			}
			dataListener, err = net.Listen("tcp", ":0")
			if err != nil {
				reply("425 Nao foi possivel abrir conexao de dados")
				continue
			}
			port := dataListener.Addr().(*net.TCPAddr).Port
			// Anuncia 127,0,0,1 e a porta em p1*256+p2 (cliente local no demo).
			p1, p2 := port/256, port%256
			reply(fmt.Sprintf("227 Entering Passive Mode (127,0,0,1,%d,%d)", p1, p2))
		case "LIST":
			if dataListener == nil {
				reply("425 Use PASV primeiro")
				continue
			}
			reply("150 Abrindo conexao de dados para listagem")
			dataConn, err := dataListener.Accept()
			if err != nil {
				reply("425 Falha na conexao de dados")
				continue
			}
			entries, _ := os.ReadDir(rootDir)
			for _, e := range entries {
				fmt.Fprintf(dataConn, "%s\r\n", e.Name())
			}
			dataConn.Close()
			reply("226 Listagem concluida")
		case "RETR":
			if dataListener == nil {
				reply("425 Use PASV primeiro")
				continue
			}
			f, err := os.Open(filepath.Join(rootDir, filepath.Base(arg)))
			if err != nil {
				reply("550 Arquivo nao encontrado")
				continue
			}
			reply("150 Abrindo conexao de dados para download")
			dataConn, derr := dataListener.Accept()
			if derr != nil {
				f.Close()
				reply("425 Falha na conexao de dados")
				continue
			}
			io.Copy(dataConn, f)
			dataConn.Close()
			f.Close()
			reply("226 Transferencia concluida")
		case "STOR":
			if dataListener == nil {
				reply("425 Use PASV primeiro")
				continue
			}
			reply("150 Abrindo conexao de dados para upload")
			dataConn, derr := dataListener.Accept()
			if derr != nil {
				reply("425 Falha na conexao de dados")
				continue
			}
			f, err := os.Create(filepath.Join(rootDir, filepath.Base(arg)))
			if err != nil {
				dataConn.Close()
				reply("550 Nao foi possivel criar o arquivo")
				continue
			}
			io.Copy(f, dataConn)
			f.Close()
			dataConn.Close()
			reply("226 Upload concluido")
		case "NOOP":
			reply("200 OK")
		case "QUIT":
			reply("221 Ate logo")
			return
		default:
			reply("502 Comando nao implementado")
		}
	}
}
