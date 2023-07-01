# Longevity

## Goal

Provide an application for discovering and managing distributed Longevity Digital Twins.

## Prerequisites

- golang >= `1.19`

### optional

- make
- docker

## General

- vendoring is used for dependency management
  - this also creates the file `go.sum`, which is the dependency tree of the overall project
- commits follow the [conventional-commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) specification

## Setup

The Makefile provides a `init` command, that initializes the project, if that isn't already done. The `setup` command enables vendoring and fetches dependencies

## Testing

The `Makefile` provides several helpful commands for testing and coverage. _short_ and _verbose_ flags for the test command can be set via the CLI flags `TEST_VERBOSE` and `TEST_SHORT`. 
The Coverage report can be generated via `make cover` and is located in the `out` directory.
Testing and coverage will always be executed regardless of caching.

## Variables

Environment variables can be set in `.env` file. This includes the location of the unix socket, the meta-repository and the ODM data directory. I would advise to change these values to your liking, depending on your operating system. All values are loaded via [godotenv](https://github.com/joho/godotenv). Do not use existing shell variables in this file, as they are not resolved correctly.

## Building

`make build` produces the binary `orchestration-manager` located in the `out` directory.
`make cli` produces the binary `odm` located in the `out` directory. This is the command line frontend application.

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