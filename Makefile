structs:
	easyjson -lower_camel_case -all http/structs.go
run:
	air
docker-dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
docker-dev-build:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml build
docker-dev-run:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml run api /bin/sh
docker-test:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml run api go test ./...
docker-prod:
	docker-compose up --build
build:
	go build -o tmp/api
