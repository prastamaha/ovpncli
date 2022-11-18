build:
	go build -o bin/ovpncli .

run:
	go run .

compile:
	GOOS=linux GOARCH=amd64 go build -o bin/ovpncli-linux-amd64 .
