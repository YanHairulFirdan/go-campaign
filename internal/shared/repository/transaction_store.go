package repository

import (
	"database/sql"
)

type TransactionStore struct {
	db *sql.DB
}

func NewTransactionStore(db *sql.DB) *TransactionStore {
	return &TransactionStore{db: db}
}

func (s *TransactionStore) ExecTx(fn func() error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}
