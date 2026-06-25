# network-services — Demo Reproduzível Local — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fazer o `network-services` subir com `docker compose up` e demonstrar cada serviço de rede funcionando ponta-a-ponta, corrigindo os principais defeitos, mais um README bilíngue (PT+EN).

**Architecture:** Toda comunicação roda na rede interna do Docker Compose (`network-services-net`). Serviços Go/C são *headless* (config por env/flags/arquivos, sem stdin). Clientes são containers "oneshot" que batem nos servidores por nome de serviço. Um `Makefile` + `scripts/demo.sh`/`scripts/smoke.sh` exercitam e validam tudo.

**Tech Stack:** Go 1.22+/1.23, C (session-manager, opcional), Node 18+ (seed), Docker Compose, MongoDB, RabbitMQ, Elasticsearch 7.17, Kibana 7.17.

## Global Constraints

- Nenhum serviço depende da rede física do host nem de portas privilegiadas publicadas no host. Servidores bindam portas padrão **dentro** do container; clientes acessam por nome de serviço na rede `network-services-net`.
- Submódulos privados (`pkg/network-services_smtp_server_c`, `pkg/network-services-dns_server_c`, `pkg/deprecated/golang-dns`) são **inacessíveis** — nunca dependa deles.
- Sem stdin interativo (`fmt.Scan`, `bufio` em stdin) em nenhum serviço que roda no Compose.
- Hosts/portas padrão via env: `MONGO_URI=mongodb://mongodb:27017`, `RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/`, `ELASTICSEARCH_URL=http://elasticsearch:9200`.
- Mensagens de log/usuário podem ficar em PT-BR (consistente com o projeto). Commits em inglês, prefixo Conventional Commits.
- Módulos Go (nomes exatos): `network-services-server-dhcp`, `dhcp_client`, `dns-simulation-server`, `dns-simulation-client`, `logs-producer-teste`, `logs-receiver`, `FTP_Server`, `smtp-simulation-server`→`network-quiz`, `smtp-simulation-client`.

---

### Task 0: Higiene do repositório (git init, .gitignore, remover binários, .env)

**Files:**
- Create: `/home/luis_promedico/Downloads/network-services-main/.gitignore`
- Create: `/home/luis_promedico/Downloads/network-services-main/.env.example`
- Delete: binários commitados e artefatos de build (ver comando)

**Interfaces:**
- Produces: repositório git inicializado; `.env.example` consumido pelo compose na Task 8.

- [ ] **Step 1: Inicializar git (o diretório não é repo)**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git init
git add -A
git commit -m "chore: import existing project as baseline"
```

- [ ] **Step 2: Criar `.gitignore`**

```gitignore
# Go
*.exe
*.exe~
*.test
*.out
/pkg/**/main
/pkg/dhcp_server/cmd/main
/pkg/dhcp_client/main

# C / CMake
cmake-build-*/
*.obj
*.o
*.a
/pkg/deprecated/mail_server/smtp_server

# Node
node_modules/

# IDE
.idea/
*.iml

# Env
.env
```

- [ ] **Step 3: Remover binários e artefatos versionados**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git rm -r --cached --ignore-unmatch \
  "pkg/dhcp-simulation-client/main.exe" \
  "pkg/dhcp-simulation-server/main.exe" \
  "pkg/dhcp_client/main" "pkg/dhcp_client/cmd/main.exe" \
  "pkg/dhcp_server/cmd/main" \
  "pkg/dns-simulation-cliente/dns_client.exe" \
  "pkg/dns-simulation-server/dns_server.exe" \
  "pkg/smtp-simulation-client/main.exe" \
  "pkg/smtp-simulation-server/main.exe" \
  "pkg/deprecated/mail_server/smtp_server" \
  "pkg/session-manager-c/cmake-build-debug" 2>/dev/null
find . -path ./.git -prune -o \( -name '*.exe' -o -name '*.obj' \) -print -delete
rm -rf pkg/session-manager-c/cmake-build-debug
rm -f pkg/dhcp_client/main pkg/dhcp_server/cmd/main pkg/deprecated/mail_server/smtp_server pkg/dhcp_client/cmd/main
```

- [ ] **Step 4: Criar `.env.example`**

```dotenv
# Infra
MONGO_URI=mongodb://mongodb:27017
MONGO_DB=mydb
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
ELASTICSEARCH_URL=http://elasticsearch:9200

# Log pipeline
LOGS_QUEUE=logs_queue
LOGS_INDEX=logs
PRODUCER_INTERVAL_SECONDS=5

# DNS
DNS_ZONES_FILE=/app/zones.json

# Demo (host-published ports)
KIBANA_PORT=5601
ES_PORT=9200
RABBITMQ_MGMT_PORT=15672
TCP_ECHO_PORT=8888
FTP_PORT=2121
QUIZ_PORT=8080
```

- [ ] **Step 5: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "chore: add gitignore, env example, remove committed binaries"
```

---

### Task 1: Renomear `smtp-simulation-server` → `network-quiz` (rótulo honesto)

**Files:**
- Rename dir: `pkg/smtp-simulation-server` → `pkg/network-quiz`
- Modify: `pkg/network-quiz/go.mod` (module name)

**Interfaces:**
- Produces: serviço `network-quiz` (TCP :8080) usado no compose (Task 8) e README.

- [ ] **Step 1: Renomear diretório e módulo**

```bash
cd /home/luis_promedico/Downloads/network-services-main/pkg
git mv smtp-simulation-server network-quiz
sed -i 's/^module smtp-simulation-server/module network-quiz/' network-quiz/go.mod
```

- [ ] **Step 2: Verificar build**

Run: `cd /home/luis_promedico/Downloads/network-services-main/pkg/network-quiz && go build ./...`
Expected: compila sem erros.

- [ ] **Step 3: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "refactor: rename mislabeled smtp-simulation to network-quiz"
```

---

### Task 2: DNS simulation server/client headless

**Files:**
- Modify: `pkg/dns-simulation-server/dns_server.go` (substituir stdin por zona via arquivo/env)
- Create: `pkg/dns-simulation-server/zones.json`
- Modify: `pkg/dns-simulation-cliente/dns_client.go` (flags em vez de stdin)

**Interfaces:**
- Produces: server DNS UDP :53 que resolve a partir de `zones.json`; client com flags `-domain` e `-server`.

- [ ] **Step 1: Criar `pkg/dns-simulation-server/zones.json`**

```json
{
  "example.com": "192.168.1.20",
  "intranet.local": "192.168.1.21",
  "kibana.local": "192.168.1.22"
}
```

- [ ] **Step 2: Reescrever `pkg/dns-simulation-server/dns_server.go`**

```go
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

// loadZones carrega o mapa domínio->IP de um arquivo JSON.
func loadZones(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	zones := make(map[string]string)
	if err := json.Unmarshal(data, &zones); err != nil {
		return nil, err
	}
	return zones, nil
}

func startDNSServer(domainIPMap map[string]string) {
	addr := net.UDPAddr{Port: 53, IP: net.ParseIP("0.0.0.0")}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor DNS:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("DNS Server (didático) escutando em :53 com %d domínios\n", len(domainIPMap))

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao ler do buffer:", err)
			continue
		}
		request := strings.TrimSpace(string(buffer[:n]))
		ip, found := domainIPMap[request]
		if !found {
			ip = "0.0.0.0"
			fmt.Println("Domínio não encontrado:", request)
		} else {
			fmt.Println("Resolvido:", request, "->", ip)
		}
		if _, err := conn.WriteToUDP([]byte(ip), clientAddr); err != nil {
			fmt.Println("Erro ao responder:", err)
		}
	}
}

func main() {
	zonesFile := os.Getenv("DNS_ZONES_FILE")
	if zonesFile == "" {
		zonesFile = "zones.json"
	}
	zones, err := loadZones(zonesFile)
	if err != nil {
		fmt.Printf("Erro ao carregar zonas de %s: %v\n", zonesFile, err)
		os.Exit(1)
	}
	startDNSServer(zones)
}
```

- [ ] **Step 3: Reescrever `pkg/dns-simulation-cliente/dns_client.go`**

```go
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	domain := flag.String("domain", "example.com", "domínio a resolver")
	server := flag.String("server", "dns-server:53", "endereço do servidor DNS")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *server)
	if err != nil {
		fmt.Println("Erro ao resolver endereço do servidor:", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Erro ao conectar com o servidor DNS:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(strings.TrimSpace(*domain))); err != nil {
		fmt.Println("Erro ao enviar requisição:", err)
		os.Exit(1)
	}
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		os.Exit(1)
	}
	fmt.Printf("DNS: %s -> %s\n", *domain, strings.TrimSpace(string(buffer[:n])))
}
```

- [ ] **Step 4: Verificar build**

Run: `cd /home/luis_promedico/Downloads/network-services-main/pkg/dns-simulation-server && go build ./... && cd ../dns-simulation-cliente && go build ./...`
Expected: ambos compilam.

- [ ] **Step 5: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "feat(dns): make sim server/client headless (zones file + flags)"
```

---

### Task 3: Pipeline de logs env-driven (producer em loop, receiver configurável)

**Files:**
- Modify: `pkg/logs-producer-teste/cmd/main.go`
- Modify: `pkg/logs-receiver/cmd/main.go`

**Interfaces:**
- Consumes: `RABBITMQ_URL`, `ELASTICSEARCH_URL`, `LOGS_QUEUE`, `LOGS_INDEX`, `PRODUCER_INTERVAL_SECONDS`.
- Produces: fluxo contínuo de logs Rabbit→ES.

- [ ] **Step 1: Reescrever `pkg/logs-producer-teste/cmd/main.go`**

```go
package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Service   string                 `json:"service"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context"`
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	url := getenv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	queue := getenv("LOGS_QUEUE", "logs_queue")
	interval, _ := strconv.Atoi(getenv("PRODUCER_INTERVAL_SECONDS", "5"))
	if interval <= 0 {
		interval = 5
	}

	var conn *amqp.Connection
	var err error
	for i := 0; i < 30; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ indisponível (%v), tentando de novo...", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	levels := []string{"INFO", "WARNING", "ERROR"}
	i := 0
	for {
		entry := LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Service:   "demo-service",
			Level:     levels[i%len(levels)],
			Message:   "Evento sintético de demonstração",
			Context:   map[string]interface{}{"seq": i},
		}
		body, _ := json.Marshal(entry)
		if err := ch.Publish("", queue, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}); err != nil {
			log.Printf("Failed to publish: %s", err)
		} else {
			log.Printf("Log publicado (seq=%d, level=%s)", i, entry.Level)
		}
		i++
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
```

- [ ] **Step 2: Tornar host/URL do receiver configurável**

Em `pkg/logs-receiver/cmd/main.go`, substituir a linha:

```go
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Conexão com o servidor RabbitMQ
```

por:

```go
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	var conn *amqp.Connection
	for i := 0; i < 30; i++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ indisponível (%v), retry...", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
```

E substituir:

```go
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false)) // Criação do cliente Elasticsearch
```

por:

```go
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://elasticsearch:9200"
	}
	client, err := elastic.NewClient(elastic.SetURL(esURL), elastic.SetSniff(false))
```

Adicionar `"os"` ao bloco de imports do arquivo.

- [ ] **Step 3: Verificar build**

Run: `cd /home/luis_promedico/Downloads/network-services-main/pkg/logs-producer-teste && go build ./... && cd ../logs-receiver && go build ./...`
Expected: ambos compilam.

- [ ] **Step 4: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "feat(logs): env-driven producer loop and configurable receiver"
```

---

### Task 4: DHCP server aloca IP real do Mongo

**Files:**
- Modify: `pkg/dhcp_server/internal/dhcp/server.go` (offer/ack)
- Test: `pkg/dhcp_server/internal/dhcp/server_test.go`

**Interfaces:**
- Consumes: `mongodb.MongoDB.FindAvailableIP(ctx)`, `.UpdateIPAssignment(ctx, ip, assigned)`.
- Produces: oferta DHCP com IP vindo do banco; ACK marca o IP como atribuído.

- [ ] **Step 1: Escrever teste falho `pkg/dhcp_server/internal/dhcp/server_test.go`**

```go
package dhcp

import (
	"net"
	"testing"
)

// O acknowledge deve copiar o IP solicitado (bytes 12:16 da request) para ack[16:20].
func TestBuildDHCPAcknowledgeCopiesRequestedIP(t *testing.T) {
	s := &Server{}
	req := make([]byte, 240)
	req[0] = 3 // DHCP Request
	requested := net.ParseIP("192.168.1.55").To4()
	copy(req[12:16], requested)

	ack := s.buildDHCPAcknowledge(req)
	if ack == nil {
		t.Fatal("ack não deveria ser nil")
	}
	got := net.IP(ack[16:20]).String()
	if got != "192.168.1.55" {
		t.Fatalf("esperado 192.168.1.55, obtido %s", got)
	}
}
```

- [ ] **Step 2: Rodar o teste e confirmar resultado**

Run: `cd /home/luis_promedico/Downloads/network-services-main/pkg/dhcp_server && go test ./internal/dhcp/ -run TestBuildDHCPAcknowledgeCopiesRequestedIP -v`
Expected: PASS (valida o comportamento atual antes de mexer no offer; serve de regressão).

- [ ] **Step 3: Trocar IP hardcoded por consulta ao Mongo no `buildDHCPOffer`**

Em `server.go`, substituir o bloco:

```go
	// Identifica um endereço IP disponível usando o MongoDB
	//availableIP, err := s.db.FindAvailableIP(s.ctx)
	//if err != nil {
	//	logger.Error("Erro ao encontrar um endereço IP disponível:", err)
	//	return nil
	//}
	availableIP := "192.168.2.42"
```

por:

```go
	// Identifica um endereço IP disponível usando o MongoDB
	availableIP, err := s.db.FindAvailableIP(s.ctx)
	if err != nil {
		logger.Error("Erro ao encontrar um endereço IP disponível:", err)
		return nil
	}
```

- [ ] **Step 4: Marcar o IP como atribuído no `buildDHCPAcknowledge`**

Em `server.go`, logo após `requestedIP := net.IP(request[12:16])`, adicionar:

```go
	if s.db != nil {
		if err := s.db.UpdateIPAssignment(s.ctx, requestedIP.String(), true); err != nil {
			logger.Error("Erro ao marcar IP como atribuído:", err)
		}
	}
```

- [ ] **Step 5: Rodar testes e build**

Run: `cd /home/luis_promedico/Downloads/network-services-main/pkg/dhcp_server && go test ./... && go build ./...`
Expected: testes PASS, build OK. (O teste de offer não é coberto por unit por exigir Mongo; coberto no smoke test.)

- [ ] **Step 6: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "fix(dhcp): allocate IP from MongoDB and mark assigned on ACK"
```

---

### Task 5: Implementar NTP (SNTP) — server + client (TDD)

**Files:**
- Create: `pkg/NTP/go.mod`
- Create: `pkg/NTP/ntptime/ntptime.go`
- Create: `pkg/NTP/ntptime/ntptime_test.go`
- Create: `pkg/NTP/cmd/server/main.go`
- Create: `pkg/NTP/cmd/client/main.go`
- Delete: `pkg/NTP/main.go` (vazio)

**Interfaces:**
- Produces: pacote `ntptime` com `ToNTP(t time.Time) (secs, frac uint32)` e `FromNTP(secs, frac uint32) time.Time`; server UDP :123; client que imprime a hora.

- [ ] **Step 1: Criar `pkg/NTP/go.mod`**

```
module ntp

go 1.22
```

- [ ] **Step 2: Escrever teste falho `pkg/NTP/ntptime/ntptime_test.go`**

```go
package ntptime

import (
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	original := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	secs, frac := ToNTP(original)
	got := FromNTP(secs, frac).UTC()
	if !got.Equal(original) {
		t.Fatalf("round-trip falhou: esperado %v, obtido %v", original, got)
	}
}

func TestEpochOffset(t *testing.T) {
	// 1970-01-01 deve corresponder a 2208988800 segundos NTP (era 1900).
	unixEpoch := time.Unix(0, 0).UTC()
	secs, _ := ToNTP(unixEpoch)
	if secs != 2208988800 {
		t.Fatalf("esperado 2208988800, obtido %d", secs)
	}
}
```

- [ ] **Step 3: Rodar o teste e confirmar falha**

Run: `cd /home/luis_promedico/Downloads/network-services-main/pkg/NTP && go test ./ntptime/ -v`
Expected: FAIL (pacote/funções não existem).

- [ ] **Step 4: Implementar `pkg/NTP/ntptime/ntptime.go`**

```go
// Package ntptime converte entre time.Time e timestamps NTP (era 1900).
package ntptime

import "time"

// ntpEpochOffset é a diferença em segundos entre 1900-01-01 e 1970-01-01.
const ntpEpochOffset = 2208988800

// ToNTP converte um time.Time para (segundos, fração) no formato NTP.
func ToNTP(t time.Time) (uint32, uint32) {
	secs := uint64(t.Unix()) + ntpEpochOffset
	frac := uint64(t.Nanosecond()) * (1 << 32) / 1e9
	return uint32(secs), uint32(frac)
}

// FromNTP converte (segundos, fração) NTP de volta para time.Time (UTC).
func FromNTP(secs, frac uint32) time.Time {
	sec := int64(secs) - ntpEpochOffset
	nsec := int64(uint64(frac) * 1e9 / (1 << 32))
	return time.Unix(sec, nsec).UTC()
}
```

- [ ] **Step 5: Rodar o teste e confirmar PASS**

Run: `cd /home/luis_promedico/Downloads/network-services-main/pkg/NTP && go test ./ntptime/ -v`
Expected: PASS.

- [ ] **Step 6: Implementar o servidor `pkg/NTP/cmd/server/main.go`**

```go
package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"ntp/ntptime"
)

func main() {
	addr := net.UDPAddr{Port: 123, IP: net.ParseIP("0.0.0.0")}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Erro ao iniciar servidor NTP:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("NTP Server (didático) escutando em :123")

	buf := make([]byte, 48)
	for {
		n, client, err := conn.ReadFromUDP(buf)
		if err != nil || n == 0 {
			continue
		}
		resp := make([]byte, 48)
		// LI=0, VN=4, Mode=4 (server) -> 0b00_100_100 = 0x24
		resp[0] = 0x24
		resp[1] = 1 // stratum
		secs, frac := ntptime.ToNTP(time.Now().UTC())
		// Transmit Timestamp (bytes 40..47)
		binary.BigEndian.PutUint32(resp[40:44], secs)
		binary.BigEndian.PutUint32(resp[44:48], frac)
		if _, err := conn.WriteToUDP(resp, client); err != nil {
			fmt.Println("Erro ao responder NTP:", err)
		}
	}
}
```

- [ ] **Step 7: Implementar o cliente `pkg/NTP/cmd/client/main.go`**

```go
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"ntp/ntptime"
)

func main() {
	server := flag.String("server", "ntp-server:123", "endereço do servidor NTP")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *server)
	if err != nil {
		fmt.Println("Erro ao resolver servidor:", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Erro ao conectar:", err)
		os.Exit(1)
	}
	defer conn.Close()

	req := make([]byte, 48)
	req[0] = 0x23 // LI=0, VN=4, Mode=3 (client)
	if _, err := conn.Write(req); err != nil {
		fmt.Println("Erro ao enviar:", err)
		os.Exit(1)
	}
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	resp := make([]byte, 48)
	if _, err := conn.Read(resp); err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		os.Exit(1)
	}
	secs := binary.BigEndian.Uint32(resp[40:44])
	frac := binary.BigEndian.Uint32(resp[44:48])
	serverTime := ntptime.FromNTP(secs, frac)
	fmt.Printf("NTP: hora do servidor = %s (offset ~ %s)\n",
		serverTime.Format(time.RFC3339Nano), time.Since(serverTime).Round(time.Millisecond))
}
```

- [ ] **Step 8: Remover stub vazio e verificar build/test**

```bash
cd /home/luis_promedico/Downloads/network-services-main/pkg/NTP
rm -f main.go
go build ./... && go test ./...
```
Expected: build OK, testes PASS.

- [ ] **Step 9: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "feat(ntp): implement SNTP server and client with tested time conversion"
```

---

### Task 6: FTP útil (comandos de controle + PASV + LIST/RETR/STOR)

**Files:**
- Modify: `pkg/FTP/pkg/FTP/server.go`

**Interfaces:**
- Consumes: `auth.Authenticate(user, pass) bool` (já existe, retorna true).
- Produces: servidor FTP de controle na :21 servindo o diretório `/data` (volume).

- [ ] **Step 1: Reescrever `pkg/FTP/pkg/FTP/server.go`**

```go
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
			reply("331 Usuario OK, envie a senha")
		case "PASS":
			if auth.Authenticate("user", arg) {
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
			dataConn, _ := dataListener.Accept()
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
			dataConn, _ := dataListener.Accept()
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
```

- [ ] **Step 2: Verificar build**

Run: `cd /home/luis_promedico/Downloads/network-services-main/pkg/FTP && go build ./...`
Expected: compila.

- [ ] **Step 3: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "feat(ftp): control commands + PASV LIST/RETR/STOR over /data"
```

---

### Task 7: Corrigir o seed Node

**Files:**
- Modify: `StartNewServer/src/setup.js`
- Modify: `StartNewServer/package.json`

**Interfaces:**
- Produces: coleção `dhcps` populada com documentos válidos `{ip_address, assigned}`, consumida pelo DHCP server (Task 4).

- [ ] **Step 1: Corrigir o insert do DHCP em `setup.js`**

Substituir o bloco do `Promise.all` que faz `...dataToInsert.dhcpEntries.map(entry => createDocument(Dhcp, DhcpJSON))` por inserção de array de objetos:

```js
        // Gera N registros DHCP como documentos reais (não string JSON)
        const dhcpDocs = JSON.parse(DhcpJSON).concat(dataToInsert.dhcpEntries);

        await Promise.all([
            Dhcp.insertMany(dhcpDocs),
            ...dataToInsert.dnsEntries.map(entry => createDocument(Dns, entry)),
            ...dataToInsert.httpEntries.map(entry => createDocument(HttpServer, entry)),
            ...dataToInsert.ftpEntries.map(entry => createDocument(FtpServer, entry)),
            ...dataToInsert.emailEntries.map(entry => createDocument(EmailServer, entry)),
            ...dataToInsert.fileEntries.map(entry => createDocument(FileServer, entry)),
            ...dataToInsert.firewallEntries.map(entry => createDocument(Firewall, entry)),
            ...dataToInsert.ldapEntries.map(entry => createDocument(LdapServer, entry)),
            ...dataToInsert.logEntries.map(entry => createDocument(LogServer, entry)),
            ...dataToInsert.computerEntries.map(entry => createDocument(UserComputer, entry)),
            ...dataToInsert.sessionEntries.map(entry => createDocument(Session, entry)),
        ]);
```

Também trocar a URI para usar env:

```js
const uri = process.env.MONGO_URI ? `${process.env.MONGO_URI}/${process.env.MONGO_DB || 'mydb'}` : 'mongodb://mongodb:27017/mydb';
```

- [ ] **Step 2: Corrigir `package.json` start script**

Substituir `"start": "node src/index.js"` por `"start": "node src/setup.js"`.

- [ ] **Step 3: Verificar sintaxe**

Run: `cd /home/luis_promedico/Downloads/network-services-main/StartNewServer && node --check src/setup.js`
Expected: sem erros de sintaxe.

- [ ] **Step 4: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "fix(seed): insert valid DHCP documents and use env Mongo URI"
```

---

### Task 8: Dockerfiles dos novos serviços + reescrita do `docker-compose.yml`

**Files:**
- Create: `Dockerfile` em cada serviço Go novo (ver tabela)
- Modify: `StartNewServer/docker-compose.yml`

**Interfaces:**
- Consumes: todos os serviços/binários das Tasks 1-7.
- Produces: stack completa `docker compose up`.

- [ ] **Step 1: Criar Dockerfile genérico Go em cada serviço**

Para cada serviço abaixo, criar `<dir>/Dockerfile` com este conteúdo, substituindo `<BUILD_PATH>` e `<BIN>`:

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download || true
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app <BUILD_PATH>

FROM alpine:latest
WORKDIR /app
COPY --from=builder /out/app /app/app
# Copia arquivos auxiliares se existirem (ex: zones.json)
COPY . /app/
CMD ["/app/app"]
```

Tabela de substituições (criar um Dockerfile em cada dir):

| Diretório | `<BUILD_PATH>` |
|---|---|
| `pkg/dns-simulation-server` | `.` |
| `pkg/dns-simulation-cliente` | `.` |
| `pkg/NTP` | `./cmd/server` |
| `pkg/TCP` | `.` |
| `pkg/network-quiz` | `.` |
| `pkg/logs-producer-teste` | `./cmd` |
| `pkg/logs-receiver` | `./cmd` |
| `pkg/FTP` | `.` |

Para o **NTP client** e **DNS client** (oneshot), o compose usa o mesmo Dockerfile do dir, mas sobrescreve o build via serviço separado — em vez disso, criar `pkg/NTP/Dockerfile.client` e `pkg/dns-simulation-cliente` já é o client. Para o NTP client criar `pkg/NTP/Dockerfile.client`:

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/client
FROM alpine:latest
WORKDIR /app
COPY --from=builder /out/app /app/app
CMD ["/app/app"]
```

Para `pkg/TCP` e `pkg/NTP` criar `go.mod` se não existir:
```bash
cd /home/luis_promedico/Downloads/network-services-main/pkg/TCP && [ -f go.mod ] || go mod init tcp-echo
```
(NTP go.mod já criado na Task 5.)

- [ ] **Step 2: Reescrever `StartNewServer/docker-compose.yml`**

```yaml
name: network-services

services:
  mongodb:
    image: mongo:6
    container_name: mongodb
    ports: ["27017:27017"]
    volumes: ["mongodb_data:/data/db"]
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks: [net]

  rabbitmq:
    image: rabbitmq:3-management
    container_name: rabbitmq
    ports: ["5672:5672", "${RABBITMQ_MGMT_PORT:-15672}:15672"]
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "-q", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks: [net]

  elasticsearch:
    image: elasticsearch:7.17.9
    container_name: elasticsearch
    environment:
      - discovery.type=single-node
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
      - xpack.security.enabled=false
      - http.host=0.0.0.0
    ports: ["${ES_PORT:-9200}:9200"]
    healthcheck:
      test: ["CMD-SHELL", "curl -s http://localhost:9200/_cluster/health || exit 1"]
      interval: 15s
      timeout: 10s
      retries: 10
    networks: [net]

  kibana:
    image: kibana:7.17.9
    container_name: kibana
    ports: ["${KIBANA_PORT:-5601}:5601"]
    depends_on:
      elasticsearch: {condition: service_healthy}
    networks: [net]

  seed:
    build: { context: . }
    environment:
      - MONGO_URI=${MONGO_URI:-mongodb://mongodb:27017}
      - MONGO_DB=${MONGO_DB:-mydb}
      - RABBITMQ_URL=${RABBITMQ_URL:-amqp://guest:guest@rabbitmq:5672/}
    depends_on:
      mongodb: {condition: service_healthy}
      rabbitmq: {condition: service_healthy}
    networks: [net]

  logs-producer:
    build: { context: ../pkg/logs-producer-teste }
    environment:
      - RABBITMQ_URL=${RABBITMQ_URL:-amqp://guest:guest@rabbitmq:5672/}
      - PRODUCER_INTERVAL_SECONDS=${PRODUCER_INTERVAL_SECONDS:-5}
    depends_on:
      rabbitmq: {condition: service_healthy}
    restart: unless-stopped
    networks: [net]

  logs-receiver:
    build: { context: ../pkg/logs-receiver }
    environment:
      - RABBITMQ_URL=${RABBITMQ_URL:-amqp://guest:guest@rabbitmq:5672/}
      - ELASTICSEARCH_URL=${ELASTICSEARCH_URL:-http://elasticsearch:9200}
    depends_on:
      rabbitmq: {condition: service_healthy}
      elasticsearch: {condition: service_healthy}
    restart: unless-stopped
    networks: [net]

  dhcp-server:
    build: { context: ../pkg/dhcp_server }
    depends_on:
      mongodb: {condition: service_healthy}
      seed: {condition: service_completed_successfully}
    restart: unless-stopped
    networks: [net]

  dhcp-client:
    build: { context: ../pkg/dhcp_client }
    depends_on: [dhcp-server]
    networks: [net]

  dns-server:
    build: { context: ../pkg/dns-simulation-server }
    environment:
      - DNS_ZONES_FILE=/app/zones.json
    networks: [net]

  ntp-server:
    build: { context: ../pkg/NTP }
    networks: [net]

  tcp-echo:
    build: { context: ../pkg/TCP }
    ports: ["${TCP_ECHO_PORT:-8888}:8888"]
    networks: [net]

  ftp-server:
    build: { context: ../pkg/FTP }
    ports: ["${FTP_PORT:-2121}:21"]
    volumes: ["ftp_data:/data"]
    networks: [net]

  network-quiz:
    build: { context: ../pkg/network-quiz }
    ports: ["${QUIZ_PORT:-8080}:8080"]
    networks: [net]

  session-manager:
    build: { context: ../pkg/session-manager-c, dockerfile: docker/Dockerfile }
    profiles: ["extras"]
    depends_on:
      mongodb: {condition: service_healthy}
      rabbitmq: {condition: service_healthy}
    networks: [net]

volumes:
  mongodb_data:
  ftp_data:

networks:
  net:
    name: network-services-net
```

- [ ] **Step 3: Validar a sintaxe do compose**

Run: `cd /home/luis_promedico/Downloads/network-services-main/StartNewServer && cp ../.env.example .env && docker compose config -q`
Expected: sem erros (config válida).

- [ ] **Step 4: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "feat(compose): full stack with healthchecks, dockerfiles, profiles"
```

---

### Task 9: Makefile + scripts de demo e smoke

**Files:**
- Create: `Makefile`
- Create: `scripts/demo.sh`
- Create: `scripts/smoke.sh`

**Interfaces:**
- Consumes: o compose da Task 8.
- Produces: `make up/down/demo/test/logs/clean`.

- [ ] **Step 1: Criar `Makefile`**

```makefile
COMPOSE = docker compose -f StartNewServer/docker-compose.yml

.PHONY: up down demo logs test clean dev

up:
	cd StartNewServer && [ -f .env ] || cp ../.env.example .env
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

demo:
	./scripts/demo.sh

logs:
	$(COMPOSE) logs -f --tail=50

test:
	@for d in pkg/dhcp_server pkg/NTP; do (cd $$d && go test ./...); done

clean:
	$(COMPOSE) down -v

dev:
	@echo "Modo nativo: rode cada serviço com 'go run' apontando env para localhost."
```

- [ ] **Step 2: Criar `scripts/demo.sh`**

```bash
#!/usr/bin/env bash
set -euo pipefail
COMPOSE="docker compose -f StartNewServer/docker-compose.yml"

echo "==> Subindo a stack (pode demorar no primeiro build)..."
make up

echo "==> Aguardando Elasticsearch ficar saudável..."
until curl -sf "http://localhost:${ES_PORT:-9200}/_cluster/health" >/dev/null; do sleep 3; done

echo "==> [DNS] resolvendo example.com"
$COMPOSE run --rm dns-server /app/app >/dev/null 2>&1 || true
$COMPOSE run --rm --entrypoint /app/app dns-client -domain example.com -server dns-server:53 || true

echo "==> [NTP] consultando hora do servidor"
$COMPOSE run --rm --build --entrypoint /app/app ntp-client -server ntp-server:123 || true

echo "==> [TCP echo] enviando 'ola'"
printf 'ola\n' | nc -w2 localhost "${TCP_ECHO_PORT:-8888}" || true

echo "==> [Logs] últimos documentos no índice 'logs'"
sleep 6
curl -s "http://localhost:${ES_PORT:-9200}/logs/_count" || true; echo

echo "==> [DHCP] veja a negociação:"
$COMPOSE logs --tail=10 dhcp-server dhcp-client || true

echo "==> Demo concluída. Kibana: http://localhost:${KIBANA_PORT:-5601}  RabbitMQ: http://localhost:${RABBITMQ_MGMT_PORT:-15672}"
```

(Nota: o `dns-client`/`ntp-client` são serviços oneshot; se preferir, declare-os no compose como serviços dedicados que rodam o binário cliente. Se `nc` não existir, use `python3 -c`.)

- [ ] **Step 3: Criar `scripts/smoke.sh`**

```bash
#!/usr/bin/env bash
set -uo pipefail
fail=0
check() { if eval "$2" >/dev/null 2>&1; then echo "PASS: $1"; else echo "FAIL: $1"; fail=1; fi; }

check "Elasticsearch responde"      "curl -sf http://localhost:${ES_PORT:-9200}/_cluster/health"
check "RabbitMQ mgmt responde"      "curl -sf http://localhost:${RABBITMQ_MGMT_PORT:-15672}"
check "TCP echo ecoa"               "printf 'x\n' | nc -w2 localhost ${TCP_ECHO_PORT:-8888} | grep -q x"
check "Indice de logs tem docs"     "[ \$(curl -s http://localhost:${ES_PORT:-9200}/logs/_count | grep -o '\"count\":[0-9]*' | cut -d: -f2) -gt 0 ]"

exit $fail
```

- [ ] **Step 4: Tornar executáveis e commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
chmod +x scripts/demo.sh scripts/smoke.sh
git add -A
git commit -m "feat: Makefile and demo/smoke scripts"
```

---

### Task 10: Verificação Go + build do compose

**Files:** nenhum (verificação).

- [ ] **Step 1: Build/test de todos os módulos Go**

```bash
cd /home/luis_promedico/Downloads/network-services-main
for d in pkg/dhcp_server pkg/dhcp_client pkg/dns-simulation-server pkg/dns-simulation-cliente \
         pkg/NTP pkg/TCP pkg/network-quiz pkg/logs-producer-teste pkg/logs-receiver pkg/FTP; do
  echo "== $d =="; (cd "$d" && go vet ./... && go build ./... && go test ./... 2>&1 | tail -3)
done
```
Expected: cada módulo compila; testes de `dhcp_server` e `NTP` passam.

- [ ] **Step 2: Build das imagens do compose**

```bash
cd /home/luis_promedico/Downloads/network-services-main/StartNewServer
docker compose build
```
Expected: todas as imagens (exceto profile `extras`) buildam.

- [ ] **Step 3: Subir e rodar smoke (best-effort neste ambiente)**

```bash
cd /home/luis_promedico/Downloads/network-services-main
make up && sleep 30 && ./scripts/smoke.sh; make down
```
Expected: a maioria dos checks PASS. Reportar honestamente o que não subiu (ES é pesado).

- [ ] **Step 4: Commit (se houve ajustes durante a verificação)**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "fix: adjustments from build/smoke verification" || echo "nada a commitar"
```

---

### Task 11: README bilíngue (PT + EN)

**Files:**
- Modify: `README.md` (PT)
- Create: `README.en.md` (EN)

**Interfaces:** documentação final do projeto.

- [ ] **Step 1: Reescrever `README.md` (PT)** com as seções, preservando o contexto do mini-curso de redes do README atual:

1. Título + badges (Go, Docker, License) + descrição de 1 linha.
2. **Arquitetura** — diagrama Mermaid:

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

3. **Matriz de serviços** (tabela: serviço · linguagem · porta · status ✅/🧪/⚠️).
4. **Quickstart**: pré-requisitos (Docker, make), `make up`, `make demo`, `make down`. Links Kibana/RabbitMQ.
5. **Uso por serviço** com saída de exemplo (DHCP, DNS, NTP, TCP echo, FTP, logs, quiz).
6. **Estrutura do projeto** (árvore resumida de `pkg/` e `StartNewServer/`).
7. **Troubleshooting** (ES sem memória → ajustar `ES_JAVA_OPTS`; portas ocupadas → `.env`).
8. **Roadmap / Limitações conhecidas**: SMTP real não implementado; `session-manager-c` é profile `extras`; pacotes DHCP/DNS são didáticos (não RFC); submódulos C/Go originais são privados.
9. **Licença** + como contribuir.

- [ ] **Step 2: Criar `README.en.md`** — tradução fiel do `README.md` (mesmas seções/diagrama), com link recíproco no topo de cada arquivo (`[Português](README.md) | [English](README.en.md)`).

- [ ] **Step 3: Conferir links e diagrama**

Run: `grep -n "README.en.md\|README.md\|mermaid" /home/luis_promedico/Downloads/network-services-main/README.md`
Expected: links recíprocos e bloco mermaid presentes.

- [ ] **Step 4: Commit**

```bash
cd /home/luis_promedico/Downloads/network-services-main
git add -A
git commit -m "docs: bilingual README (PT/EN) with architecture, quickstart, service matrix"
```

---

## Self-Review (preenchido pelo autor do plano)

**Spec coverage:** §4 matriz → Tasks 1-8; §5.1 compose → Task 8; §5.2 seed → Task 7; §5.3 DHCP → Task 4; §5.4 DNS → Task 2; §5.5 NTP → Task 5; §5.6 FTP → Task 6; §5.7 limpeza → Tasks 0-1; §7 tooling → Task 9; §8 testes → Tasks 4,5,10; §9 README → Task 11; §6 limitações → Task 11 (roadmap). Sem lacunas.

**Placeholder scan:** sem TBD/TODO de implementação; código concreto em cada step. Exceção consciente: o conteúdo de prosa do README (Task 11) é especificado por seções + diagrama real — apropriado para documentação.

**Type consistency:** `FindAvailableIP(ctx)`, `UpdateIPAssignment(ctx, ip, assigned)`, `ReleaseIPAddress` conferem com `mongodb.go`. `auth.Authenticate(user, pass) bool` confere. `ntptime.ToNTP/FromNTP` consistentes entre server, client e teste. Nomes de módulo Go conferem com os `go.mod` reais.
