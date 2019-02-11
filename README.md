# ![foag](https://raw.githubusercontent.com/lnsp/foag/master/docs/logo.png)

foag (pronounced as *[ˈfoːɐ̯k]*) is a Function-as-a-Service platform built upon Docker containers. It was built as a proof-of-concept within an hour and is not suited for production usage.

## Running the daemon

1. Run `go get github.com/lnsp/foag/foagd`
2. Start using `$GOPATH/bin/foagd`, the server should immediately listen on port 8080

Running an instance of `foagd` requires a working local Docker installation. Especially the user running `foagd` must have access to the `docker` command.

## Running the frontend

After starting up the daemon, you can choose to interact with it using either the CLI or the web UI. When you choose the CLI, remember to set the environment variable `FOAG_ENDPOINT` to your daemons endpoint. Same goes for the web UI, before building and serving it, remember to configure the `.env` file to point to your daemon.

## Deploying a function

To deploy a function to a foagd service you have the option to choose between C, Swift, JavaScript and Go. For minimal startup time functions implemented in C, Swift or Go are recommended. Assuming we have implemented our functionality similar to the ones in the `docs/examples` folder (just read and write from standard input and output), we are ready to push and build our function.

```bash
# Set up foagd endpoint first, in this case locally
export FOAG_ENDPOINT=http://localhost:8080
# Assuming you have cloned this repository
cd docs/examples/hello-c
# and now we can deploy!
foag-cli deploy --lang=c main.c
Your ap is now running on http://localhost:8080/trigger/65e64c928f3c6c84e0bcf96fe93d2f05579a2bb47d0e39a7245e1bd310599fba
```

You may follow the build progress by using the `foag-cli builds logs [deployment]` command. After success, you may test your function.

```bash
curl http://localhost:8080/trigger/65e64c928f3c6c84e0bcf96fe93d2f05579a2bb47d0e39a7245e1bd310599fba
hello from C!
```

Congratulations, you just deployed your first function!
