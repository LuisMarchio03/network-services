# Documento de Desenvolvimento e Estudo: Servidor DHCP

## 1. Introdução

Este documento descreve o processo de desenvolvimento e estudo de um servidor DHCP (Dynamic Host Configuration Protocol), utilizado para atribuir configurações de rede dinâmicas aos clientes. O servidor DHCP é implementado em Go (Golang) e utiliza o MongoDB como banco de dados para gerenciar os endereços IP disponíveis.

## 2. Objetivos

O objetivo deste projeto é criar um servidor DHCP funcional capaz de:

- Receber mensagens DHCP de clientes.
- Atribuir endereços IP disponíveis aos clientes.
- Responder a solicitações DHCP com ofertas de configuração.
- Gerenciar o status de atribuição de endereços IP no MongoDB.

## 3. Tecnologias Utilizadas

As principais tecnologias utilizadas neste projeto são:

- **Go (Golang)**: Linguagem de programação principal para implementação do servidor DHCP.
- **MongoDB**: Banco de dados NoSQL utilizado para armazenar registros de endereços IP.
- **Context**: Pacote padrão do Go utilizado para gerenciamento de contexto e cancelamento de operações.
- **UDP (User Datagram Protocol)**: Protocolo de transporte utilizado para comunicação entre o servidor DHCP e os clientes.

## 4. Implementação

### 4.1 Estrutura do Projeto

O projeto é estruturado da seguinte forma:

- **`main.go`**: Arquivo principal que inicia o servidor DHCP.
- **`dhcp/`**: Pacote contendo a lógica do servidor DHCP.
- **`mongodb/`**: Pacote para interação com o MongoDB.
- **`schemas/`**: Definições de estrutura para os dados no MongoDB.

### 4.2 Componentes Principais

#### 4.2.1 `main.go`

- Inicializa o logger e a conexão com o MongoDB.
- Cria uma nova instância do servidor DHCP e inicia a escuta por mensagens DHCP.

#### 4.2.2 `dhcp/`

- **`server.go`**: Implementação do servidor DHCP que recebe e processa mensagens DHCP.
- **`message.go`**: Lógica para processar mensagens DHCP e gerar respostas.

#### 4.2.3 `mongodb/`

- **`mongo.go`**: Conecta-se ao MongoDB e fornece métodos para interação com o banco de dados.

#### 4.2.4 `schemas/`

- **`ip_record.go`**: Define a estrutura de dados para registros de endereços IP no MongoDB.

## 5. Processo de Desenvolvimento

O desenvolvimento do servidor DHCP seguiu as seguintes etapas:

1. **Planejamento**: Definição dos requisitos e estrutura básica do servidor DHCP.
2. **Implementação**: Codificação das funcionalidades principais do servidor em Go.
3. **Testes Unitários**: Criação de testes para garantir a correção e robustez do código.
4. **Integração com MongoDB**: Implementação da lógica para interagir com o banco de dados.
5. **Depuração e Otimização**: Identificação e correção de erros, bem como otimização de desempenho.
6. **Documentação**: Elaboração de documentação técnica para facilitar o entendimento e manutenção do código.

## 6. Estudo e Aprendizado

Durante o desenvolvimento deste projeto, foram adquiridos os seguintes conhecimentos e habilidades:

- **Go (Golang)**: Aprofundamento na linguagem Go, incluindo conceitos avançados de concorrência e gerenciamento de pacotes.
- **MongoDB**: Aprendizado sobre o uso do MongoDB com o driver oficial para Go.
- **Protocolos de Rede**: Compreensão dos princípios e funcionamento do protocolo DHCP e do UDP.

## 7. Conclusão

O servidor DHCP desenvolvido demonstra a aplicação prática dos conceitos estudados e a habilidade de implementar soluções de rede usando tecnologias modernas. Este projeto serviu como uma oportunidade de aprimorar habilidades de programação em Go e expandir o conhecimento em sistemas distribuídos e gerenciamento de dados.
