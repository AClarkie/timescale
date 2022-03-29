# Timescale exercise

## Prerequisites

1. Install Docker, if you don't already have it. For packages and instructions, see the [Docker installation documentation](https://docs.docker.com/get-docker/)
2. Install go, if you don't already have it. For versions and instructions, see the [Download and install documentation](https://go.dev/doc/install)
## Setup

1. Setup the database
 
    Run the following makefile target:
    * `make start-database`

2. Build the application

    Run the following makefile target to generate a binary name `app`
    * `make build`

## Running the application

Run the application in your terminal as below:
```bash
./app -queryParams query_params.csv -goroutineCount 2
```

There are additional flags you can pass to the application, the full list is below:

| Name | Default | Type | Description |
| ---- | ------- | ---- | ----------- |
| queryParams | query_params.csv | string | Path to the input csv |
| verbose | false | bool | Enable verbose logging |
| goroutineCount | 2 | int | The number of goroutines to use |
| dbHost | localhost | string | The database host |
| dbName | homework | string | The database name |
| dbUser | postgres | string | The database user |
| dbPassword | password | string | The database password |
| dbSSLMode | disable | string | The database sslmode |


## Running the tests
Run the following makefile target to runs the tests:
```
make test
```

## Clean up

Run the following makefile target to stop the database:
```
make stop-database
```