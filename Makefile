.PHONY: start-server
start-server:
	@ air -c .air.server.toml

.PHONY: start-client
start-client:
	@ air -c .air.client.toml

.PHONY: docker-start
docker-start:
	@ docker compose up --remove-orphans
	@ docker compose down --remove-orphans
