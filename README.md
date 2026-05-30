# rprinter

`rprinter` is a lightweight Go-based client daemon that polls and prints receipt jobs from a central gRPC-based queueing service (`printqueue`).

## Features

- **gRPC Registration**: Connects to the remote print service (`print.brotherlogic-backend.com`) and registers the host machine as a receipt printer destination (`DESTINATION_RECEIPT`).
- **Job Processing**: Fetches queued print jobs containing raw lines of text.
- **Local Spooling**: Spools jobs to a local temporary file and prints them using the system's `lp` command line utility.
- **Job Acknowledgment**: Confirms successful prints back to the gRPC backend.

## Architecture

1. **Register**: Sends a `RegisterPrinterRequest` using the local hostname to identify the print node.
2. **Execute**: Automatically downloads any pending print jobs.
3. **Spool & Print**: For each job, writes the content lines to a temporary file under the system temp directory and runs:
   ```bash
   lp /tmp/printdetailsXXXXXX
   ```
4. **Acknowledge**: Signals a successful print by sending an acknowledgment (`Ack`) back to the print queue service.

## Prerequisites

- **Go**: Version 1.19 or higher.
- **CUPS / `lp`**: The system must have the `lp` command line printing utility installed and configured with a default receipt printer.

## Installation & Running

Clone the repository and build or run the program directly:

```bash
# Build the binary
go build -o rprinter main.go

# Run the printer agent
./rprinter
```

## Configuration

The printer client is currently configured to connect to:
- **Server Address**: `print.brotherlogic-backend.com:80`
- **Destination Type**: `DESTINATION_RECEIPT`
- **Credentials**: Insecure/Plaintext gRPC connection.
