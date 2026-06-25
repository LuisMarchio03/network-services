#!/usr/bin/env bash

fail=0

check() {
  local label="$1"
  if eval "$2" >/dev/null 2>&1; then
    echo "PASS: $label"
  else
    echo "FAIL: $label"
    fail=1
  fi
}

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

check "Elasticsearch responde" \
  "curl -sf http://localhost:${ES_PORT:-9200}/_cluster/health"

check "RabbitMQ mgmt responde" \
  "curl -sf http://localhost:${RABBITMQ_MGMT_PORT:-15672}"

check "TCP echo ecoa" \
  'out=$(tcp_echo); printf "%s" "$out" | grep -q ola'

check "Índice de logs tem docs" \
  "[ \"\$(curl -s http://localhost:${ES_PORT:-9200}/logs/_count | grep -o '\"count\":[0-9]*' | cut -d: -f2)\" -gt 0 ]"

exit $fail
