# Architecture Description

This files serves to explain the general architecture of the product store.

## Explaining HTTP Flow

When the programm is run a server is started (per default at localhost:8080) that, depending on the environment variables, sets up either a postgress or in-memory store and a gorilla/mux router.
Available API routes are defined in the handlers [handler.go](/internal/handler/handler.go) and [postgres_handler.go](./internal/handler/postgres_handler.go).

When a http request hits the server and matches one of the defined routes, the matched route dispatches one of the methods defined in the handlers. Handlers deserialize the json request, validate the contents (e.g. name is non-empty) and then calls a store method defined in one of the [stores](./internal/store/)

Depending on which store is used results from the store method are either stored and retrieved from the postgres data base or simply stored and retrieved from an in memory Product map.

The return value of the store methods is then returned as a JSON response by the handlers along with a status code.

## Difference between postgres and memory store

Depending on which way the application is started `go run ./cmd/api` or `docker compose up --build` the application will use either the in-memory store or the postgres store. The in-memory store is a simple implementation that uses a map to store products, while the postgres store interacts with a PostgreSQL database to persist product data. Essentially the in-memory store will be empty everytime the application is started, while the postgres store will persist data between sessions.