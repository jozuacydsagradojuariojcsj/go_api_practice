package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/db"
)

var (
	ErrInsufficientFunds = errors.New("Insufficient Funds")
	ErrInvalidTransfer   = errors.New("Invalid Transfer")
	ErrUserNotFound      = errors.New("user not found")
)

type Store struct {
	pool *pgxpool.Pool
	*db.Queries
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool, Queries: db.New(pool)}
}

func (s *Store) Deposit(ctx context.Context, userID uuid.UUID, amount int64) error {
	if amount <= 0 {
		return ErrInvalidTransfer
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)
	q := s.Queries.WithTx(tx)

	if _, err := q.AddBalance(ctx, db.AddBalanceParams{ID: userID, Amount: amount}); err != nil {
		return err
	}

	if _, err := q.RecordTransaction(ctx, db.RecordTransactionParams{ToUserID: uuid.NullUUID{UUID: userID, Valid: true}, Amount: amount}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) Transfer(ctx context.Context, from, to uuid.UUID, amount int64) error {
	if amount <= 0 || from == to {
		return ErrInvalidTransfer
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	q := s.Queries.WithTx(tx)

	locked, err := q.LockUsers(ctx, []uuid.UUID{from, to})
	if err != nil {
		return err
	}
	if len(locked) != 2 {
		return errors.New("User not found")
	}

	if _, err := q.AddBalance(ctx, db.AddBalanceParams{ID: from, Amount: -amount}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return ErrInsufficientFunds
		}
		return err
	}
	if _, err := q.AddBalance(ctx, db.AddBalanceParams{ID: to, Amount: amount}); err != nil {
		return err
	}
	if _, err := q.RecordTransaction(ctx, db.RecordTransactionParams{
		FromUserID: uuid.NullUUID{UUID: from, Valid: true},
		ToUserID:   uuid.NullUUID{UUID: to, Valid: true},
		Amount:     amount,
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
