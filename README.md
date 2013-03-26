# Go (golang) Bindings for LXC (Linux Containers)

This package implements [Go](http://golang.org) (golang) bindings for the [LXC](http://lxc.sourceforge.net/) C API.

## Installing

lxc only tested on Ubuntu 12.10 with following packages

```
ii  liblxc0                                   0.8.0~rc1-4ubuntu39.12.10.2                       amd64        Linux Containers userspace tools (library)
ii  lxc                                       0.8.0~rc1-4ubuntu39.12.10.2                       amd64        Linux Containers userspace tools
ii  lxc-dev                                   0.8.0~rc1-4ubuntu39.12.10.2                       amd64        Linux Containers userspace tools (development)
```

    go get github.com/caglar10ur/lxc

## Example

See [lxc_test.go](https://github.com/caglar10ur/lxc/blob/master/lxc_test.go) as an example.

## TODO

    Implement missing API...
