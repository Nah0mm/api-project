package api

import (
	"api-worker/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	jobs      chan model.Job
	wg        sync.WaitGroup
	apiClient ApiClient
}

func NewServer(baseURL string) *Server {
	return &Server{
		jobs:      make(chan model.Job, 10),
		apiClient: *NewApiClient(baseURL),
	}
}

func (server *Server) Shutdown() {
	close(server.jobs)
	server.wg.Wait()

}

func (server *Server) PostMethodHandler(w http.ResponseWriter, r *http.Request) {
	var job model.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
	}
	server.jobs <- job
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "message accepted",
		"job id":  job.ID,
	})
}

func (server *Server) StartWorkers() {
	for i := 1; i <= 3; i++ {
		server.wg.Add(1)
		go func(id int) {
			defer server.wg.Done()
			for job := range server.jobs {
				var err error
				fmt.Printf("Worker %d processing job %s\n", id, job.ID)
				for attempt := 1; attempt <= 3; attempt++ {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					err = server.apiClient.Process(ctx, job)
					cancel()
					if err == nil {
						fmt.Printf("Worker %d finished job %s\n", id, job.ID)
						break
					}
					fmt.Printf(
						"Worker %d attempt %d failed for job %s: %v\n",
						id,
						attempt,
						job.ID,
						err,
					)
					if attempt < 3 {
						backoff := time.Duration(attempt) * time.Second

						fmt.Printf(
							"Worker %d retrying job %s in %v\n",
							id,
							job.ID,
							backoff,
						)

						select {
						case <-time.After(backoff):
							break
						case <-ctx.Done():
							return
						}
					}
				}
				if err != nil {
					fmt.Printf(
						"Worker %d permanently failed job %s\n",
						id,
						job.ID,
					)
				}
			}
		}(i)
	}
}
