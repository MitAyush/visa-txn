# visa-txn

HTTP API for accounts and card-like transactions, backed by **SQLite**. Schema is applied automatically on startup from embedded SQL.

## Requirements

- **Go 1.25.x** (see `go.mod`)
- **CGO** and a C toolchain when building on your machine (`github.com/mattn/go-sqlite3` needs it). Docker handles this inside the image.
- **enable CGO** using `export CGO_ENABLED=1` and have C compiler "GCC"

## Run locally

```bash
go mod tidy
go run ./cmd/server
```

`go mod tidy` adds missing modules and drops unused ones so your module matches the code. `go run ./cmd/server` compiles and starts the API without leaving a binary in the repo root.

Alternatively, build and run via Make:

```bash
make start
```

The server listens on **port 8080** by default. SQLite file defaults to **`storage/visa.db`**. Create the directory once if it does not exist: `mkdir -p storage`.

### Environment variables

The server reads **`os.Getenv`** only (no `.env` loader in the binary). Use **`.env.example`** as a template: copy to `.env` and export variables in your shell before `go run`, or set them inline as below.

| Variable | Default           | Purpose                                  |
| -------- | ----------------- | ---------------------------------------- |
| `PORT`   | `8080`            | HTTP listen port                         |
| `DB_URL` | `storage/visa.db` | SQLite DSN/path passed to `database/sql` |

Example:

```bash
PORT=3000 DB_URL=/tmp/visa.db go run ./cmd/server
```

## Run with Docker

Build and start a container (detached, port mapped):

```bash
make docker-build
make docker-run
```

Override port or database path:

```bash
docker run --rm -p 3000:3000 -e PORT=3000 -e DB_URL=/app/storage/custom.db visa-txn
```

## API

| Method | Path             | Description                                                                       |
| ------ | ---------------- | --------------------------------------------------------------------------------- |
| `POST` | `/accounts`      | Create account (JSON: `document_number`)                                          |
| `GET`  | `/accounts/{id}` | Get account by numeric `account_id`                                               |
| `POST` | `/transactions`  | Create transaction (optional header `X-Idempotency-Key`; duplicate key → **409**) |

**Operation types** (for `operation_type_id` in transaction body) are defined in code: **1**–**4** (1–3 debit, 4 credit). Request **amount** must be positive; debit types are stored as negative amounts.

Errors are JSON: `{"error":"..."}` with appropriate HTTP status (e.g. 400, 404, 409).

### Try the API with `curl`

Examples use `http://localhost:8080`. Change the host or port if your server differs.

**Create account**

```bash
curl -sS -X POST "http://localhost:8080/accounts" -H 'Content-Type: application/json' -d '{"document_number":"12345678901"}'
```

Duplicate `document_number` returns **409** with `{"error":"account already exists"}`.

**Get account** (replace `1` with the `account_id` from the create response)

```bash
curl -sS "http://localhost:8080/accounts/1"
```

Unknown id → **404**; invalid path id → **400**.

**Create transaction** (use an `account_id` that exists; `operation_type_id` **1**–**4**; `amount` must be positive in the API)

With idempotency: the first successful `POST` with a given `X-Idempotency-Key` creates the transaction. Sending **another** `POST` with the **same** key when a transaction already exists for that key returns **409** with `{"error":"transaction already exists"}` (it does not return the existing transaction).

```bash
curl -sS -X POST "http://localhost:8080/transactions" -H 'Content-Type: application/json' -H 'X-Idempotency-Key: my-key-001' -d '{"account_id":1,"operation_type_id":1,"amount":100.50,"event_date":"2026-04-18T12:00:00Z"}'
```

Without `event_date` (server sets it):

```bash
curl -sS -X POST "http://localhost:8080/transactions" -H 'Content-Type: application/json' -H 'X-Idempotency-Key: my-key-002' -d '{"account_id":1,"operation_type_id":4,"amount":25}'
```

If you omit `X-Idempotency-Key`, the server generates a UUID per request.

Typical errors: **400** (validation / unknown operation type), **404** (account missing), **409** (duplicate `document_number` on account create, or duplicate `X-Idempotency-Key` when a transaction already exists).

**Pretty-print** (optional):

```bash
curl -sS "http://localhost:8080/accounts/1" | jq .
```

## Tests

```bash
make test
```

## Makefile targets

| Target              | Description                                                               |
| ------------------- | ------------------------------------------------------------------------- |
| `make build`        | Build `./visa-txn` binary                                                 |
| `make run`          | Run `./visa-txn`                                                          |
| `make docker-build` | Build Docker image `visa-txn`                                             |
| `make docker-run`   | Run container **`visa-txn-app`** from image **`visa-txn`** on `8080:8080` |
| `make docker-clean` | `docker stop` / `docker rm` container **`visa-txn-app`**                  |
| `make test`         | Run tests with verbose output                                             |
| `make deps`         | `go mod tidy` and `go mod vendor`                                         |
