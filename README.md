# GOFW
**gofw** (GO! Forward) is a simple TCP forwarder for local traffic, nothing more, nothing less.
It supports multiple connections to the listening port.

# Build
[golang](https://golang.org/) is required to build this application.

    > make build

This will build the application and moves the binary into the newly created *bin*/ directory.

Now you can use the binary directly or copy it to a directory where it can be accessed all the time (e.g. */usr/local/bin/*).

# Usage
Scenario: Accept TCP connections from 1337 and forward them to 631.

    gofw -from 1337 -to 631
