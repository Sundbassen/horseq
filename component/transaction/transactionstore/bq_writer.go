package transactionstore

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/gofrs/uuid/v5"
	"github.com/sundbassen/horseq/component/transaction"
)

type bqWriter struct {
	client *bigquery.Client
}

type bqTransaction struct { //nolint:go-staticcheck
	ID        uuid.UUID
	Timestamp string
	ProjectID string
	ValueUSD  float64
}

func (t *bqTransaction) Save() (map[string]bigquery.Value, string, error) { //nolint:go-staticcheck
	return map[string]bigquery.Value{
		"transactions_id": t.ID.String(), // Convert uuid.UUID to STRING
		"project_id":      t.ProjectID,
		"timestamp":       t.Timestamp,
		"value_usd":       t.ValueUSD,
		"created_at":      time.Now(),
		"updated_at":      time.Now(),
	}, "", nil
}

func NewBQWriter(client *bigquery.Client) *bqWriter {
	return &bqWriter{
		client: client,
	}
}

func toBQTransaction(t *transaction.Transaction) *bqTransaction {
	return &bqTransaction{
		ID:        t.ID,
		Timestamp: t.Timestamp.Format(time.RFC3339),
		ProjectID: t.ProjectID,
		ValueUSD:  t.ValueUSD,
	}
}

func (b *bqWriter) Create(ctx context.Context, transactions []*transaction.Transaction) error {
	bqTransactions := make([]*bqTransaction, len(transactions))

	for i, t := range transactions {
		bqTransactions[i] = toBQTransaction(t)
	}

	inserter := b.client.Dataset("horseq").Table("transactions").Inserter()
	return inserter.Put(ctx, bqTransactions)
}
