# Shell (Go)

A Unix shell implemented from scratch in Go following the Codecrafters
"Build Your Own Shell" challenge.

The project focuses on understanding how a shell works internally:
command parsing, process execution, terminal interaction, environment
handling, and job control.

## Implemented Features

## Core Shell

- Interactive REPL
- Command prompt
- Exit handling
- Invalid command handling
- Execute external programs
- PATH executable lookup
- TAB to autocomplete builtins, external binaries, paths, file names
- Builtins:
  - `echo`
  - `type`
  - `pwd`
  - `cd`
  - `echo`
  - `complete`
  - `history`
  - `jobs`
  - `declare`
  - `exit`


## Quoting & Escaping

- Single quotes
- Double quotes
- Backslash escaping
- Quoted executable execution

## Redirection

- Redirect stdout
- Redirect stderr
- Append stdout
- Append stderr

Examples:

```bash
echo hello > file.txt
cat file.txt
echo error 2> errors.txt
echo hello >> file.txt
```

## Command Completion

- Builtin completion
- Executable completion
- Partial completion
- Multiple matches
- Filename completion
- Directory completion

## Programmable Completion

- `complete` builtin
- Register completion handlers
- Environment-aware completion
- Multiple completion candidates
- Longest common prefix

## Background Jobs

- Run commands in background

```bash
sleep 10 &
```

- `jobs` builtin
- Job tracking
- Process cleanup
- Job number reuse

## History

- `history` builtin
- Limit history output

```bash
history 10
```

- Arrow key navigation
- Execute previous commands
- Persistent history

Supported:

```bash
history -r file
history -w file
history -a file
```

## History Persistence

- Read history from file
- Write history to file
- Append history
- Load history using `HISTFILE`
- Save history on exit

## Shell Variables

- `declare` builtin
- Variable storage
- Identifier validation
- Parameter expansion

Examples:

```bash
declare NAME=mk
echo $NAME
echo ${NAME}123
```

## In Progress

## Pipelines

Currently working on:

```bash
cat file | wc
```

Planned:

- Dual command pipelines
- Builtins inside pipelines
- Multi-command pipelines

## Tech Stack

- Go
- Standard library only

Main packages used:

- `os/exec` — process execution
- `os` — filesystem and environment
- `io` — streams and pipes
- terminal input handling

## Purpose

This project is a learning implementation of shell internals:

- How processes are created
- How stdin/stdout/stderr are connected
- How terminals handle input
- How shells manage state
- How command execution pipelines work
