check: test lint coverage

test:
	go test -race ./...

lint:
	golangci-lint run

coverage:
	go-carpet -mincov 99.9
