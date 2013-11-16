all: format vet lint
format:
	@gofmt -s -w *.go
	@goimports -w *.go
test:
	@sudo `which go` test -v
doc:
	@`which godoc` github.com/caglar10ur/lxc | less
vet:
	`which go` vet .
lint:
	`which golint` .
ctags:
	@ctags -R --languages=c,go
