# GOFW
**gofw** (GO! Forward) is a simple TCP traffic forwarder, nothing more, nothing less.
It supports multiple connections to the listening port.

# Build
[golang](https://golang.org/) is required to build this application.

    > make build

This will build the application and moves the binary into the newly created *bin*/ directory.

Now you can use the binary directly or copy it to a directory where it can be accessed all the time (e.g. */usr/local/bin/*).

# Usage examples
Accept TCP connections from localhost on port 1337 and forward them to localhost on port 631.

    gofw -s 1337 -d 631
Accept TCP connections from different hosts in the network on port 1337 and forward them to localhost on port 443

    gofw -s 0.0.0.0:1337 -d 443

Accept TCP connections from different hosts in the network on port 1339 and forward them to the host with the IP 192.168.1.157 on port 5900

    gofw -s 0.0.0.0:1339 -d 192.168.1.157:5900

# But why?
I just learned GO a few days ago and this seemed like an interesting little project.
