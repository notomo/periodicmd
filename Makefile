test:
	go test -v ./...

lint:
	go vet ./...

start:
	go run main.go -config ./example_config.json -start-date 2023-12-29

help:
	go run main.go -h
