run:
	@echo "> Starting server"
	go run server.go

test:
	@echo "> Running tests"
	go test ./... -v

coverage:
	@echo "> Checking coverage"
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

format:
	@echo "> Formatting the source"
	gofmt -d -e

build:
	@echo "> Building binary"
	go build -o bin/main .

clean:
	@echo "> Cleaning build cache and temporary files generated from tests"
	go clean
	rm -rf tmp bin *.rdb *.trdb *.out

.PHONY: run test coverage format clean build