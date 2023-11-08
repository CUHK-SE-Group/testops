# Makefile

HOSTNAME ?= debian  # replace this !!! or pass this !! :   make start-all HOSTNAME=debian
export HOSTNAME

# Define services
SERVICES := openobserve redpanda

# Define commands for each service
.PHONY: $(SERVICES) start stop restart ps logs start-all stop-all restart-all ps-all logs-all

# Start services
start: $(SERVICES)
$(SERVICES):
	@echo "Starting $@..."
	@docker compose -f $@/docker-compose.yaml up -d

# Stop services
stop: $(addprefix stop-,$(SERVICES))
stop-%:
	@echo "Stopping $*..."
	@docker compose -f $*/docker-compose.yaml down

# Restart services
restart: $(addprefix restart-,$(SERVICES))
restart-%:
	@echo "Restarting $*..."
	@docker compose -f $*/docker-compose.yaml down
	@docker compose -f $*/docker-compose.yaml up -d

# Service status
ps: $(addprefix ps-,$(SERVICES))
ps-%:
	@echo "Status of $*..."
	@docker compose -f $*/docker-compose.yaml ps

# Service logs
logs: $(addprefix logs-,$(SERVICES))
logs-%:
	@echo "Logs of $*..."
	@docker compose -f $*/docker-compose.yaml logs -f

# Start all services at once
start-all:
	@echo "Starting all services..."
	@$(foreach service,$(SERVICES),docker compose -f $(service)/docker-compose.yaml up -d;)

# Stop all services at once
stop-all:
	@echo "Stopping all services..."
	@$(foreach service,$(SERVICES),docker compose -f $(service)/docker-compose.yaml down;)

# Restart all services at once
restart-all: stop-all start-all

# Status of all services
ps-all:
	@echo "Status of all services..."
	@$(foreach service,$(SERVICES),docker compose -f $(service)/docker-compose.yaml ps;)

# Logs of all services
logs-all:
	@echo "Fetching logs of all services..."
	@$(foreach service,$(SERVICES),docker compose -f $(service)/docker-compose.yaml logs -f &)


network:
	docker network create redpanda_network