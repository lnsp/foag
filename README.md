# foag

foag (pronounced as *[ˈfoːɐ̯k]*) is a Function-as-a-Service platform built upon Docker containers. It was built as a proof-of-concept within an hour and is not suited for production usage.

## Running the daemon

1. Run `go get github.com/lnsp/foag/foagd`
2. Start using `$GOPATH/bin/foagd`, the server should immediately listen on port 8080

## Running the web frontend

Running an instance of `foagd` requires a working local Docker installation. Especially the user running `foagd` must have access to the `docker` command.