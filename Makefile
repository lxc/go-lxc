NO_COLOR=\033[0m
OK_COLOR=\033[0;32m

all: format vet lint

format:
	@echo "$(OK_COLOR)==> Formatting the code $(NO_COLOR)"
	@gofmt -s -w *.go
	@goimports -w *.go

test:
	@echo "$(OK_COLOR)==> Running go test $(NO_COLOR)"
	@sudo `which go` test -v

# requires https://codereview.appspot.com/34680044/
cover:
	@sudo `which go` test -v -covermode=count -coverprofile=coverage.out

doc:
	@`which godoc` github.com/caglar10ur/lxc | less

vet:
	@echo "$(OK_COLOR)==> Running go vet $(NO_COLOR)"
	@`which go` vet .

lint:
	@echo "$(OK_COLOR)==> Running golint $(NO_COLOR)"
	@`which golint` .

ctags:
	@ctags -R --languages=c,go

.PHONY: all format test doc vet lint ctags
