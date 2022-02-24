package postgres

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

func NewDb(conStr string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", conStr)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
	}

	return db, nil
}
