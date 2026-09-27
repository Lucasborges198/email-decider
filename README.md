# Email Decider

`email-decider` is a Go HTTP API that receives email messages and forwards them to a separate Python service powered by Laya. It is a learning project that demonstrates a small service-to-service flow over HTTP.

```text
Client
  │ POST /decision
  ▼
Go API (:8080)
  │ POST /laya/decision
  ▼
Laya API (:9090)
  ▼
Typed email-intent decision
```

## What this project does

- Exposes a `POST /decision` endpoint on port `8080`.
- Loads the Laya service base URL from `LAYA_URL`.
- Forwards the incoming JSON payload to `POST /laya/decision`.
- Relays the upstream HTTP status code and JSON body back to the caller.
- Applies a 10-second timeout to requests sent to the Laya service.

The Python service lives in the companion repository: [`laya-api`](../laya/laya-api).

## Requirements

- Go `1.27.1` or later, as declared in `go.mod`.
- A running `laya-api` instance. Its default local address is `http://localhost:9090`.

## Project structure

```text
.
├── main.go                 # Starts the HTTP server and registers routes
├── handler/handler.go      # Request parsing and Laya HTTP call
├── dto/dto.go              # Request and response data structures
├── .env.example            # Local environment-variable template
├── go.mod                  # Module and dependency declarations
└── go.sum                  # Dependency checksums
```

## Setup

1. Clone the repository and enter the project directory.

2. Create the local environment file:

   ```bash
   cp .env.example .env
   ```

3. Set the Laya service address in `.env`:

   ```dotenv
   LAYA_URL="http://localhost:9090"
   ```

4. Download the Go module dependencies:

   ```bash
   go mod download
   ```

5. Start the API:

   ```bash
   go run .
   ```

The server listens on `http://localhost:8080`.

## Run the complete local flow

Start the Python service first, from the `laya-api` repository:

```bash
uv sync --no-install-project
PYTHONPATH=src .venv/bin/python -m flask \
  --app laya_api.index:app run --port 9090 --no-reload
```

Then start this Go API:

```bash
go run .
```

Send a request through the Go gateway:

```bash
curl --request POST http://localhost:8080/decision \
  --header "Content-Type: application/json" \
  --data '{
    "decisionType": "email_intent",
    "message": "I cannot access my account after changing my password."
  }'
```

The Go API forwards the JSON to Laya and returns the upstream JSON response.

## API contract

### `POST /decision`

Request body:

```json
{
  "decisionType": "email_intent",
  "message": "I would like to schedule a meeting next week."
}
```

| Field | Type | Current behavior |
|---|---|---|
| `decisionType` | string | Part of the DTO and forwarded to Laya; it is not validated or used by the current Python classifier. |
| `message` | string | Required by the Python service and used as the text to classify. |

Successful responses preserve the status code and JSON body sent by `laya-api`.

## Error behavior

| Situation | Current response |
|---|---|
| Invalid JSON received by the Go API | HTTP `500` with a plain-text error message |
| `LAYA_URL` is not configured | HTTP `500` |
| Laya service cannot be reached within 10 seconds | HTTP `502` |
| Laya returns an HTTP response | Its status code and body are forwarded to the client |

## Debugging in VS Code

The repository includes a `Debug Email Decider API` launch configuration. Open the repository root as the VS Code workspace, select that configuration, and press `F5`.

The configuration targets `${workspaceFolder}`, where `go.mod` and `main.go` are located. This avoids accidentally asking Delve to build the `.vscode` directory.

## Development commands

```bash
# Compile the project
go build .

# Run all package tests
go test ./...

# Format Go source files
gofmt -w main.go handler/*.go dto/*.go
```

## Current limitations

- The project has no automated tests yet.
- `decisionType` is not validated by the handler.
- The Go service relays Laya's response instead of mapping it to a stable Go-owned response schema.
- The API has no authentication, persistence, observability, or production deployment configuration.

## Environment and generated files

- Commit `go.mod` and `go.sum`.
- Do not commit `.env`, generated binaries, or Delve debug binaries.
- `.env.example` is safe to commit because it contains only the variable name and a placeholder value.

## License

No license has been selected yet.
