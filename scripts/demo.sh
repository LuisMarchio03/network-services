#!/usr/bin/env bash
set -euo pipefail

COMPOSE="docker compose -f StartNewServer/docker-compose.yml"

tcp_echo() {
  local port="${TCP_ECHO_PORT:-8888}"
  timeout 4 bash -c '
    p="$1"
    exec 3<>"/dev/tcp/localhost/$p" || exit 1
    printf "ola\n" >&3
    IFS= read -r -t 3 out <&3
    printf "%s\n" "$out"
  ' _ "$port" 2>/dev/null
}

echo "==> Subindo a stack (pode demorar no primeiro build)..."
make up

echo "==> Aguardando Elasticsearch ficar saudável..."
until curl -sf "http://localhost:${ES_PORT:-9200}/_cluster/health" >/dev/null; do
  echo "    ... aguardando ES"
  sleep 3
done

echo "==> [DNS] resolvendo example.com"
$COMPOSE --profile clients run --rm dns-client || true

echo "==> [NTP] consultando hora do servidor"
$COMPOSE --profile clients run --rm ntp-client || true

echo "==> [TCP echo] enviando 'ola'"
tcp_echo || true

echo "==> [Logs] documentos no índice 'logs' (aguardando 6s para acumular)"
sleep 6
curl -s "http://localhost:${ES_PORT:-9200}/logs/_count" || true
echo

echo "==> [DHCP] últimos eventos de negociação:"
$COMPOSE logs --tail=10 dhcp-server dhcp-client || true

echo
echo "==> Demo concluída."
echo "    Kibana:   http://localhost:${KIBANA_PORT:-5601}"
echo "    RabbitMQ: http://localhost:${RABBITMQ_MGMT_PORT:-15672}"
