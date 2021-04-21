.PHONY:
.SILENT:

build:
	go build -o ./.bin/tg-minter2 cmd/service/main.go

run: build
	./.bin/tg-minter2

build-test:
	go build -o ./.bin/tg-minter2-test cmd/test/main.go

run-test: build-test
	./.bin/tg-minter2-test