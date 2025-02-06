package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"backend-engineering/pkg/httpx/middleware"
)

var jobs sync.Map

func main() {
	http.HandleFunc("/submit", middleware.Func(submit))
	http.HandleFunc("/checkstatus", middleware.Func(checkstatus))

	slog.Info("server started", "port", "5432")
	slog.Error(http.ListenAndServe(":5432", nil).Error())
}

func submit(w http.ResponseWriter, r *http.Request) {
	jobID := fmt.Sprintf("%d", time.Now().Unix())

	jobs.Store(jobID, 0)
	go updateJob(jobID, 0)

	fmt.Fprintf(w, "gobID: %s", jobID)
}

func checkstatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job_id")

	v, ok := jobs.Load(jobID)
	if !ok {
		fmt.Fprintf(w, "jobID: %s missing ", jobID)
	}

	for v.(int) < 100 {
		v, _ = jobs.Load(jobID)

		slog.Info("job in progress", "ID", jobID, "status", v)

		time.Sleep(1 * time.Second)
	}

	fmt.Fprintf(w, "progress for jobID: %s => %d", jobID, v)
}

func updateJob(jobID string, progress int) {
	for progress < 100 {
		progress += 10
		jobs.Store(jobID, progress)

		time.Sleep(3 * time.Second)
	}
}
