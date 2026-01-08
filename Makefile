dev:
	air

build:
	go build -o foglio-v2 main.go

test:
	go test ./...