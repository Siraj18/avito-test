package main

import (
	"time"

	"github.com/siraj18/avito-test/internal/handlers"
	"github.com/siraj18/avito-test/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	handler := handlers.NewHandler()

	server := server.NewServer(":8000", handler.InitRoutes(), time.Second*10)
	if err := server.Run(); err != nil {
		logrus.Fatal(err)
	}

}
