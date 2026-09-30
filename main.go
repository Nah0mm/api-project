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

func MethodHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		server.PostMethodHanlder(w, r)
	}
}

func main() {
	http.HandleFunc("/jobs", MethodHandler)
	server.StartWorkers()
	fmt.Println("Starting server on port 8080")
	myHttpServer := &http.Server{
		Addr: ":8080",
	}
	go func() {
		if err := myHttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error occurred starting server on port 8080\n%v\n", err)
			return
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	fmt.Println("Shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := myHttpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Error occurred shuttind down server\n%v\n", err)
	}
	server.Shutdown()
	fmt.Printf("Server shutdown gracefully\n")

}
