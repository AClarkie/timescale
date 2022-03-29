# Timescale exercise

### Prerequisites

1. Install Docker, if you don't already have it. For packages and instructions, see the [Docker installation documentation](https://docs.docker.com/get-docker/)
### Setup

1. Setup the database
  Run the following makefile target:
  * `make start-database`

  This will build the database docker container and start the database, exposing port 5432 locally.

2. 