# URL Shortener

An easily configurable, locally host-able URL shortener written in Golang. A (small) weekend project.

## Table of Contents
 - [Requirements](#requirements)
 - [Configuration](#configuration)
 - [Deployment](#deployment)
 - [Database Management](#database)
 - [Usage](#usage)

## Requirements

To run and deploy this application, you will need a few things:
 - [Postgres](https://hub.docker.com/_/postgres) 
 - [Go 1.18+](https://go.dev/doc/install)
 - [Docker](https://docs.docker.com/get-docker/)
 - [Docker-compose](https://docs.docker.com/compose/install/)

## Configuration

| Name | Environment Variable | Values | Description |
| ---- | -------------------- | ------ | ----------- |
| encoding | `URL_LOGGER_ENCODING` | `console` or `json` | The encoding of the logger. |
| level | `URL_LOGGER_LEVEL` | `info`, `debug`, `warning`, `error`, `dpanic`, `panic`, and `fatal` | The level at which logs are recorded. | 
| REST host | `URL_REST_HOST` | `localhost` or `0.0.0.0` | The hostname of the REST server, should be `localhost` outside of container, and `0.0.0.0` inside a container. |
| REST port | `URL_REST_PORT` | `8080` | The port of the REST server. |
| DB Host | `URL_DATABASE_HOST` | `localhost` or `database` | The host of the PostgreSQL instance, should be `localhost` outside of a container, and `database` inside a container. |
| DB Port | `URL_DATABASE_PORT` | `5432` | The port of the PostgreSQL instance, this is the default port. |
| DB User | `URL_DATABASE_USER` | `postgres` | The user of the database, this is the default user. |
| DB Password | `URL_DATABASE_PASSWORD` | `postgres` | The password to the user, this can be changed to whatever value is required. |
| DB Name | `URL_DATABASE_NAME` | `urlshortener`| The name of the database within the PostgreSQL instance that the application uses. |
| DB Capacity | `URL_DATABASE_CAPACITY` | `10` |The complexity of generated URLs. Directly impacts the length of generated ID's. |

## Deployment

This project makes use of [docker-compose](https://docs.docker.com/compose/install/), so the following commands can be used:

 * To build the containers
 ```bash
 docker-compose build --no-cache
 ```

 * To run the containers
 ```bash
 docker-compose up -d
 ```

 * To bring the containers down
 ```bash
 docker-compose down
 ```

 * To view the logs of the containers
 ```bash
 docker-compose logs -f [container-name]
 ```

## Database

This project uses PostgreSQL for data persistence. It saves any data from the instance to a `./db/` directory placed in the project directory.

 * To create the database within the PostgreSQL instance:
 ```bash
 make db-create
 ```

 * Then, to populate the database with the `urls` table:
 ```bash
 make db-init
 ```

 * To wipe all data from the database:
 ```bash
 make db-wipe
 ```

 * Additionally, you can exec into the container using:
 ```bash
 make db-exec
 ```

 There is a known issue where the service is running, and the database is created for the first time, the application may not recognise the database as created. In this case, restart the application container. This can be done with:
```bash
docker-compose restart rest-server
```

## Usage

It is what it says on the tin -- it's a basic URL shortener. Links have an expiry time added to them by default, but this can be altered easily [here](./internal/datastore/postgres.go). I may add some more customisability in the future.

For  both of the below scenarios, we assume that the application is running on `localhost:8080`.

### Create a new short URL

URL: `POST http://localhost:8080/api/v1/shortener/`

Request:
```json
{
    "url":"https://github.com/alexmerren"
}
```
Response:
```json
{
    "id": "WTWBY9Rr5C"
}
```

### Retrieve a full URL

URL: `GET http://localhost:8080/api/v1/shortener/WTWBY9Rr5C/`

Request:
```json
// This is deliberately blank...
```

Response:
```json
{
    "url": "https://github.com/alexmerren"
}
```
