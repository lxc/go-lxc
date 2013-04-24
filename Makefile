format:
	gofmt -s -w *.go
test:
	sudo `which go` test -v
test_clone:
	sudo `which go` test -run TestClone -v
test_concurrent:
	sudo `which go` test -run TestConcurrentDefined_Negative -v
	sudo `which go` test -run TestConcurrentCreate -v
	sudo `which go` test -run TestConcurrentDefined_Positive -v
	sudo `which go` test -run TestConcurrentDestroy -v
