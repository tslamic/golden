FROM golangci/golangci-lint:latest
WORKDIR /golden
COPY ./ .
RUN make test
