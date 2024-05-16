# Testando o Servidor de Eco TCP em C

## Passos para Compilar e Executar o Servidor em C

### Salve o Código em um Arquivo

Salve o código C em um arquivo chamado `echo_server.c`.

### Compile o Código

Abra um terminal e navegue até o diretório onde o arquivo `echo_server.c` está localizado. Compile o código usando `gcc`:

```bash
gcc -o echo_server echo_server.c
```

### Execute o Servidor

Execute o servidor compilado:

```bash
./echo_server
```

## Usando Telnet para Testar o Servidor em C

**Abra um Novo Terminal:**

- Abra um novo terminal enquanto o servidor está em execução.

**Conecte ao Servidor Usando Telnet:**

- No novo terminal, use o comando telnet para conectar ao servidor:

```bash
telnet localhost 8888
```

**Envie Mensagens:**

- Digite uma mensagem e pressione Enter. Você deve ver a mesma mensagem ecoada de volta pelo servidor.

**Saia do Telnet:**

- Para sair da sessão Telnet, você pode usar Ctrl+] para entrar no modo de comando Telnet e depois digitar quit.

# Testando o Servidor de Eco TCP em Go

## Passos para Executar o Servidor em Go

**Salve o Código em um Arquivo:**

Salve o código Go em um arquivo chamado echo_server.go.

**Execute o Servidor:**

Abra um terminal e navegue até o diretório onde o arquivo echo_server.go está localizado. Execute o servidor usando o comando go run:

```bash
go run echo_server.go
```

## Usando Telnet para Testar o Servidor em Go

**Abra um Novo Terminal:**

- Abra um novo terminal enquanto o servidor está em execução.

**Conecte ao Servidor Usando Telnet:**

- No novo terminal, use o comando telnet para conectar ao servidor:

```bash
telnet localhost 8888
```

### Envie Mensagens

- Digite uma mensagem e pressione Enter. Você deve ver a mesma mensagem ecoada de volta pelo servidor.

### Saia do Telnet

- Para sair da sessão Telnet, você pode usar Ctrl+] para entrar no modo de comando Telnet e depois digitar quit.

# Exemplo de Sessão Telnet

## Terminal do Servidor (C ou Go)

```bash
Servidor de Eco TCP esperando por conexões na porta 8888...
```

## Terminal do Servidor (C ou Go)

```bash
$ telnet localhost 8888
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
Hello, Server!
Hello, Server!
```
