.PHONY: build test lint clean release

build:
	go build -o bin/afk-x ./cmd/afk-x

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf bin/

release:
	goreleaser release --clean

run: build
ifeq ($(shell uname),Darwin)
	script -q /dev/null ./bin/afk-x
else
	script -q -c "./bin/afk-x" /dev/null
endif