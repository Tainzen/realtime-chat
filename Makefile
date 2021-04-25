build:
	CGO_ENABLED=0 go build -o bin/server

swagger:
	swag init -g main.go

run:
	./bin/server