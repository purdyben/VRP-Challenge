package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	// server := http.Server{
	// 	Addr:    ":9091",
	// 	Handler: &customHandler{"Hello world from the custom handler"},
	// }
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Starting server")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
	fmt.Println("Exiting")
}
