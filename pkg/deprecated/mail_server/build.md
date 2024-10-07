## 1. Salve o Código-fonte: Copie o código em C para um arquivo com a extensão .c, por exemplo, smtp_server.c

## 2. Compilação do Código: Abra um terminal e navegue até o diretório onde o arquivo .c está localizado. Em seguida, use um compilador C, como o gcc, para compilar o código. Por exemplo

gcc -o smtp_server smtp_server.c

## 3. Execução do Programa: Após a compilação bem-sucedida, você pode executar o programa compilado. Basta digitar o nome do arquivo executável no terminal e pressionar Enter. Por exemplo

./smtp_server

## 4. Teste do Servidor: Para testar o servidor, você pode usar um cliente SMTP, como o telnet, para estabelecer uma conexão com o servidor e enviar comandos SMTP. Por exemplo

telnet localhost 25
