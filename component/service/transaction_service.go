package service

import (
	"context"

	"github.com/sundbassen/horseq/component/transaction"
)

type TransactionService struct {
	readStore  transaction.ReadStore
	writeStore transaction.WriteStore
}

func NewTransactionService(rs transaction.ReadStore, ws transaction.WriteStore) *TransactionService {
	return &TransactionService{
		readStore:  rs,
		writeStore: ws,
	}
}

func (s *TransactionService) List(ctx context.Context) ([]*transaction.Transaction, error) {
	return s.readStore.List(ctx)
}

func (s *TransactionService) MapToNew(ctx context.Context) error {
	transactions, err := s.readStore.List(ctx)
	if err != nil {
		return err
	}

	return s.writeStore.Create(ctx, transactions)
}
