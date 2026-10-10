package service

import (
	"context"
	"database/sql"
	"errors"

	"go-payment-gateway/internal/repository"
)

// CreateTransactionDTO contains the client-provided fields used to create a
// transaction.
type CreateTransactionDTO struct {
	MerchantID    string  `json:"merchant_id"`    // Merchant that owns the transaction.
	Amount        float64 `json:"amount"`         // Transaction amount.
	Currency      string  `json:"currency"`       // ISO currency code for the amount.
	PaymentMethod string  `json:"payment_method"` // Payment method requested by the client.
	ReferenceID   string  `json:"reference_id"`   // Merchant-provided transaction reference.
}

// TransactionService defines the business operations available for
// transactions.
type TransactionService interface {
	// CreateTransaction validates the request and persists a pending
	// transaction.
	CreateTransaction(ctx context.Context, dto CreateTransactionDTO) (*repository.Transaction, error)
}

// paymentService implements transaction business logic using a database and
// transaction repository.
type paymentService struct {
	db   *sql.DB
	repo repository.TransactionRepository
}

// NewTransactionService constructs a TransactionService backed by db and repo.
func NewTransactionService(db *sql.DB, repo repository.TransactionRepository) TransactionService {
	return &paymentService{db: db, repo: repo}
}

// CreateTransaction validates dto and stores the new transaction atomically.
func (s *paymentService) CreateTransaction(ctx context.Context, dto CreateTransactionDTO) (*repository.Transaction, error) {
	if dto.Amount <= 0 {
		return nil, errors.New("nominal transaksi harus lebih besar dari 0")
	}
	if dto.ReferenceID == "" {
		return nil, errors.New("reference_id dari merchant wajib diisi")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	trx := &repository.Transaction{
		MerchantID:    dto.MerchantID,
		Amount:        dto.Amount,
		Currency:      dto.Currency,
		Status:        "pending",
		PaymentMethod: dto.PaymentMethod,
		ReferenceID:   dto.ReferenceID,
	}

	if err := s.repo.Create(ctx, tx, trx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return trx, nil
}
