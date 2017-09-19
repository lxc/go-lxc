NO_COLOR=\033[0m
OK_COLOR=\033[0;32m

all: format vet lint

format:
	@echo "$(OK_COLOR)==> Formatting the code $(NO_COLOR)"
	@gofmt -s -w *.go
	@goimports -w *.go || true

test:
	@echo "$(OK_COLOR)==> Running go test $(NO_COLOR)"
	@sudo `which go` test -v -x

test-race:
	@echo "$(OK_COLOR)==> Running go test $(NO_COLOR)"
	@sudo `which go` test -race -v -x

test-unprivileged:
	@echo "$(OK_COLOR)==> Running go test for unprivileged user$(NO_COLOR)"
	@`which go` test -v -x

test-unprivileged-race:
	@echo "$(OK_COLOR)==> Running go test for unprivileged user$(NO_COLOR)"
	@`which go` test -race -v -x

cover:
	@sudo `which go` test -v -coverprofile=coverage.out
	@`which go` tool cover -func=coverage.out

doc:
	@`which godoc` gopkg.in/lxc/go-lxc.v2 | less

vet:
	@echo "$(OK_COLOR)==> Running go vet $(NO_COLOR)"
	@`which go` vet .

lint:
	@echo "$(OK_COLOR)==> Running golint $(NO_COLOR)"
	@`which golint` . || true

escape-analysis:
	@go build -gcflags -m

ctags:
	@ctags -R --languages=c,go

setup-test-cgroup:
	for d in /sys/fs/cgroup/*; do \
	    [ -f $$d/cgroup.clone_children ] && echo 1 | sudo tee $$d/cgroup.clone_children; \
	    [ -f $$d/cgroup.use_hierarchy ] && echo 1 | sudo tee $$d/cgroup.use_hierarchy; \
	    sudo mkdir -p $$d/lxc; \
	    sudo chown -R $$USER: $$d/lxc; \
	done

.PHONY: all format test doc vet lint ctags
