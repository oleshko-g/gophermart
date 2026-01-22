package balance

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/oleshko-g/oggophermart/internal/oglog"
	_ "github.com/oleshko-g/oggophermart/internal/storage"
	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql"
	moqAccrual "github.com/oleshko-g/oggophermart/internal/gen/transport/moq/accrual"

)

var (
	svc *balanceSvc
)

func TestMain(m *testing.M) {
	cfg := configueStorage()
	storage, err := sql.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	svc = New(oglog.NewLoggingCtx(), storage, nil, &moqAccrual.AccrualMock{})
}

func Test_processAccrual(t *testing.T) {
	orderID, _ := uuid.NewRandom()
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	err := svc.processAccrual(ctx, orderID)
	if err != nil {
		t.Error(err)
	}
}

func configueStorage() *db.Config {
	var cfg db.Config

	cfg.DSN().Set("postgres://gennadyoleshko:@localhost:5432/oggophermart?sslmode=disable")

	return &cfg
}
