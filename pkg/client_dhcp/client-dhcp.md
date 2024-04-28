# Documentação: Cliente DHCP em Go

## 1. Introdução

Este documento descreve a implementação de um cliente DHCP em Go, responsável por enviar uma mensagem DHCP Discover para um servidor DHCP local e receber uma resposta.

## 2. Funcionamento do Cliente DHCP

O cliente DHCP funciona da seguinte maneira:

1. **Conexão com o Servidor DHCP**
    - Estabelece uma conexão UDP com o servidor DHCP na porta 67 (padrão para DHCP).

2. **Construção da Mensagem DHCP Discover**
    - Define os parâmetros necessários para construir a mensagem DHCP Discover:
        - `messageType`: Tipo 1 (DHCP Discover).
        - `hardwareType`: Tipo 1 (Ethernet).
        - `hardwareLength`: Comprimento do endereço MAC (6 bytes para Ethernet).
        - `xid`: Número de transação DHCP.
        - `clientIP`, `yourIP`, `serverIP`, `gatewayIP`: Endereços IP utilizados na negociação DHCP.
        - `clientHardwareAddr`: Endereço MAC do cliente.
    - Utiliza um buffer de bytes (`bytes.Buffer`) para construir a mensagem DHCP.

3. **Envio da Mensagem DHCP para o Servidor**
    - Envia a mensagem DHCP Discover para o servidor DHCP utilizando a conexão estabelecida.

4. **Recepção da Resposta do Servidor**
    - Aguarda e recebe a resposta do servidor DHCP.
    - Extrai o endereço IP oferecido (`yourIP`) da resposta recebida.

## 3. Detalhes do Código

O código é estruturado da seguinte forma:

- **`main()`**: Função principal que realiza as etapas descritas acima.
    - Inicia uma conexão UDP com o servidor DHCP.
    - Constrói a mensagem DHCP Discover.
    - Envia a mensagem para o servidor e aguarda a resposta.
    - Extrai e exibe o endereço IP oferecido pela resposta DHCP.

- **Variáveis e Estruturas Utilizadas**:
    - `conn`: Representa a conexão UDP com o servidor DHCP.
    - `buf`: Buffer de bytes utilizado para construir a mensagem DHCP.
    - `messageType`, `hardwareType`, `xid`, `clientHardwareAddr`, etc.: Parâmetros da mensagem DHCP.

## 4. Uso e Execução

Para executar o cliente DHCP:

1. Certifique-se de ter um servidor DHCP em execução na máquina local ou no endereço especificado (`localhost:67`).
2. Execute o código Go fornecido.
3. O cliente enviará a mensagem DHCP Discover e exibirá o endereço IP oferecido pelo servidor DHCP.

## 5. Considerações Finais

Este cliente DHCP demonstra como construir e enviar uma mensagem DHCP Discover utilizando a biblioteca padrão do Go. Ele é útil para testar e interagir com servidores DHCP durante o desenvolvimento e depuração de redes.

---

Este documento oferece uma visão detalhada do funcionamento e implementação do cliente DHCP em Go, destacando cada etapa do processo e explicando as funcionalidades específicas utilizadas no código.
