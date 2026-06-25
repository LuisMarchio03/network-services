[Português](README.md) | [English](README.en.md)

# network-services

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose%20v2-2496ED?logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

Projeto didático que recria serviços de rede (DHCP, DNS, NTP, FTP, TCP echo, pipeline de logs) do zero em Go/C, executável como demo local reproduzível via Docker Compose.

---

## Arquitetura

```mermaid
flowchart LR
  seed[Seed Node] --> mongo[(MongoDB)]
  dhcpc[DHCP client] -- UDP --> dhcps[DHCP server] --> mongo
  ftp[FTP server] --> mongo
  prod[logs-producer] -- AMQP --> rabbit[(RabbitMQ)]
  rabbit -- AMQP --> recv[logs-receiver] -- HTTP --> es[(Elasticsearch)] --> kib[Kibana]
  dnsc[DNS client] -- UDP --> dns[DNS server]
  ntpc[NTP client] -- UDP --> ntp[NTP server]
  echo[TCP echo]
```

---

## Matriz de Serviços

| Serviço | Linguagem | Porta interna | Porta host | Status |
|---|---|---|---|---|
| MongoDB | — | 27017 | 27017 | infra |
| RabbitMQ (+mgmt) | — | 5672 / 15672 | 5672 / 15672 | infra |
| Elasticsearch | — | 9200 | 9200 | infra |
| Kibana | — | 5601 | 5601 | infra |
| seed (Node) | JS | — | — | one-shot (seed MongoDB) |
| logs-producer | Go | — | — | ✅ publica logs via AMQP |
| logs-receiver | Go | — | — | ✅ consome AMQP → Elasticsearch |
| dhcp-server | Go | UDP 67 | — (interno) | ✅ aloca IP do MongoDB |
| dhcp-client | Go | — | — | ✅ solicita IP via DHCP |
| dns-server | Go | UDP 53 | — (interno) | 🧪 didático (zones.json) |
| dns-client | Go | — | — | 🧪 consulta DNS didático |
| ntp-server | Go | UDP 123 | — (interno) | ✅ SNTP (timestamps era-1900) |
| ntp-client | Go | — | — | ✅ consulta NTP |
| tcp-echo | Go | 8888 | 8888 | ✅ servidor de eco TCP |
| ftp-server | Go | 21 | 2121 | 🧪 controle + PASV LIST/RETR/STOR |
| network-quiz | Go | 8080 | 8080 | 🧪 quiz TCP de redes |
| session-manager | C | — | — | ⚙️ opcional (profile `extras`) |

> **Legenda:** ✅ funcionando · 🧪 didático/experimental · ⚙️ opcional

---

## Quickstart

### Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) + Docker Compose v2
- `make`

### Subir o stack

```bash
make up
```

Copia `.env` de `.env.example` automaticamente se não existir, depois executa:

```bash
docker compose -f StartNewServer/docker-compose.yml up -d --build
```

### Demo completo

```bash
make demo
```

O script `scripts/demo.sh` sobe o stack, exercita DNS, NTP, TCP echo, exibe a contagem de logs no Elasticsearch e acompanha os logs do DHCP.

### Outros alvos `make`

| Alvo | O que faz |
|---|---|
| `make down` | Para e remove os containers |
| `make logs` | Acompanha logs dos containers (últimas 50 linhas) |
| `make test` | Executa testes Go em `pkg/dhcp_server` e `pkg/NTP` |
| `make clean` | `down -v` (remove volumes) |
| `make dev` | Orientação para rodar serviços nativos com `go run` |

### Links úteis após `make up`

| Interface | URL |
|---|---|
| Kibana | http://localhost:5601 |
| RabbitMQ Management | http://localhost:15672 (guest/guest) |
| Elasticsearch | http://localhost:9200 |

---

## Uso por Serviço

### Pipeline de Logs (logs-producer → RabbitMQ → logs-receiver → Elasticsearch → Kibana)

```bash
# Ver logs chegando no Elasticsearch
curl -s http://localhost:9200/logs/_count | python3 -m json.tool
# {"count": 42, ...}

# Verificar logs dos containers
make logs
```

### DHCP

```bash
# Acompanhar alocações de IP (servidor interno)
docker compose -f StartNewServer/docker-compose.yml logs dhcp-server

# Exemplo de saída
# dhcp-server  | Received DHCP Discover from 02:42:ac:12:00:05
# dhcp-server  | Offering IP 10.0.0.10 to 02:42:ac:12:00:05
# dhcp-server  | ACK sent: 10.0.0.10 assigned
```

### DNS (profile clients)

```bash
docker compose -f StartNewServer/docker-compose.yml --profile clients run --rm dns-client

# Exemplo de saída
# Querying DNS for example.local...
# Response: example.local -> 10.0.0.1
```

### NTP (profile clients)

```bash
docker compose -f StartNewServer/docker-compose.yml --profile clients run --rm ntp-client

# Exemplo de saída
# NTP response received
# Server time: 2024-01-15 14:32:10 UTC
# Offset: +0.003s
```

### TCP Echo

```bash
# Conectar ao servidor de eco
nc localhost 8888
Hello network-services!
# Hello network-services!   <- eco do servidor
```

### FTP

```bash
ftp -p localhost 2121
# Connected to localhost
# 220 Go FTP Server ready
# Name: anonymous
# 331 Password required
# Password: (enter)
# 230 Login successful
ftp> ls
ftp> get arquivo.txt
ftp> quit
```

### Network Quiz

```bash
# Acessar o quiz de redes via TCP
nc localhost 8080
# Bem-vindo ao Quiz de Redes!
# Pergunta 1: O que é um endereço MAC?
```

---

## Estrutura do Projeto

```
network-services/
├── Makefile                        # Alvos: up, down, demo, logs, test, clean
├── .env.example                    # Variáveis de ambiente de referência
├── scripts/
│   ├── demo.sh                     # Demo end-to-end automatizado
│   └── smoke.sh                    # Smoke tests básicos
├── pkg/                            # Implementações Go dos serviços
│   ├── dhcp_server/                # Servidor DHCP (aloca IPs do MongoDB)
│   │   ├── cmd/                    # Entrypoint
│   │   ├── internal/               # Lógica interna
│   │   ├── mongo/                  # Integração MongoDB
│   │   └── go.mod
│   ├── dhcp-simulation-client/     # Cliente DHCP
│   ├── dhcp-simulation-server/     # Servidor DHCP (versão simulação)
│   ├── dns-simulation-server/      # Servidor DNS didático
│   ├── dns-simulation-cliente/     # Cliente DNS
│   ├── NTP/                        # Servidor e cliente NTP (SNTP)
│   │   ├── cmd/                    # Entrypoints server/client
│   │   ├── ntptime/                # Conversão de timestamps era-1900
│   │   └── go.mod
│   ├── FTP/                        # Servidor FTP (controle + PASV)
│   ├── TCP/                        # Servidor de eco TCP
│   ├── network-quiz/               # Quiz TCP de redes
│   ├── logs-producer-teste/        # Produtor de logs via AMQP
│   ├── logs-receiver/              # Receptor de logs → Elasticsearch
│   ├── session-manager-c/          # Gerenciador de sessões em C (opcional)
│   └── deprecated/                 # Código legado não utilizado
└── StartNewServer/                 # Orquestração Docker Compose
    ├── docker-compose.yml          # Definição de todos os serviços
    ├── src/                        # Seed Node.js (seed MongoDB)
    │   ├── setup.js
    │   └── rabbitmq-elasticsearch.js
    ├── elasticsearch.yml           # Configuração do Elasticsearch
    └── Dockerfile                  # Imagem base
```

---

## Variáveis de Ambiente

Copie `.env.example` para `StartNewServer/.env` (feito automaticamente pelo `make up`):

| Variável | Padrão | Descrição |
|---|---|---|
| `MONGO_URI` | `mongodb://mongodb:27017` | URI de conexão MongoDB |
| `MONGO_DB` | `mydb` | Nome do banco de dados |
| `RABBITMQ_URL` | `amqp://guest:guest@rabbitmq:5672/` | URL RabbitMQ |
| `ELASTICSEARCH_URL` | `http://elasticsearch:9200` | URL Elasticsearch |
| `LOGS_QUEUE` | `logs_queue` | Nome da fila AMQP |
| `LOGS_INDEX` | `logs` | Índice Elasticsearch |
| `PRODUCER_INTERVAL_SECONDS` | `5` | Intervalo de produção de logs (segundos) |
| `DNS_ZONES_FILE` | `/app/zones.json` | Arquivo de zonas DNS |
| `KIBANA_PORT` | `5601` | Porta host do Kibana |
| `ES_PORT` | `9200` | Porta host do Elasticsearch |
| `RABBITMQ_MGMT_PORT` | `15672` | Porta host do RabbitMQ Management |
| `TCP_ECHO_PORT` | `8888` | Porta host do TCP echo |
| `FTP_PORT` | `2121` | Porta host do FTP |
| `QUIZ_PORT` | `8080` | Porta host do Network Quiz |

---

## Troubleshooting

### Elasticsearch não sobe (falta de memória)

O Elasticsearch exige no mínimo 1 GB de heap. Ajuste via variável de ambiente antes de `make up`:

```bash
export ES_JAVA_OPTS="-Xms512m -Xmx512m"
make up
```

Ou adicione ao `.env`:

```
ES_JAVA_OPTS=-Xms512m -Xmx512m
```

### Conflito de portas

Se alguma porta já estiver em uso no host, edite `StartNewServer/.env` antes de `make up`:

```bash
# Exemplo: mudar porta do Kibana para 5602
KIBANA_PORT=5602
```

### Serviços DHCP/DNS/NTP não respondem externamente

Por design, DHCP (UDP 67), DNS (UDP 53) e NTP (UDP 123) ficam apenas na rede interna Docker para evitar conflitos com serviços do sistema. Use `make demo` para exercitá-los ou acesse via outro container na mesma rede.

### `make test` falha

Certifique-se de ter Go 1.21+ instalado localmente:

```bash
go version
# Deve ser >= go1.21
```

---

## Roadmap / Limitações Conhecidas

- **SMTP real não implementado.** O antigo `smtp-simulation-server` era na verdade um quiz TCP de redes; foi renomeado para `network-quiz` para refletir sua função real.
- **`session-manager-c` é opcional** (`--profile extras`) devido ao build pesado de vcpkg/mongo-c-driver. Não é necessário para a demo principal.
- **Pacotes DHCP/DNS são didáticos**, não conformes ao formato wire RFC. São adequados para fins de aprendizado, mas não para uso em produção.
- **DHCP/DNS/NTP não são publicados para o host** para evitar conflitos com serviços privilegiados do sistema operacional.
- **Submódulos C/Go originais** (`smtp_server_c`, `dns_server_c`, `golang-dns`) eram privados/inacessíveis e foram removidos do projeto ativo.
- **Próximos passos:** conformidade RFC para DNS/DHCP; testes de integração automatizados; dashboard Kibana pré-configurado; SMTP real.

---

## Mini Curso de Redes

Este projeto nasceu de um mini curso didático sobre redes de computadores, ministrado em ambiente acadêmico.

### Parte 1 — Fundamentos de Redes (Semana de SI)

**Tema:** Introdução aos fundamentos de redes de computadores.

- Visão geral dos princípios básicos de redes
- Implementação de um servidor TCP simples usando sockets
- Demonstração de comunicação cliente-servidor
- **Atividade prática:** servidor de eco TCP (hoje: `pkg/TCP/echo_server.go`)

### Parte 2 — Tecnologias Avançadas e Segurança (Semana Universitária)

**Tema:** Tecnologias avançadas e segurança em redes.

- Servidor HTTP com gerenciamento de rotas em Go
- Conceitos avançados de roteamento e segurança
- **Atividade prática:** servidor HTTP com roteamento de URLs e discussão sobre proteção contra ataques comuns

### Abordagem Pedagógica

- **Integração conceitual:** cada módulo combina teoria de redes com implementação prática em Go ou C.
- **Progressão lógica:** dos fundamentos TCP/IP até serviços completos (DHCP, DNS, NTP, FTP).
- **Projeto cumulativo:** os participantes desenvolvem servidores funcionais ao longo do curso.
- **Foco em redes, não em programação:** o código é apresentado como ilustração dos conceitos — não como objetivo final.

> O desenvolvimento prático é uma apresentação da lógica de criação dos serviços; o foco é o entendimento de redes de computadores.

---

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

Contribuições são bem-vindas via pull requests. Para mudanças maiores, abra uma issue primeiro para discutir o que você gostaria de alterar.
