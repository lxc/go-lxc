# Go Bindings for LXC (Linux Containers)

This package implements [Go](http://golang.org) bindings for the [LXC](http://lxc.sourceforge.net/) C API.

## Installing

This package requires [LXC 0.9](http://lxc.git.sourceforge.net/git/gitweb.cgi?p=lxc/lxc;a=summary) and [Go 1.x](https://code.google.com/p/go/downloads/list).

It has been tested on Ubuntu 12.10 (quantal) by manually installing LXC 0.9 or on Ubuntu 13.04 (raring) by using distribution [provided packages](https://launchpad.net/ubuntu/raring/+package/lxc).

    go get github.com/caglar10ur/lxc

## Documentation

Checkout the documentation at [GoDoc](http://godoc.org/github.com/caglar10ur/lxc)

## Example

See the [examples](https://github.com/caglar10ur/lxc/tree/master/examples) directory.

## Note

Note that as we don’t have full user namespaces support at the moment, any code using the LXC API needs to run as root.

## Development branch

If you are interested with upcoming LXC version (staging tree) then please use the devel branch.
