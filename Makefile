COMPOSE_PROD=docker compose -f docker-compose.prod.yml

prod-up:
	$(COMPOSE_PROD) up -d --build

prod-down:
	$(COMPOSE_PROD) down

prod-logs:
	$(COMPOSE_PROD) logs -f

prod-ps:
	$(COMPOSE_PROD) ps

prod-restart:
	$(COMPOSE_PROD) restart

prod-migrate-up:
	$(COMPOSE_PROD) --profile tools run --rm migrate

prod-db-shell:
	$(COMPOSE_PROD) exec postgres psql -U postgres -d relayops

prod-api-logs:
	$(COMPOSE_PROD) logs -f api

prod-worker-logs:
	$(COMPOSE_PROD) logs -f worker

prod-nginx-logs:
	$(COMPOSE_PROD) logs -f nginx