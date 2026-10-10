# Go Payment Gateway

A Go-based payment gateway project built to grow with additional payment
features and services, backed by PostgreSQL.

## Requirements

- Go 1.25.4 or later
- Docker Compose

## Run locally

1. Start PostgreSQL:

   ```sh
   docker compose up -d postgres
   ```

   Compose publishes the database on host port `5435` and stores its data in
   the `postgres_data` volume.

2. Create or update a `.env` file in the repository root with settings
   matching `docker-compose.yml`:

   ```dotenv
   PORT=8080
   DB_HOST=localhost
   DB_PORT=5435
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=payment_db
   DB_SSLMODE=disable
   ```

   `.env` is ignored by Git. Keep these settings in sync with the PostgreSQL
   service. In particular, `DB_PORT` is `5435` for connections from the host;
   the container listens on `5432` internally.

3. Apply the initial schema:

   ```powershell
   Get-Content .\migrations\000001_create_transactions_table.up.sql |
     docker exec -i payment_postgres psql -U postgres -d payment_db
   ```

   The migration creates the `transaction_status` enum and `transactions`
   table. It is intended to be applied once to a new database. The matching
   down migration drops the table and enum, and therefore deletes transaction
   data.

4. Start the API:

   ```sh
   go run ./cmd/api
   ```

   The server listens on port `8080` by default. Set `PORT` to use another
   port. The application loads `.env` from its working directory and requires
   a reachable PostgreSQL database at startup.

## API

### Create a transaction

`POST /v1/transactions`

Example request:

```json
{
  "merchant_id": "11111111-1111-4111-8111-111111111111",
  "amount": 125000.00,
  "currency": "IDR",
  "payment_method": "bank_transfer",
  "reference_id": "order-123"
}
```

`merchant_id` must be a valid UUID. `amount` must be greater than zero, and
`reference_id` must be non-empty and unique. Include `currency` and
`payment_method`; the API passes these values to the database when inserting a
transaction.

On success, the API returns HTTP `201 Created` with a JSON response containing
a message and the created transaction:

```json
{
  "message": "Transaction created successfully",
  "data": {
    "id": "generated-uuid",
    "merchant_id": "11111111-1111-4111-8111-111111111111",
    "amount": 125000,
    "currency": "IDR",
    "status": "pending",
    "payment_method": "bank_transfer",
    "reference_id": "order-123",
    "created_at": "2026-10-10T00:00:00Z",
    "updated_at": "2026-10-10T00:00:00Z"
  }
}
```

Malformed JSON and transaction creation errors are returned as HTTP `400 Bad
Request`.

## Database and migrations

The PostgreSQL service uses database `payment_db`, user `postgres`, and
password `postgres` by default for local development. Do not use these
development credentials in a production environment.

PostgreSQL persists data in the named `postgres_data` volume. Changing
`POSTGRES_PASSWORD` in Compose does not change the password of an already
initialized database volume; the database credentials and `DB_PASSWORD` must
match.

SQL migrations are in `migrations/`. Apply the `.up.sql` migration to create
the schema. Use the corresponding `.down.sql` only when intentionally removing
the schema and its data.

## Project layout

```text
cmd/api/                         API entry point and route setup
internal/database/               PostgreSQL connection setup
internal/handler/                HTTP request handling
internal/service/                Transaction validation and business logic
internal/repository/             PostgreSQL transaction persistence
migrations/                      Database up/down SQL scripts
docker-compose.yml               Local PostgreSQL service
```
