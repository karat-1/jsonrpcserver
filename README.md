# JRPC2 Server

This repository contains a simple JSON-RPC 2.0 server implementation in Go using the [jrpc2](https://pkg.go.dev/github.com/creachadair/jrpc2) library. The server provides several methods for string operations and alert handling.

## Features

- **CountString**: Counts the number of Unicode code points in a given string.
- **Status**: Returns the status of the server.
- **Alert**: Logs an alert message.

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/karat-1/jrpc2server.git
   cd jrpc2server
   ```

2. Build the project:
   ```sh
   go build
   ```

## Usage

1. Start the server with the default address `localhost:8080` and a maximum of 1000 concurrent tasks:
   ```sh
   ./jrpc2server
   ```

2. Customize the address and maximum concurrent tasks using command-line flags:
   ```sh
   ./jrpc2server -address="localhost:8080" -max=500
   ```

## Graceful Shutdown

The server can be gracefully shut down by sending an interrupt signal (Ctrl+C). It will log the total number of received messages before exiting.

## Logging

The server logs various events including startup, received alerts, and shutdown messages. All logs are printed to the standard output.

## Dependencies

- [jrpc2](https://pkg.go.dev/github.com/creachadair/jrpc2)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [creachadair/jrpc2](https://github.com/creachadair/jrpc2) for the JSON-RPC 2.0 library.

## Author

- [karat-1](https://github.com/karat-1)
