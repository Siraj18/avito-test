package main

import (
	"fmt"
	"net/http"

	"github.com/siraj18/avito-test/internal/handlers"
)

func main() {

	server := &http.Server{
		Addr:    ":8000",
		Handler: handlers.NewHandler().InitRoutes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
