# Order service

Golang microservice to manage orders.

## Prerequisites

Ensure you have the following installed on your system:

- Golang (>=1.24)
- protoc (Protocol Buffers compiler)
- make (latest version recommended)

## Installation

### 1. Clone the project

### 2. Go to the project directory

### 3. Install dependencies

```bash
make i
```

### 4. Customize environment

```bash
cp .env.example .env
```

And setup env vars according to your needs.

## Configuration

```bash
GRPC_PORT=50051        # gRPC server port
LOG_LEVEL=info         # logging severity (debug, info, warn, error)
```

## Running

### Build + run

```bash
make run
```

### Build

```bash
make build
```

### gRPC code generation

```bash
make generate
```
