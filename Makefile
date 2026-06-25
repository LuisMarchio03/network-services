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
