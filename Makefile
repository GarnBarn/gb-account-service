run:
	go run .

compose-up:
	COMPOSE_PROJECT_NAME=garnbarn-account docker-compose up -d

compose-down:
	COMPOSE_PROJECT_NAME=garnbarn-account docker-compose down