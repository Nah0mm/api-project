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
	apiClient *APIClient
}

func NewServer(baseURL string) *Server {
	return &Server{
		jobs:      make(chan model.Job, 10),
		apiClient: NewAPIClient(baseURL),
	}
}

func (s *Server) Shutdown() {
	close(s.jobs)
	s.wg.Wait()
}

func (s *Server) PostMethodHanlder(w http.ResponseWriter, r *http.Request) {
	var job model.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	s.jobs <- job
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "job accepted",
		"id":      job.ID,
	})
}

func (s *Server) StartWorkers() {
	for i := 1; i <= 3; i++ {
		s.wg.Add(1)
		go func(id int) {
			defer s.wg.Done()
			for job := range s.jobs {
				fmt.Printf("Worker %d processing job %s\n", id, job.ID)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := s.apiClient.Process(ctx, job)
				cancel()
				if err != nil {
					fmt.Printf("Error calling external API on job %s", job.ID)
					continue
				}
				fmt.Printf("Worker %d finished job %s\n", id, job.ID)
			}
		}(i)
	}
}
