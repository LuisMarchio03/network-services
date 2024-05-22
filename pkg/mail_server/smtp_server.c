// Importação das bibliotecas necessárias
// <stdio.h>: Biblioteca padrão de entrada e saída em C, usada para funções como printf e perror.
// <stdlib.h>: Biblioteca padrão para funções utilitárias, como malloc, free e exit.
// <string.h>: Biblioteca para manipulação de strings, como strlen, strcpy, strcmp, etc.
// <unistd.h>: Biblioteca para chamadas do sistema POSIX, como close e read.
// <arpa/inet.h>: Biblioteca para operações de rede, incluindo criação e manipulação de sockets.
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>

// Definição de constantes
// PORT: Define a porta que o servidor SMTP irá utilizar para escutar conexões. Neste caso, é a porta 25, que é a porta padrão para SMTP.
// BUFFER_SIZE: Define o tamanho do buffer utilizado para armazenar dados recebidos e enviados. Neste caso, é 1024 bytes.
#define PORT 25
#define BUFFER_SIZE 1024

// Definição da estrutura SMTPServer
// A estrutura SMTPServer armazena informações sobre o servidor SMTP.
// socket: Um descritor de arquivo para o socket do servidor. Sockets são usados para comunicação entre o servidor e os clientes.
// address: Uma estrutura sockaddr_in que contém o endereço e a porta do servidor. Esta estrutura é usada para configurar o socket do servidor e associá-lo a um endereço e porta específicos.
typedef struct
{
    int socket;
    struct sockaddr_in address;
} SMTPServer;

void initialize_server(SMTPServer *server)
{
    // Cria um socket do tipo stream (TCP)
    // AF_INET: Família de endereços para IPv4
    // SOCK_STREAM: Tipo de socket para TCP
    // 0: Protocolo padrão para SOCK_STREAM (TCP)
    server->socket = socket(AF_INET, SOCK_STREAM, 0);
    if (server->socket < 0)
    {
        // Exibe mensagem de erro e termina o programa se o socket não puder ser criado
        perror("Failed to create socket");
        exit(EXIT_FAILURE);
    }

    // Configura o endereço do servidor
    // sin_family: Família de endereços (IPv4)
    // sin_addr.s_addr: Endereço IP (INADDR_ANY permite receber conexões em qualquer endereço de rede do host)
    // sin_port: Porta do servidor, convertida para o formato de rede com htons
    server->address.sin_family = AF_INET;
    server->address.sin_addr.s_addr = INADDR_ANY;
    server->address.sin_port = htons(PORT);

    // Associa o socket a um endereço e porta
    // bind: Vincula o socket ao endereço e porta especificados
    if (bind(server->socket, (struct sockaddr *)&server->address, sizeof(server->address)) < 0)
    {
        // Exibe mensagem de erro, fecha o socket e termina o programa se a associação falhar
        perror("Failed to bind socket");
        close(server->socket);
        exit(EXIT_FAILURE);
    }

    // Coloca o socket em modo de escuta para aceitar conexões
    // listen: Habilita o socket a aceitar conexões
    // 5: Tamanho da fila de conexões pendentes
    if (listen(server->socket, 5) < 0)
    {
        // Exibe mensagem de erro, fecha o socket e termina o programa se a escuta falhar
        perror("Failed to listen on socket");
        close(server->socket);
        exit(EXIT_FAILURE);
    }

    // Exibe uma mensagem indicando que o servidor foi inicializado com sucesso e está escutando na porta especificada
    printf("SMTP server initialized and listening on port %d\n", PORT);
}

// Função para processar conexões de clientes
// Esta função recebe comandos dos clientes, processa-os e envia respostas apropriadas.
// Implementa os comandos básicos do SMTP: HELO, MAIL FROM, RCPT TO, DATA e QUIT.

void process_connection(int client_socket)
{
    char buffer[BUFFER_SIZE]; // Buffer para armazenar dados recebidos do cliente
    int bytes_read;           // Variável para armazenar o número de bytes lidos do socket

    // Envia uma mensagem de boas-vindas ao cliente
    // "220" é o código de status para um serviço pronto no protocolo SMTP
    send(client_socket, "220 Simple SMTP Server\r\n", 24, 0);

    // Loop para processar os comandos recebidos do cliente
    while ((bytes_read = recv(client_socket, buffer, BUFFER_SIZE, 0)) > 0)
    {
        buffer[bytes_read] = '\0';      // Adiciona um terminador nulo ao final dos dados recebidos
        printf("Received: %s", buffer); // Exibe os dados recebidos no console

        // Verifica se o comando recebido é "HELO" ou "EHLO"
        if (strncmp(buffer, "HELO", 4) == 0 || strncmp(buffer, "EHLO", 4) == 0)
        {
            // Envia uma resposta "250 Hello" ao cliente
            // "250" é o código de status para uma solicitação bem-sucedida no protocolo SMTP
            send(client_socket, "250 Hello\r\n", 11, 0);
        }
        // Verifica se o comando recebido é "MAIL FROM:"
        else if (strncmp(buffer, "MAIL FROM:", 10) == 0)
        {
            // Envia uma resposta "250 OK" ao cliente
            send(client_socket, "250 OK\r\n", 7, 0);
        }
        // Verifica se o comando recebido é "RCPT TO:"
        else if (strncmp(buffer, "RCPT TO:", 8) == 0)
        {
            // Envia uma resposta "250 OK" ao cliente
            send(client_socket, "250 OK\r\n", 7, 0);
        }
        // Verifica se o comando recebido é "DATA"
        else if (strncmp(buffer, "DATA", 4) == 0)
        {
            // Envia uma resposta "354 End data with <CR><LF>.<CR><LF>" ao cliente
            // "354" é o código de status indicando que o servidor está pronto para receber a mensagem
            send(client_socket, "354 End data with <CR><LF>.<CR><LF>\r\n", 36, 0);
            // Lê os dados da mensagem do cliente
            bytes_read = recv(client_socket, buffer, BUFFER_SIZE, 0);
            buffer[bytes_read] = '\0';        // Adiciona um terminador nulo ao final dos dados recebidos
            printf("Email data: %s", buffer); // Exibe os dados do e-mail no console
            // Envia uma resposta "250 OK" ao cliente após receber os dados do e-mail
            send(client_socket, "250 OK\r\n", 7, 0);
        }
        // Verifica se o comando recebido é "QUIT"
        else if (strncmp(buffer, "QUIT", 4) == 0)
        {
            // Envia uma resposta "221 Bye" ao cliente
            // "221" é o código de status para encerramento da conexão no protocolo SMTP
            send(client_socket, "221 Bye\r\n", 9, 0);
            break; // Sai do loop e encerra a conexão
        }
        // Se o comando recebido não for reconhecido
        else
        {
            // Envia uma resposta "500 Unrecognized command" ao cliente
            // "500" é o código de status para erro de sintaxe no protocolo SMTP
            send(client_socket, "500 Unrecognized command\r\n", 26, 0);
        }
    }

    // Fecha o socket do cliente após processar a conexão
    close(client_socket);
}

// Função principal do programa
int main()
{
    SMTPServer server;          // Declaração de uma variável do tipo SMTPServer
    initialize_server(&server); // Inicializa o servidor SMTP

    // Loop infinito para aceitar e processar conexões de clientes
    while (1)
    {
        // Aceita uma conexão de cliente
        int client_socket = accept(server.socket, NULL, NULL);
        if (client_socket < 0)
        {
            // Exibe uma mensagem de erro se a conexão falhar
            perror("Failed to accept connection");
            continue; // Continua para a próxima iteração do loop
        }

        // Processa a conexão do cliente
        process_connection(client_socket);
    }

    // Fecha o socket do servidor após sair do loop (embora este ponto nunca seja alcançado devido ao loop infinito)
    close(server.socket);
    return 0; // Retorna 0 para indicar que o programa terminou corretamente
}
