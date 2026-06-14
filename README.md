# net-cat

A concurrent TCP group chat server built in Go, recreating the functionality of the `nc` (NetCat) command in a server-client architecture.

## Overview

This project implements a group chat system where one server can handle multiple client connections over TCP. It uses Go routines, channels, and mutexes to manage concurrency safely.

## Features

- **TCP connections** — one server handles up to 10 simultaneous clients
- **Name registration** — every client must choose a non-empty, unique name on join
- **Message broadcasting** — clients send messages to the group; empty messages are ignored
- **Timestamps & identification** — every message is formatted as `[YYYY-MM-DD HH:MM:SS][username]: message`
- **Chat history** — new clients receive all previous messages on connection
- **Join/leave notifications** — server announces when clients join or leave the chat
- **Name changes** — clients can use `/name <newname>` to rename themselves (bonus feature)
- **File logging** — all messages are appended to `chat.log` for persistence (bonus feature)

## Project Structure

```
net/
├── cmd/
│   └── server/
│       └── main.go       # Entry point — loads config and starts the server
├── internal/
│   └── chat/
│       ├── ascii.go      # Welcome banner (ASCII penguin logo)
│       ├── client.go     # Client connection handling, read/write loops, commands
│       ├── message.go    # Message formatting (timestamps, system messages)
│       └── server.go     # Server state, connection handling, broadcasting
├── pkg/
│   └── config/
│       └── config.go     # Command-line argument parsing and validation
├── Makefile              # Build, test, run, and clean targets
├── go.mod                # Go module definition
└── chat.log              # Chat log output (generated at runtime)
```

## Usage

Run the server:

```bash
make run               # listens on default port 8989
make run port=2525     # listens on port 2525
```

Build the binary:

```bash
make build
./net-cat 2525
```

Run tests:

```bash
make test
```

Connect as a client using NetCat (or any TCP client):

```bash
nc localhost 8989
```

You will see the welcome banner and be prompted to enter your name:

```text
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
```

## Client Commands

| Command | Description |
|---------|-------------|
| `/name <name>` | Change your display name |

## Error Handling

- Starting with no arguments uses the default port (8989).
- Starting with more than one argument prints: `[USAGE]: ./TCPChat $port`
- The port must be a number between 1 and 65535.
- The server rejects connections when the chat is full (max 10 clients).
- Duplicate names are rejected and the client is disconnected.

## Implementation Notes

- **Concurrency** — each client connection runs in its own goroutine for reading and writing; a dedicated broadcaster goroutine forwards messages
- **Synchronization** — a `sync.Mutex` protects the shared clients map and message history
- **Channels** — each client has a buffered message channel (`chan string`, capacity 100) used as an inbox
- **Max clients** — enforced at 10 via `maxClients` constant
- **Allowed packages** — `fmt`, `log`, `os`, `net`, `sync`, `time`, `bufio`, `errors`, `strings` (plus `reflect`, `io`)
