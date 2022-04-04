package main

import (
	"os"
	"time"

	"github.com/siraj18/avito-test/internal/currency"
	"github.com/siraj18/avito-test/internal/db/postgresdb"
	"github.com/siraj18/avito-test/internal/handlers"
	"github.com/siraj18/avito-test/internal/server"
	"github.com/siraj18/avito-test/pkg/postgres"
	"github.com/siraj18/avito-test/pkg/rediscache"
	"github.com/sirupsen/logrus"
)

func main() {
	conStr := os.Getenv("connection_string_postgres")
	redisConStr := os.Getenv("connection_string_redis")
	apiToken := os.Getenv("api_token")
	address := os.Getenv("address")

	db, err := postgres.NewDb(conStr, 10)
	if err != nil {
		logrus.Fatal(err)
	}

	rep, err := postgresdb.NewSqlRepository(db)
	if err != nil {
		logrus.Fatal(err)
	}

	rdb, err := rediscache.NewRedis(redisConStr, "", 0)
	if err != nil {
		logrus.Fatal(err)
	}

	cr := currency.NewCurrency(rdb, apiToken, 10*time.Minute)

	handler := handlers.NewHandler(rep, cr)

	server := server.NewServer(address, handler.InitRoutes(), time.Second*10)
	if err := server.Run(); err != nil {
		logrus.Fatal(err)
	}

}
