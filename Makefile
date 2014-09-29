NO_COLOR=\033[0m
OK_COLOR=\033[0;32m

all: format vet lint

format:
	@echo "$(OK_COLOR)==> Formatting the code $(NO_COLOR)"
	@gofmt -s -w *.go
	@goimports -w *.go || true

test:
	@echo "$(OK_COLOR)==> Running go test $(NO_COLOR)"
	@sudo `which go` test -v

test-race:
	@echo "$(OK_COLOR)==> Running go test $(NO_COLOR)"
	@sudo `which go` test -race -v

test-unprivileged:
	@echo "$(OK_COLOR)==> Running go test for unprivileged user$(NO_COLOR)"
	@`which go` test -v

test-unprivileged-race:
	@echo "$(OK_COLOR)==> Running go test for unprivileged user$(NO_COLOR)"
	@`which go` test -race -v

# requires https://codereview.appspot.com/34680044/
cover:
	@sudo `which go` test -v -coverprofile=coverage.out
	@`which go` tool cover -func=coverage.out

doc:
	@`which godoc` github.com/lxc/go-lxc | less

vet:
	@echo "$(OK_COLOR)==> Running go vet $(NO_COLOR)"
	@`which go` vet .

lint:
	@echo "$(OK_COLOR)==> Running golint $(NO_COLOR)"
	@`which golint` . || true

ctags:
	@ctags -R --languages=c,go

.PHONY: all format test doc vet lint ctags
