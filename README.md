# ![foag](https://raw.githubusercontent.com/lnsp/foag/master/docs/logo.png)

foag (pronounced as *[ˈfoːɐ̯k]*) is a Function-as-a-Service platform built upon Docker containers. It was built as a proof-of-concept within an hour and is not suited for production usage.

## Running the daemon

1. Run `go get github.com/lnsp/foag/foagd`
2. Start using `$GOPATH/bin/foagd`, the server should immediately listen on port 8080

Running an instance of `foagd` requires a working local Docker installation. Especially the user running `foagd` must have access to the `docker` command.

## Running the frontend

After starting up the daemon, you can choose to interact with it using either the CLI or the web UI. When you choose the CLI, remember to set the environment variable `FOAG_ENDPOINT` to your daemons endpoint. Same goes for the web UI, before building and serving it, remember to configure the `.env` file to point to your daemon.
