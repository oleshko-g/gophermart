package balance

import (
	"context"
	"embed"
	"fmt"
	"log"
	"path"
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
	t.Cleanup(func() {
		ctx := context.Background()
		storageBalance.TearDown(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	tests := []struct {
		name                  string
		orderNumber           string
		orderIDString         string
		fetchOrderAccrualFunc fetchOrderAccrualFunc
		transport.FetchOrderAccrualResult
	}{
		{
			name:                   "empty accrual",
			orderIDString: "019c22a8e94f7d75bcb973b05f070535",
			orderNumber:    "568082882086285",
			fetchOrderAccrualFunc: emptyFetchOrderAccrual,
		},
		{
			name: "accrual status INVALID accrual",
			orderIDString: "019c22a8ea2974ab87a39959d1301bf0",
			orderNumber: "73568464702",
			fetchOrderAccrualFunc: invalidFetchOrderAccrualFunc,
		},
		{
			name:  "accrual status REGISTERED",
			orderIDString: "019c38078bdf7ed4a187d33b193aa2d9",
			orderNumber: "5441240171117",
			fetchOrderAccrualFunc: registeredFetchOrderAccrual,
		},
		{
			name: "accrual status PROCESSING",
			orderIDString: "019c22a8e55774419eae11fbe458db69",
			orderNumber: "167862623743756",
			fetchOrderAccrualFunc: processingFetchOrderAccrual,
		},
		{
			name: "accrual status PROCESSED no accrual",
			orderIDString: "019c27b58f537dcb8aa9fd9ca0d61f93",
			orderNumber: "860548181160",
			fetchOrderAccrualFunc: processedNoAccrualFetchOrderAccrual,
		},
		{
			name: "accrual status PROCESSED with accrual",
			orderIDString: "019c380787e87a77b9f8d392ab5797df",
			orderNumber: "2402188500348",
			fetchOrderAccrualFunc: processedFetchOrderAccrual,
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

// newTestStorage creates a new [sql.Storage] and runs sql statements from ./testdata/{testName} or returns an error
func newTestStorage(testName string) (*sql.Storage, error) {
	tn := strings.ToLower(testName)
	testDSN := fmt.Sprintf(
		"postgres://gennadyoleshko:@localhost:5432/gophermart_%s?sslmode=disable", tn)

	cfg := &db.Config{}
	err := cfg.DSN().Set(testDSN)
	if err != nil {
		return nil, err
	}

	storageBalance, err := sql.New(cfg)
	if err != nil {
		return nil, err
	}

	sqlStmt, err := testDataFS.ReadFile(path.Join("testdata", testName+".sql"))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = storageBalance.Exec(ctx, string(sqlStmt))
	if err != nil {
		return nil, err
	}

	return storageBalance, nil
}

type fetchOrderAccrualFunc func(context.Context, transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error)

func emptyFetchOrderAccrual(ctx context.Context, payload transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error) {
	return nil, nil
}

func invalidFetchOrderAccrualFunc(ctx context.Context, payload transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error) {
	return &transport.FetchOrderAccrualResult{
		Order:   payload.Number,
		Status:  transport.OrderAccrualStatusInvalid,
		Accrual: nil,
	}, nil
}

func registeredFetchOrderAccrual(ctx context.Context, payload transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error) {
	return &transport.FetchOrderAccrualResult{
		Order:   payload.Number,
		Status:  transport.OrderAccrualStatusRegistered,
		Accrual: nil,
	}, nil
}

func processingFetchOrderAccrual(ctx context.Context, payload transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error) {
	return &transport.FetchOrderAccrualResult{
		Order:   payload.Number,
		Status:  transport.OrderAccrualStatusProcessing,
		Accrual: nil,
	}, nil
}

func processedNoAccrualFetchOrderAccrual(ctx context.Context, payload transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error) {
	return &transport.FetchOrderAccrualResult{
		Order:   payload.Number,
		Status:  transport.OrderAccrualStatusProcessed,
		Accrual: nil,
	}, nil
}

func processedFetchOrderAccrual(ctx context.Context, payload transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error) {
	accrual := float64(729.98)
	return &transport.FetchOrderAccrualResult{
		Order:   payload.Number,
		Status:  transport.OrderAccrualStatusProcessed,
		Accrual: &accrual,
	}, nil
}
