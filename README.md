# Go Bindings for LXC (Linux Containers)

This package implements [Go](http://golang.org) bindings for the [LXC](http://linuxcontainers.org/) C API.

## Requirements

This package requires [LXC 1.0+](https://github.com/lxc/lxc/releases) and [Go 1.x](https://code.google.com/p/go/downloads/list).

It has been tested on

+ Ubuntu 13.04 (raring) by using distribution [provided packages](https://launchpad.net/ubuntu/raring/+package/lxc)
+ Ubuntu 13.10 (saucy) by using distribution [provided packages](https://launchpad.net/ubuntu/saucy/+package/lxc)
+ Ubuntu 14.04 (trusty) by using distribution [provided packages](https://launchpad.net/ubuntu/trusty/+package/lxc)

## Installing

The typical `go get github.com/lxc/go-lxc` will install LXC Go Bindings.

## Documentation

Documentation can be found at [GoDoc](http://godoc.org/github.com/lxc/go-lxc)

## Examples

See the [examples](https://github.com/lxc/go-lxc/tree/master/examples) directory for some.

## Contributing

We'd love to see go-lxc improve. To contribute to go-lxc;

* **Fork** the repository
* **Modify** your fork
* Ensure your fork **passes all tests**
* **Send** a pull request
	* Bonus points if the pull request includes *what* you changed, *why* you changed it, and *tests* attached.
	* For the love of all that is holy, please use `go fmt` *before* you send the pull request.

We'll review it and merge it in if it's appropriate.
