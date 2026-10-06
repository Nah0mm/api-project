package main

import (
	"api-worker/api"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var server = api.NewServer("http://localhost:9090/process")

func main() {
	http.HandleFunc("/jobs", MethodHandlers)
	server.StartWorkers()
	myHttpServer := &http.Server{
		Addr: ":8080",
	}
	fmt.Println("Starting server on 8080")
	go func() {
		if err := myHttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server on 8080")
			return
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	fmt.Println("Shutdown server received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := myHttpServer.Shutdown(ctx); err != nil {
		fmt.Println("Error shutting down server")
	}
	server.Shutdown()
	fmt.Println("Server shutdown successfully")
}

func MethodHandlers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		server.PostMethodHandler(w, r)

	}
}
