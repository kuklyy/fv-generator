run:
	go build -o ./bin/main cmd/fv-generator/main.go && ./bin/main

test:
	go test ./...