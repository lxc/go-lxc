all: format vet lint
format:
	@gofmt -s -w *.go
test:
	@sudo `which go` test -v
doc:
	@`which godoc` github.com/caglar10ur/lxc | less
vet:
	`which go` vet .
lint:
	`which golint` .
