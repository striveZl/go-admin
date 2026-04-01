.PHONY: migrate
migrate:
	go run . migrate -d configs -c dev

.PHONY: migrate-docker
migrate-docker:
	go run . migrate -d configs -c docker
