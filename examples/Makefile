all:
	go build concurrent_create.go
	go build concurrent_destroy.go
	go build concurrent_start.go
	go build concurrent_stop.go
	go build concurrent_shutdown.go
	go build create.go
	go build destroy.go
	go build freeze.go
	go build list.go
	go build shutdown.go
	go build start.go
	go build stats.go
	go build stop.go
	go build unfreeze.go
clean:
	rm -f concurrent_create
	rm -f concurrent_destroy
	rm -f concurrent_start
	rm -f concurrent_stop
	rm -f concurrent_shutdown
	rm -f create
	rm -f destroy
	rm -f freeze
	rm -f list
	rm -f shutdown
	rm -f start
	rm -f stats
	rm -f stop
	rm -f unfreeze
format:
	gofmt -s -w *.go
