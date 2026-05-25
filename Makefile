.PHONY: app-build
app-build:
	docker compose up --build app

.PHONY: migrate-docker
migrate-docker:
	docker compose run --rm app migrate -d /app/configs -c dev

.PHONY: fmt
fmt:
	gofmt -w .
