# Design — `network-services`: demo reproduzível local + README aprimorado

- **Data:** 2026-06-18
- **Autor:** Claude (brainstorming com o mantenedor)
- **Status:** Aprovado para implementação

## 1. Objetivo

Transformar o `network-services` (projeto didático que recria serviços de rede do zero em
Go/C) em algo que **roda com um único comando** e **demonstra cada serviço funcionando**,
mantendo o caráter educacional. Mais um README bilíngue (PT + EN) substancialmente melhor.

Decisões já tomadas no brainstorming:

- **Alvo:** demo reproduzível local — `docker compose up` sobe toda a stack; um `make demo`
  exercita cada serviço. Portável, sem depender da rede física, CI-friendly. (Abordagem A.)
- **Ambição:** fazer rodar ponta-a-ponta **e corrigir os principais defeitos** (DHCP alocando
  IP de verdade, NTP implementado, FTP minimamente útil). Sem reescrever tudo do zero.
- **README:** bilíngue PT + EN.

## 2. Não-objetivos (YAGNI)

- **Não** rodar serviços reais ligados à interface física do host (DHCP/DNS/FTP servindo a LAN
  real em portas privilegiadas com root). Tudo acontece dentro da rede do Docker Compose.
- **Não** implementar conformidade total com as RFCs (pacotes DHCP/DNS permanecem em formato
  didático simplificado, como já estão). O foco é clareza e demonstração, não interoperabilidade real.
- **Não** implementar um servidor SMTP real agora (ver §6, limitações documentadas).
- **Não** refatorar código que não serve ao objetivo.

## 3. Visão de arquitetura

Toda a comunicação ocorre na rede interna do Compose (`network-services-net`). Serviços Go/C
são *headless*: configuração por variáveis de ambiente / flags / arquivos, nunca por stdin.

```
                ┌─────────────────────────── docker compose ───────────────────────────┐
                │                                                                        │
   seed (Node) ─┼─► MongoDB ◄── DHCP server ◄─UDP─ DHCP client                           │
                │      ▲           │ aloca IP (FindAvailableIP/UpdateIPAssignment)        │
                │      ├───────────┴─ FTP server (auth via Mongo)                         │
                │      │                                                                  │
  logs-producer ─AMQP─► RabbitMQ ─AMQP─► logs-receiver ─HTTP─► Elasticsearch ─► Kibana    │
                │                                                                         │
   DNS server ◄─UDP─ DNS client          NTP server ◄─UDP─ NTP client      TCP echo:8888  │
                │                                                                         │
                │   (profile "extras") session-manager-c ◄AMQP─ RabbitMQ ─► MongoDB       │
                └─────────────────────────────────────────────────────────────────────────┘
```

## 4. Matriz de serviços

| Serviço | Linguagem | Estado hoje | Ação | Porta interna |
|---|---|---|---|---|
| MongoDB / RabbitMQ / Elasticsearch / Kibana | — | em `docker-compose.yml` | manter; rede nomeada + healthchecks | 27017 / 5672+15672 / 9200 / 5601 |
| Seed | Node | `setup.js` com bug; `package.json` aponta p/ `index.js` inexistente | corrigir insert do DHCP; `start` → `setup.js`; aguardar Mongo healthy | — |
| Pipeline de logs | Go | producer/receiver apontam `localhost` | host por env (`rabbitmq`/`elasticsearch`); producer em loop periódico | — |
| DHCP server + client | Go | server hardcoda IP; alocação no Mongo comentada | usar `FindAvailableIP` + `UpdateIPAssignment`; corrigir release | UDP 67 (interno) |
| DNS server + client | Go | lê domínios via stdin interativo | headless: zona via arquivo/env; client por flags | UDP 53 |
| NTP | Go | `main.go` **vazio** | **implementar** SNTP mínimo (server + client) | UDP 123 |
| TCP echo | Go | funciona | incluir no demo | 8888 |
| FTP | Go | só faz login | USER/PASS/PWD/LIST/QUIT (+STOR/RETR em volume); auth via Mongo | 21 |
| network-quiz (ex-`smtp-simulation`) | Go | quiz TCP rotulado errado como SMTP | renomear honestamente; incluir como demo opcional | 8080 |
| session-manager-c | C | consome Rabbit, grava sessões no Mongo; build vcpkg pesado | profile `extras` do compose; documentado | — |

## 5. Mudanças por componente

### 5.1 Orquestração (`StartNewServer/docker-compose.yml`)
- Rede nomeada `network-services-net`; remover `networks: default: null`.
- `healthcheck` em mongo/rabbit/elasticsearch; `depends_on: condition: service_healthy`.
- Adicionar serviços: `logs-producer`, `logs-receiver`, `dns-server`, `dns-client` (oneshot),
  `ntp-server`, `ntp-client` (oneshot), `tcp-echo`, `ftp-server`, `network-quiz`.
- DHCP server **não** publica `67` no host (evita conflito com dnsmasq); cliente acessa pela
  rede interna. Documentar.
- `session-manager-c` sob `profiles: [extras]`.
- `.env` / `.env.example` com hosts e portas.

### 5.2 Seed Node
- `setup.js`: inserir registros DHCP corretamente (array de documentos, não string JSON única).
- `package.json`: `start` → `node src/setup.js`; revisar dependências (remover `elasticsearch`
  legado se não usado).

### 5.3 DHCP (`pkg/dhcp_server`)
- `server.go`: substituir IP hardcoded por `db.FindAvailableIP(ctx)`; ao confirmar (ACK),
  `db.UpdateIPAssignment(ctx, ip, true)`; release já chama `ReleaseIPAddress`. Tratar
  "sem IP disponível" com log claro.
- Dockerfile: garantir `go.sum` presente; `dhcp_client` idem.
- Teste unitário: `buildDHCPOffer`/`buildDHCPAcknowledge` produzem os bytes esperados.

### 5.4 DNS (`pkg/dns-simulation-server`, `pkg/dns-simulation-cliente`)
- Server: carregar mapa domínio→IP de um arquivo de zona (`zones.txt`/JSON) ou env; sem stdin.
- Client: domínio e endereço do server por flags/env; modo oneshot p/ o demo.

### 5.5 NTP (`pkg/NTP`) — novo
- Servidor SNTP mínimo (UDP): responde pacote de 48 bytes com timestamps NTP (era 1900).
- Client: envia request, calcula offset/round-trip, imprime hora.
- TDD: teste de conversão `time.Time` ↔ timestamp NTP (segundos + fração).

### 5.6 FTP (`pkg/FTP`)
- Estender `handleConnection`: loop de comandos USER/PASS/PWD/LIST/QUIT; opcional STOR/RETR
  num diretório de volume. Auth via Mongo (`auth.Authenticate` já existe).

### 5.7 Limpeza
- Remover binários commitados: `**/*.exe`, `pkg/dhcp_client/main`, `pkg/dhcp_server/cmd/main`,
  `pkg/session-manager-c/cmake-build-debug/`, `**/*.obj`, `pkg/deprecated/mail_server/smtp_server`.
- Criar `.gitignore` (binários Go, build C, `.idea/`, etc.).
- `.gitmodules`: remover submódulos privados inacessíveis / `deprecated`; documentar como histórico.
- Renomear `pkg/smtp-simulation-server` → `pkg/network-quiz` (e o client correspondente).

## 6. Limitações documentadas (no README)
- **SMTP real:** não implementado; o antigo "smtp-simulation" era um quiz. Roadmap.
- **session-manager-c:** opcional (profile `extras`) por causa do build vcpkg pesado.
- **Pacotes DHCP/DNS:** formato didático simplificado, não conforme RFC.
- **Submódulos C/Go originais:** repositórios privados, fora do escopo do demo.

## 7. Tooling
- **`Makefile`** (raiz): `up`, `down`, `demo`, `logs`, `test`, `clean`, `dev` (modo nativo opcional).
- **`scripts/demo.sh`:** sobe a stack e exercita cada serviço, imprimindo a saída.
- **`scripts/smoke.sh`:** asserções PASS/FAIL por serviço (base para CI).
- **`.env.example`.**

## 8. Testes & verificação
- Go: `go vet` + `go build` + `go test ./...` em cada módulo; testes unitários novos (NTP, DHCP, DNS).
- `docker compose build` de tudo; `docker compose up` + `scripts/smoke.sh`.
- Relatar honestamente o que rodou neste ambiente (a stack Elasticsearch é pesada).

## 9. README (PT + EN)
`README.md` (PT) e `README.en.md` (EN), cada um com:
1. Título, badges, descrição de 1 linha.
2. Diagrama de arquitetura (Mermaid).
3. Matriz de serviços (status / porta / linguagem).
4. Quickstart: pré-requisitos, `docker compose up`, `make demo`.
5. Uso por serviço, com saída de exemplo.
6. Estrutura do projeto.
7. Troubleshooting.
8. Roadmap / limitações conhecidas.
9. Licença + contribuição.

Preservar o contexto do mini-curso de redes (Parte 1/Parte 2) que já existe no README atual.

## 10. Fases de implementação (alto nível)
1. Higiene: `.gitignore`, remover binários, renomear quiz, `.env.example`.
2. Tornar serviços headless + env-driven (DNS, logs, DHCP wiring).
3. Implementar NTP (TDD) e estender FTP.
4. Reescrever `docker-compose.yml` (rede, healthchecks, todos os serviços, profile extras).
5. Corrigir seed Node.
6. Makefile + `scripts/demo.sh` + `scripts/smoke.sh`.
7. Testes Go + verificação `compose build`/`up`.
8. README bilíngue.
