package main

import (
	"time"

	"github.com/siraj18/avito-test/internal/db/postgresdb"
	"github.com/siraj18/avito-test/internal/handlers"
	"github.com/siraj18/avito-test/internal/server"
	"github.com/siraj18/avito-test/pkg/postgres"
	"github.com/sirupsen/logrus"
)

func main() {
	// docker run --name some-postgres -p 5433:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres:14.2-alpine
	// 05802b752dcc15626f922580104bad3a - TOKEN
	conStr := "postgresql://postgres:mysecretpassword@localhost:5433/postgres?sslmode=disable"
	db, err := postgres.NewDb(conStr)
	if err != nil {
		logrus.Fatal(err)
	}

	rep, err := postgresdb.NewSqlRepository(db)
	if err != nil {
		logrus.Fatal(err)
	}

	handler := handlers.NewHandler(rep)

	server := server.NewServer(":8000", handler.InitRoutes(), time.Second*10)
	if err := server.Run(); err != nil {
		logrus.Fatal(err)
	}

}
