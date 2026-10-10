package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Transaction represents a payment transaction stored in the database.
type Transaction struct {
	ID            string    `json:"id"`             // Unique transaction identifier.
	MerchantID    string    `json:"merchant_id"`    // Merchant that owns the transaction.
	Amount        float64   `json:"amount"`         // Transaction amount.
	Currency      string    `json:"currency"`       // ISO currency code for the amount.
	Status        string    `json:"status"`         // Current transaction status.
	PaymentMethod string    `json:"payment_method"` // Payment method used for the transaction.
	ReferenceID   string    `json:"reference_id"`   // Merchant-provided transaction reference.
	CreatedAt     time.Time `json:"created_at"`     // Time the transaction was created.
	UpdatedAt     time.Time `json:"updated_at"`     // Time the transaction was last updated.
}

// TransactionRepository defines persistence operations for transactions.
type TransactionRepository interface {
	// Create inserts trx within tx. A nil tx uses the repository's database
	// connection directly.
	Create(ctx context.Context, tx *sql.Tx, trx *Transaction) error
}

// PostgresTransactionRepository persists transactions in PostgreSQL.
type PostgresTransactionRepository struct {
	// DB is the PostgreSQL connection used for queries outside a transaction.
	DB *sql.DB
}

// NewPostgresTransactionRepository creates a PostgreSQL-backed transaction
// repository.
func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{DB: db}
}

// Create inserts trx and populates its ID and timestamps.
func (r *PostgresTransactionRepository) Create(ctx context.Context, tx *sql.Tx, trx *Transaction) error {
	trx.ID = uuid.New().String()
	query := `
		INSERT INTO transactions (id, merchant_id, amount, currency, status, payment_method, reference_id)
		VALUES ($1, $2, $3, $4, 'pending', $5, $6)
		RETURNING created_at, updated_at
	`

	var executor interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	} = r.DB

	if tx != nil {
		executor = tx
	}

	return executor.QueryRowContext(
		ctx, query,
		trx.ID, trx.MerchantID, trx.Amount, trx.Currency, trx.PaymentMethod, trx.ReferenceID,
	).Scan(&trx.CreatedAt, &trx.UpdatedAt)
}
