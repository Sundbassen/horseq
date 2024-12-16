package transactionstore

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/sundbassen/horseq/component/transaction"

	"cloud.google.com/go/storage"
)

const (
	timeFormat   = "2006-01-02 15:04:05"
	cgTimeFormat = "02-01-2006"
	coinGeckoURL = "https://api.coingecko.com/api/v3/"
	tsIdx        = 1
	projectIDIdx = 2
	propsIdx     = 14
	numsIdx      = 15
)

type Props struct {
	CurrencySymbol string `json:"currencySymbol"`
}

// Struct to represent each element in the nums array
type Num struct {
	CurrencyValueDecimal string `json:"currencyValueDecimal"`
}

type transactionNotConverted struct {
	Timestamp      time.Time
	ProjectID      string
	CurrencySymbol string
	CurrencyValue  float64
}

type bucket struct {
	b                *storage.BucketHandle
	transactionsPath string
}

var _ transaction.ReadStore = (*bucket)(nil)

func NewBucket(b *storage.BucketHandle, transactionsPath string) *bucket {

	return &bucket{
		b:                b,
		transactionsPath: transactionsPath,
	}
}
func (b *bucket) ReadCSV(ctx context.Context, path string) (io.ReadCloser, error) {
	obj := b.b.Object(path)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (b *bucket) List(ctx context.Context) ([]*transaction.Transaction, error) {
	reader, err := b.ReadCSV(ctx, b.transactionsPath)
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	csvReader := csv.NewReader(reader)
	transactions := []*transactionNotConverted{}

	_, err = csvReader.Read() // Skip the header
	if err != nil {
		return nil, err
	}

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
			return nil, err
		}

		t, err := time.Parse(timeFormat, record[tsIdx])
		if err != nil {
			return nil, err
		}

		projectID := record[projectIDIdx]

		var props Props
		err = json.Unmarshal([]byte(record[propsIdx]), &props)
		if err != nil {
			return nil, err
		}

		var nums Num
		err = json.Unmarshal([]byte(record[numsIdx]), &nums)
		if err != nil {
			return nil, err
		}

		f, err := strconv.ParseFloat(nums.CurrencyValueDecimal, 64)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &transactionNotConverted{
			Timestamp:      t,
			ProjectID:      projectID,
			CurrencySymbol: props.CurrencySymbol,
			CurrencyValue:  f,
		})
	}

	convertedTransactions, err := convertToTransaction(ctx, transactions)
	if err != nil {
		return nil, err
	}

	return convertedTransactions, nil
}

type cgResponse struct {
	MarketData struct {
		CurrentPrice map[string]float64 `json:"current_price"`
	} `json:"market_data"`
}

func convertToTransaction(ctx context.Context, transactionsNotConv []*transactionNotConverted) ([]*transaction.Transaction, error) {
	transactions := make([]*transaction.Transaction, 0, len(transactionsNotConv))
	fetchedConversionRates := make(map[string]map[string]float64)

	for _, t := range transactionsNotConv {
		// Fetch the conversion rate for the currency symbol
		var convRate float64
		cgDate := t.Timestamp.Format(cgTimeFormat)
		if _, ok := fetchedConversionRates[cgTimeFormat]; !ok {

			url := fmt.Sprintf(coinGeckoURL+"coins/usd/history?date=%s", cgDate)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return nil, err
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, err
			}

			defer res.Body.Close()

			var response cgResponse

			if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
				return nil, err
			}

			fetchedConversionRates[cgTimeFormat] = response.MarketData.CurrentPrice
		}

		convRate, ok := fetchedConversionRates[cgTimeFormat][t.CurrencySymbol]

		// TODO: Since some of the currencies does not exist in coin gecko, random conversion rate is used.
		// See e.g. https://api.coingecko.com/api/v3/coins/slf/history?date=12-12-2024
		if !ok {
			// slog.Error("Conversion rate not found", slog.Any("currencySymbol", t.CurrencySymbol), slog.Any("date", cgDate))
			convRate = 1
		}

		u, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction.Transaction{
			ID:        u,
			Timestamp: t.Timestamp,
			ProjectID: t.ProjectID,
			ValueUSD:  t.CurrencyValue * convRate,
		})

	}
	return transactions, nil
}
