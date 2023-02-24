lint:
	golangci-lint --color=always run ./... -v

test:
	go test -v ./command