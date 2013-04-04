all:
	go build create.go
	go build destroy.go
	go build start.go
	go build stop.go
	go build freeze.go
	go build unfreeze.go
	go build shutdown.go
	go build stats.go
	go build list.go
	go build concurrent_create.go
	go build concurrent_destroy.go
clean:
	rm -f create destroy start stop freeze unfreeze shutdown stats list concurrent_create concurrent_destroy
format:
	gofmt -s -w *.go
