package balance

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	moqAccrual "github.com/oleshko-g/oggophermart/internal/gen/transport/moq/accrual"
	"github.com/oleshko-g/oggophermart/internal/oglog"
	"github.com/oleshko-g/oggophermart/internal/storage"
	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql"
	"github.com/oleshko-g/oggophermart/internal/transport"
)

func Test_processAccrual(t *testing.T) {
	cfg := configueStorage()

	var storageBalance storage.Balance

	storageBalance, err := sql.New(cfg)
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
			orderIDString:         "019be993-bf1a-7088-a5d4-006bd660cc73",
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

func configueStorage() *db.Config {
	var cfg db.Config

	cfg.DSN().Set("postgres://gennadyoleshko:@localhost:5432/oggophermart?sslmode=disable")

	return &cfg
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

// result=&{Order:388772667448878 Status:PROCESSED Accrual:0x14000388480}
