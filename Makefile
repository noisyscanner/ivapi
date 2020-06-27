structs:
	easyjson -lower_camel_case -all http/structs.go
run:
	air
docker-dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
docker-prod:
	docker-compose up --build
build:
	go build -o tmp/api
