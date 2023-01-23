# Longevity

## Goal

Provide an application for discovering and managing distributed longevity digital twins.

## Prerequisites

- golang >= `1.19`
- a database
  - postgres or sqlite3

### optional

- make
- docker

## General

- vendoring is used for dependency management
  - this also creates the file `go.sum`, which is the dependency tree of the overall project
- commits follow the [conventional-commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) specification

## Setup

The Makefile provides a `init` command, that initializes the project, if that isn't already done. The `setup` command enables vendoring and fetches existing requirements

## Testing

The `Makefile` provides several helpful commands for testing, short testing, verbose testing, and also providing a coverage report located in the `out` directory.

## Building

The command from the `Makefile` produces a binary located in the `out` directory.

## Database

The application requires a database. Currently supported are postgres and sqlite3. I recommend a postgres library, and sqlite3 only for testing purposes.

The postgres database needs to expose port `5432` and is currently accessible via the default user `postgres`. The database password needs to be specified in the database setup method call. 

Easiest way to get the database started is via docker:
`docker run -d --rm -p 5432:5432 -e POSTGRES_PASSWORD=foobar --name longevity-db postgres`
