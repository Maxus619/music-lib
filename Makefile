build:
	docker-compose build music-lib

run:
	docker-compose up music-lib

test:
	go test -v ./...

migrate:
	migrate -path ./schema -database 'postgres://postgres:test@0.0.0.0:5436/music-lib?sslmode=disable' up

swag:
	swag init -g cmd/main.go