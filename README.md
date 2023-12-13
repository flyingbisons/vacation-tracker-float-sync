# Vacation Tracker & Float Integration

## Overview

This is a Cloudflare Worker that integrates with Float and Vacation Tracker to automatically update Float with vacation time.
The Intention of this tool is to help our managers and team leads to not have to manually update Float with vacation time.

Integration is made in Go and compiled to WebAssembly thanks to [syumai/workers](https://github.com/syumai/workers).
Looking forward to your thoughts and suggestions.

## How to use
To use this tool, you will need to create a Vacation Tracker account with API access and a Float account. 
You will also need to create a Cloudflare account and a Cloudflare Workers D1 database.

## How it works
Because of Vacation Tracker is not providing webhooks, we need to use a workaround to get the data from Vacation Tracker.
This is done by using a Cloudflare Worker to scrape the data from the Vacation Tracker API and then send it to Float.
Scrapping is done by using a Go Wasm binary that is compiled to WebAssembly and then run in the Cloudflare Worker.
To store processed leave requests, we use a Cloudflare D1 database.

You can configure frequency of the scrapping by changing the cron expression, check `wrangler.toml` file. 
There is also an option to run the scrapping manually by sending a request to the worker endpoint.

## Usage
Copy `wrangler.toml.example` to `wrangler.toml` and fill in the blanks.

- `main.go` includes simple Cron task and HTTP server handler (for dev).
- update const `internal/flaot/float.go` with your Float TimeTypes identifiers.
- update const `internal/vacation/client.go` with your Vacation Tracker leave request types identifiers.

## Limitations
There is no easy option to get actual time in Cloudflare Workers, so we need to use a workaround, I have used api call to http://worldtimeapi.org.

With webassembly, there is a problem with unit testing, even adding `GOOS=js GOARCH=wasm` in fornt of test command  produce `exec format error`. So I could not write unit tests for this project (I'd love to hear your comments on how to improve it).

## Development

The Core of a project is based on [syumai/workers](https://github.com/syumai/workers), plase check this repo for more details (specially examples directory).

You can query your database using command:

for production:
```shell
npx wrangler d1 execute [db_name] --command="SELECT count(*) FROM Requests" 
```

for local:
```shell
npx wrangler d1 execute [db_name] --command="SELECT count(*) FROM Requests"  --local
```

### Commands

```
make dev     # run dev server (plase edit wrangler.toml file and main.go)
make deploy  # deploy worker
make test    # run tests
make init-db # init database (production)
make init-db-init-db-local # init database (locally)
```

For local development, is easier to use http endpoint. Check commended code in `main.go` file.


