# Go (golang) Bindings for LXC (Linux Containers)

This package implements [Go](http://golang.org) (golang) bindings for the [LXC](http://lxc.sourceforge.net/) C API.

## Installing

This package requires [LXC 0.9](http://lxc.git.sourceforge.net/git/gitweb.cgi?p=lxc/lxc;a=summary) or newer versions.

It has been tested on Ubuntu 12.10 (quantal) by manually installing LXC 0.9 or on Ubuntu 13.04 (raring) by using distribution [provided packages](https://launchpad.net/ubuntu/raring/+package/lxc).

    go get github.com/caglar10ur/lxc

## Documentation

Checkout the documentation at [GoDoc](http://godoc.org/github.com/caglar10ur/lxc)

## Example

Checkout the [examples](https://github.com/caglar10ur/lxc/tree/master/examples) directory.

## Note

Note that as we don’t have full user namespaces support at the moment, any code using the LXC API needs to run as root.
