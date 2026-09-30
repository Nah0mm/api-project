package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Job struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

func main() {
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		var job Job

		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf(
			"Mock API received job: ID=%s Type=%s Amount=%d\n",
			job.ID,
			job.Type,
			job.Amount,
		)

		select {
		case <-time.After(4 * time.Second):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "processed",
				"id":     job.ID,
			})
			fmt.Printf("Job %s processed\n", job.ID)
		case <-r.Context().Done():
			fmt.Println("Request cancelled")
			return
		}
	})

	myProcessHttpServer := &http.Server{
		Addr: ":9090",
	}
	fmt.Println("Starting mock server on port 9090")
	if err := myProcessHttpServer.ListenAndServe(); err != nil {
		fmt.Printf("Starting mock server on port 9090 encountered error %v", err)
		return
	}
}
