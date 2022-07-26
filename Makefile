.PHONY: vet test coverage

vet:
	go vet
	golangci-lint run ./...

test: vet
	go test -v -race ./...

coverage: vet
	go test -coverprofile=c.out -v -race ./...
	gocov convert c.out | gocov report