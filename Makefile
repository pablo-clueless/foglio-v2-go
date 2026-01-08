dev:
	air

build:
	go build -o foglio-v2 main.go

test:
	go test ./...

test-unit:
	go test -v ./src/...

test-e2e:
	go test -v ./tests/...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-ci:
	cp .env.test .env
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f foglio-v2 coverage.out coverage.html

docker-test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down

.PHONY: dev build test test-unit test-e2e test-coverage test-ci lint clean docker-test