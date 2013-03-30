format:
	gofmt -s -w *.go
test:
	sudo `which go` test -v
testmem:
	sudo `which go` test -run TestGetMemoryUsageInBytes -v
