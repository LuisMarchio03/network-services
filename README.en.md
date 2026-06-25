[Português](README.md) | [English](README.en.md)

# network-services

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose%20v2-2496ED?logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

A didactic project that recreates network services (DHCP, DNS, NTP, FTP, TCP echo, log pipeline) from scratch in Go/C, runnable as a reproducible local demo via Docker Compose.

---

## Architecture

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

## Service Matrix

| Service | Language | Internal port | Host port | Status |
|---|---|---|---|---|
| MongoDB | — | 27017 | 27017 | infra |
| RabbitMQ (+mgmt) | — | 5672 / 15672 | 5672 / 15672 | infra |
| Elasticsearch | — | 9200 | 9200 | infra |
| Kibana | — | 5601 | 5601 | infra |
| seed (Node) | JS | — | — | one-shot (seeds MongoDB) |
| logs-producer | Go | — | — | ✅ publishes logs via AMQP |
| logs-receiver | Go | — | — | ✅ consumes AMQP → Elasticsearch |
| dhcp-server | Go | UDP 67 | — (internal) | ✅ allocates IP from MongoDB |
| dhcp-client | Go | — | — | ✅ requests IP via DHCP |
| dns-server | Go | UDP 53 | — (internal) | 🧪 didactic (zones.json) |
| dns-client | Go | — | — | 🧪 queries didactic DNS |
| ntp-server | Go | UDP 123 | — (internal) | ✅ SNTP (era-1900 timestamps) |
| ntp-client | Go | — | — | ✅ queries NTP |
| tcp-echo | Go | 8888 | 8888 | ✅ TCP echo server |
| ftp-server | Go | 21 | 2121 | 🧪 control + PASV LIST/RETR/STOR |
| network-quiz | Go | 8080 | 8080 | 🧪 TCP networking quiz |
| session-manager | C | — | — | ⚙️ optional (profile `extras`) |

> **Legend:** ✅ working · 🧪 didactic/experimental · ⚙️ optional

---

## Quickstart

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) + Docker Compose v2
- `make`

### Start the stack

```bash
make up
```

Automatically copies `.env` from `.env.example` if it does not exist, then runs:

```bash
docker compose -f StartNewServer/docker-compose.yml up -d --build
```

### Full demo

```bash
make demo
```

The `scripts/demo.sh` script brings up the stack, exercises DNS, NTP, TCP echo, displays the Elasticsearch logs count, and tails DHCP logs.

### Other `make` targets

| Target | What it does |
|---|---|
| `make down` | Stop and remove containers |
| `make logs` | Follow container logs (last 50 lines) |
| `make test` | Run Go tests in `pkg/dhcp_server` and `pkg/NTP` |
| `make clean` | `down -v` (removes volumes) |
| `make dev` | Guidance for running services natively with `go run` |

### Useful links after `make up`

| Interface | URL |
|---|---|
| Kibana | http://localhost:5601 |
| RabbitMQ Management | http://localhost:15672 (guest/guest) |
| Elasticsearch | http://localhost:9200 |

---

## Per-Service Usage

### Log Pipeline (logs-producer → RabbitMQ → logs-receiver → Elasticsearch → Kibana)

```bash
# Check logs arriving in Elasticsearch
curl -s http://localhost:9200/logs/_count | python3 -m json.tool
# {"count": 42, ...}

# Follow container logs
make logs
```

### DHCP

```bash
# Follow IP allocations (internal server)
docker compose -f StartNewServer/docker-compose.yml logs dhcp-server

# Example output
# dhcp-server  | Received DHCP Discover from 02:42:ac:12:00:05
# dhcp-server  | Offering IP 10.0.0.10 to 02:42:ac:12:00:05
# dhcp-server  | ACK sent: 10.0.0.10 assigned
```

### DNS (clients profile)

```bash
docker compose -f StartNewServer/docker-compose.yml --profile clients run --rm dns-client

# Example output
# Querying DNS for example.local...
# Response: example.local -> 10.0.0.1
```

### NTP (clients profile)

```bash
docker compose -f StartNewServer/docker-compose.yml --profile clients run --rm ntp-client

# Example output
# NTP response received
# Server time: 2024-01-15 14:32:10 UTC
# Offset: +0.003s
```

### TCP Echo

```bash
# Connect to the echo server
nc localhost 8888
Hello network-services!
# Hello network-services!   <- echoed by the server
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
ftp> get file.txt
ftp> quit
```

### Network Quiz

```bash
# Access the networking quiz via TCP
nc localhost 8080
# Welcome to the Network Quiz!
# Question 1: What is a MAC address?
```

---

## Project Structure

```
network-services/
├── Makefile                        # Targets: up, down, demo, logs, test, clean
├── .env.example                    # Reference environment variables
├── scripts/
│   ├── demo.sh                     # End-to-end automated demo
│   └── smoke.sh                    # Basic smoke tests
├── pkg/                            # Go service implementations
│   ├── dhcp_server/                # DHCP server (allocates IPs from MongoDB)
│   │   ├── cmd/                    # Entrypoint
│   │   ├── internal/               # Internal logic
│   │   ├── mongo/                  # MongoDB integration
│   │   └── go.mod
│   ├── dhcp-simulation-client/     # DHCP client
│   ├── dhcp-simulation-server/     # DHCP server (simulation version)
│   ├── dns-simulation-server/      # Didactic DNS server
│   ├── dns-simulation-cliente/     # DNS client
│   ├── NTP/                        # NTP server and client (SNTP)
│   │   ├── cmd/                    # Server/client entrypoints
│   │   ├── ntptime/                # Era-1900 timestamp conversion
│   │   └── go.mod
│   ├── FTP/                        # FTP server (control + PASV)
│   ├── TCP/                        # TCP echo server
│   ├── network-quiz/               # TCP networking quiz
│   ├── logs-producer-teste/        # AMQP log producer
│   ├── logs-receiver/              # Log receiver → Elasticsearch
│   ├── session-manager-c/          # Session manager in C (optional)
│   └── deprecated/                 # Legacy unused code
└── StartNewServer/                 # Docker Compose orchestration
    ├── docker-compose.yml          # All services definition
    ├── src/                        # Node.js seed (seeds MongoDB)
    │   ├── setup.js
    │   └── rabbitmq-elasticsearch.js
    ├── elasticsearch.yml           # Elasticsearch configuration
    └── Dockerfile                  # Base image
```

---

## Environment Variables

Copy `.env.example` to `StartNewServer/.env` (done automatically by `make up`):

| Variable | Default | Description |
|---|---|---|
| `MONGO_URI` | `mongodb://mongodb:27017` | MongoDB connection URI |
| `MONGO_DB` | `mydb` | Database name |
| `RABBITMQ_URL` | `amqp://guest:guest@rabbitmq:5672/` | RabbitMQ URL |
| `ELASTICSEARCH_URL` | `http://elasticsearch:9200` | Elasticsearch URL |
| `LOGS_QUEUE` | `logs_queue` | AMQP queue name |
| `LOGS_INDEX` | `logs` | Elasticsearch index |
| `PRODUCER_INTERVAL_SECONDS` | `5` | Log production interval (seconds) |
| `DNS_ZONES_FILE` | `/app/zones.json` | DNS zones file path |
| `KIBANA_PORT` | `5601` | Host port for Kibana |
| `ES_PORT` | `9200` | Host port for Elasticsearch |
| `RABBITMQ_MGMT_PORT` | `15672` | Host port for RabbitMQ Management |
| `TCP_ECHO_PORT` | `8888` | Host port for TCP echo |
| `FTP_PORT` | `2121` | Host port for FTP |
| `QUIZ_PORT` | `8080` | Host port for Network Quiz |

---

## Troubleshooting

### Elasticsearch does not start (out of memory)

Elasticsearch requires at least 1 GB of heap. Adjust via environment variable before `make up`:

```bash
export ES_JAVA_OPTS="-Xms512m -Xmx512m"
make up
```

Or add to `.env`:

```
ES_JAVA_OPTS=-Xms512m -Xmx512m
```

### Port conflicts

If any port is already in use on the host, edit `StartNewServer/.env` before `make up`:

```bash
# Example: change Kibana port to 5602
KIBANA_PORT=5602
```

### DHCP/DNS/NTP servers not reachable externally

By design, DHCP (UDP 67), DNS (UDP 53), and NTP (UDP 123) are only on the internal Docker network to avoid conflicts with system services. Use `make demo` to exercise them, or access via another container on the same network.

### `make test` fails

Make sure Go 1.21+ is installed locally:

```bash
go version
# Should be >= go1.21
```

---

## Roadmap / Known Limitations

- **Real SMTP not implemented.** The former `smtp-simulation-server` was actually a TCP networking quiz; it was renamed `network-quiz` to reflect its actual function.
- **`session-manager-c` is optional** (`--profile extras`) due to the heavy vcpkg/mongo-c-driver build. Not required for the main demo.
- **DHCP/DNS packets are didactic**, not RFC wire-format compliant. Suitable for learning purposes, not for production use.
- **DHCP/DNS/NTP are not published to the host** to avoid conflicts with privileged OS system services.
- **Original C/Go submodules** (`smtp_server_c`, `dns_server_c`, `golang-dns`) were private/inaccessible and have been removed from the active project.
- **Next steps:** RFC compliance for DNS/DHCP; automated integration tests; pre-configured Kibana dashboard; real SMTP.

---

## Mini Networking Course

This project originated from a didactic mini course on computer networking, delivered in an academic setting.

### Part 1 — Networking Fundamentals (Computer Science Week)

**Topic:** Introduction to computer networking fundamentals.

- Overview of basic networking principles
- Implementation of a simple TCP server using sockets
- Client-server communication demonstration
- **Hands-on activity:** TCP echo server (today: `pkg/TCP/echo_server.go`)

### Part 2 — Advanced Technologies and Security (University Week)

**Topic:** Advanced networking technologies and security.

- HTTP server with URL routing in Go
- Advanced routing and security concepts
- **Hands-on activity:** HTTP server with URL routing and discussion of protection against common attacks

### Teaching Approach

- **Conceptual integration:** each module combines networking theory with practical implementation in Go or C.
- **Logical progression:** from TCP/IP fundamentals to full services (DHCP, DNS, NTP, FTP).
- **Cumulative project:** participants build functional servers throughout the course.
- **Focus on networks, not programming:** code is presented to illustrate concepts — not as an end goal.

> Hands-on development is a presentation of service creation logic; the focus is on understanding computer networks.

---

## License

This project is licensed under the [MIT License](LICENSE).

Contributions are welcome via pull requests. For major changes, please open an issue first to discuss what you would like to change.
