package balance

import (
	"context"
	"embed"
	"fmt"
	"log"
	"strings"

	"testing"

	"github.com/google/uuid"
	moqAccrual "github.com/oleshko-g/oggophermart/internal/gen/transport/moq/accrual"
	"github.com/oleshko-g/oggophermart/internal/oglog"
	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql"
	"github.com/oleshko-g/oggophermart/internal/transport"
	_ "github.com/pressly/goose/v3"
)

//go:embed testdata/*.sql
var testDataFS embed.FS

func Test_processAccrual(t *testing.T) {
	storageBalance, err := newTestStorage(t.Name())
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name                  string
		orderNumber           string
		orderIDString         string
		fetchOrderAccrualFunc fetchOrderAccrualFunc
	}{
		{
			name:                  "order processed with accrual",
			orderNumber:           "388772667448878",
			orderIDString:         "019bfe24-c85e-7c58-bca4-a9dfd7b95a5f",
			fetchOrderAccrualFunc: processedFetchOrderAccrualFunc,
		},
		{
			name:                  "empty order accrual",
			orderNumber:           "757483714",
			orderIDString:         "019be994-0e1f-7d0f-b288-67c60cc3218c",
			fetchOrderAccrualFunc: emptyFetchOrderAccrualFunc,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := New(oglog.NewLoggingCtx(), storageBalance, nil, &moqAccrual.AccrualMock{FetchOrderAccrualFunc: test.fetchOrderAccrualFunc})
			orderID, err := uuid.Parse(test.orderIDString)
			if err != nil {
				t.Error(err)
			}

			svc.processAccrual(context.Background(), orderID)
			if err != nil {
				t.Error(err)
			}
		})
	}

}

// newTestStorage creates a new [sql.Storage] and runs the ./testdata/{testName} migrations on it
func newTestStorage(testName string) (*sql.Storage, error) {
	testName = strings.ToLower(testName)
	testDSN := fmt.Sprintf(
		"postgres://gennadyoleshko:@localhost:5432/gophermart_%s?sslmode=disable", testName)

	cfg := &db.Config{}
	err := cfg.DSN().Set(testDSN)
	if err != nil {
		return nil, err
	}

	storageBalance, err := sql.New(cfg)
	if err != nil {
		return nil, err
	}
	sql, err := testDataFS.ReadFile(testName + ".sql")
	_ = err
	ctx := context.Background()
	err = storageBalance.Exec(ctx, string(sql))
	if err != nil {
		return nil, err
	}

	return storageBalance, nil
}

type fetchOrderAccrualFunc func(context.Context, transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error)

func emptyFetchOrderAccrualFunc(ctx context.Context, payload transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error) {
	return nil, nil
}

func processedFetchOrderAccrualFunc(ctx context.Context, payload transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error) {
	accrual := float64(729.98)
	return &transport.FetchOrderAccrualResult{
		Order:   payload.Number,
		Status:  transport.OrderAccrualStatusProcessed,
		Accrual: &accrual,
	}, nil
}
