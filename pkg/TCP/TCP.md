# Parte 1: Fundamentos de Redes de Computadores

## Visão Geral dos Princípios Básicos de Redes de Computadores

### Modelos de Referência

- **OSI Model (Modelo OSI)**: Um modelo em 7 camadas que serve como referência para a comunicação entre sistemas de redes.
  - Camadas: Física, Enlace de Dados, Rede, Transporte, Sessão, Apresentação, Aplicação.
- **TCP/IP Model (Modelo TCP/IP)**: Um modelo de 4 camadas usado na internet.
  - Camadas: Acesso à Rede, Internet, Transporte, Aplicação.

### Protocolos de Rede

- **TCP (Transmission Control Protocol)**: Protocolo de transporte confiável, orientado à conexão.
- **UDP (User Datagram Protocol)**: Protocolo de transporte não confiável, sem conexão.
- **HTTP/HTTPS**: Protocolos de aplicação para transferência de hipertexto.
- **DNS (Domain Name System)**: Serviço de resolução de nomes.
- **DHCP (Dynamic Host Configuration Protocol)**: Protocolo para atribuição dinâmica de endereços IP.

## Aplicação Prática: Implementação de Servidor TCP

### Título do Curso: "Construindo e Compreendendo Redes - Parte 1"

### Conteúdo

- Visão geral dos princípios básicos de redes de computadores.
- Implementação de um servidor TCP simples.
- Demonstração de comunicação básica cliente-servidor.
- Uso de sockets para comunicação de rede.

### Atividade Prática

- Desenvolvimento de um servidor de eco TCP.
- Exercícios para entender o fluxo de dados entre cliente e servidor.

## Desenvolvimento de um Servidor de Eco TCP

### Implementação em C

Vamos implementar um servidor de eco TCP em C. Este servidor vai receber mensagens de um cliente e enviar a mesma mensagem de volta (eco).

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>

#define PORT 8888
#define BUFFER_SIZE 1024

int main() {
    int server_socket, client_socket;
    struct sockaddr_in server_addr, client_addr;
    char buffer[BUFFER_SIZE];
    socklen_t client_addr_len;

    // Criar socket TCP
    server_socket = socket(AF_INET, SOCK_STREAM, 0);
    if (server_socket == -1) {
        perror("Erro ao criar socket");
        exit(EXIT_FAILURE);
    }

    // Configurar endereço do servidor
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(PORT);

    // Vincular o socket ao endereço e à porta
    if (bind(server_socket, (struct sockaddr *)&server_addr, sizeof(server_addr)) == -1) {
        perror("Erro ao vincular o socket");
        close(server_socket);
        exit(EXIT_FAILURE);
    }

    // Aguardar por conexões de clientes
    if (listen(server_socket, 5) == -1) {
        perror("Erro ao aguardar por conexões");
        close(server_socket);
        exit(EXIT_FAILURE);
    }

    printf("Servidor de Eco TCP esperando por conexões na porta %d...\n", PORT);

    while (1) {
        client_addr_len = sizeof(client_addr);
        client_socket = accept(server_socket, (struct sockaddr *)&client_addr, &client_addr_len);
        if (client_socket == -1) {
            perror("Erro ao aceitar conexão");
            continue;
        }

        printf("Cliente conectado: %s:%d\n", inet_ntoa(client_addr.sin_addr), ntohs(client_addr.sin_port));

        ssize_t bytes_received;
        while ((bytes_received = recv(client_socket, buffer, BUFFER_SIZE, 0)) > 0) {
            buffer[bytes_received] = '\0';
            printf("Mensagem recebida do cliente: %s\n", buffer);
            send(client_socket, buffer, bytes_received, 0);
        }

        if (bytes_received == 0) {
            printf("Cliente desconectado\n");
        } else if (bytes_received == -1) {
            perror("Erro ao receber dados do cliente");
        }

        close(client_socket);
    }

    close(server_socket);
    return 0;
}
```

### Implementação em Go

Agora, vamos implementar um servidor de eco TCP em Go. Este servidor também vai receber mensagens de um cliente e enviar a mesma mensagem de volta (eco).

```go
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

func handleConnection(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Cliente desconectado.")
            return
        }
        fmt.Printf("Mensagem recebida do cliente: %s", message)
        conn.Write([]byte(message))
    }
}

func main() {
    listener, err := net.Listen("tcp", PORT)
    if err != nil {
        fmt.Println("Erro ao iniciar o servidor:", err)
        os.Exit(1)
    }
    defer listener.Close()

    fmt.Printf("Servidor de Eco TCP esperando por conexões na porta %s...\n", PORT)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Erro ao aceitar conexão:", err)
            continue
        }

        go handleConnection(conn)
    }
}

```

## Explicação das Implementações

## Em C

### Criação do Socket

`socket(AF_INET, SOCK_STREAM, 0)` cria um socket TCP.

### Configuração do Endereço do Servidor

`server_addr` é configurado com `AF_INET`, `INADDR_ANY` para aceitar conexões de qualquer endereço, e a porta definida.

### Vinculação e Escuta

- `bind` vincula o socket ao endereço e porta.
- `listen` faz o socket escutar por conexões de clientes.

### Aceitação de Conexões e Comunicação

- `accept` aceita conexões de clientes.
- `recv` recebe dados do cliente.
- `send` envia dados de volta ao cliente.

### Encerramento

- `close` fecha o socket do cliente e do servidor.

## Em Go

### Criação do Listener

`net.Listen("tcp", PORT)` cria um listener TCP na porta especificada.

### Aceitação de Conexões

`listener.Accept()` aceita conexões de clientes.

### Tratamento da Conexão

`handleConnection` lê dados do cliente com `bufio.NewReader(conn).ReadString('\n')` e envia de volta com `conn.Write([]byte(message))`.

### Concorrência

`go handleConnection(conn)` usa goroutines para tratar múltiplas conexões simultaneamente.

## Conclusão

Essas implementações de servidores de eco TCP em C e Go fornecem uma base sólida para aprender sobre comunicação cliente-servidor, sockets e programação de rede. Ao entender e experimentar esses exemplos, você estará bem preparado para expandir seus conhecimentos e habilidades em redes de computadores e desenvolvimento de servidores.
