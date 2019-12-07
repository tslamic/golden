.PHONY: vet
vet:
	@golangci-lint \
		-E golint \
		-E gofmt \
		-E goimports \
		-E bodyclose \
		-E gosec \
		-E unconvert \
		-E misspell \
		-E whitespace \
		run ./...

.PHONY: test
test: vet
	go test -coverprofile=c.out -v -race ./...
	gocov convert c.out | gocov report
