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

# Using LXC as an unprivileged user - https://gist.github.com/caglar10ur/8429502
# 
#echo 1 | sudo tee -a /sys/fs/cgroup/memory/memory.use_hierarchy > /dev/null
#
#for entry in /sys/fs/cgroup/*/cgroup.clone_children; do
#    echo 1 | sudo tee -a $entry > /dev/null
#done
#
#for controller in /sys/fs/cgroup/*; do
#    sudo mkdir -p $controller/$USER
#    sudo chown -R $USER $controller/$USER
#    echo $$ > $controller/$USER/tasks
#done

#cat /etc/lxc/lxc-usernet
#caglar veth lxcbr0 1

#cat /etc/subuid
#caglar:100000:100000

#cat /etc/subgid
#caglar:100000:100000

#cat ~/.config/lxc/default.conf
#lxc.network.type = veth
#lxc.network.link = lxcbr0
#lxc.network.flags = up
#
#lxc.id_map = u 0 100000 100000
#lxc.id_map = g 0 100000 100000
test-unprivileged:
	@echo "$(OK_COLOR)==> Running go test for unprivileged user$(NO_COLOR)"
	@`which go` test -v

# requires https://codereview.appspot.com/34680044/
cover:
	@sudo `which go` test -v -covermode=count -coverprofile=coverage.out
	@`which go` tool cover -func=coverage.out

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
