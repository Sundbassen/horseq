package transaction

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Transaction struct {
	ID        uuid.UUID
	Timestamp time.Time
	ProjectID string
	ValueUSD  float64
}

type ReadStore interface {
	List(ctx context.Context) ([]*Transaction, error)
	// Delete(ctx context.Context, ID uuid.UUID) error
}

type WriteStore interface {
	Create(ctx context.Context, transactions []*Transaction) error

	// Update(ctx context.Context, t *Transaction) error
}
