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

The `Makefile` provides several helpful commands for testing and coverage. _short_ and _verbose_ flags for the test command can be set via the CLI flags `TEST_VERBOSE` and `TEST_SHORT`. 
The Coverage report can be generated via `make cover` and is located in the `out` directory.
Testing and coverage will always be executed regardless of caching.

## Building

The command from the `Makefile` produces a binary located in the `out` directory.

## Database

The application requires a database. Currently supported are postgres and sqlite3. I recommend a postgres library, and sqlite3 only for testing purposes.

The postgres database needs to expose port `5432` and is currently accessible via the default user `postgres`. The database password needs to be specified in the database setup method call. 

Easiest way to get the database started is via docker:
`docker run -d --rm -p 5432:5432 -e POSTGRES_PASSWORD=foobar --name longevity-db postgres`

## Commands

This section lists useful commands which help during execution and debugging of the project.

### Stdout / Stderr from LDT

In case the ODM does not redirect the output during process creation, the output of LDTs can be accessed like this:

Getting stdout from an LDT on macOS:
```bash
dtrace -p $LDT_PID -qn 'syscall::write*:entry /pid ==  && (arg0 == 1 || arg0 == 2)/ { printf(%s, copyinstr(arg1, arg2)); }'
```

Getting stdout from an LDT on Linux:
1: stdout
2: stderr
```bash
tail -f /proc/$LDT_PID/fd/1
```