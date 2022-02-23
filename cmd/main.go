package main

import (
	"fmt"

	"github.com/siraj18/avito-test/internal/handlers"
	"github.com/siraj18/avito-test/internal/server"
)

func main() {

	handler := handlers.NewHandler()
	server := server.NewServer(":8000", handler.InitRoutes())
	err := server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
