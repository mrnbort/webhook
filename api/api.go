package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/mrnbort/webhook.git/task"
	"log"
	"net/http"
	"time"
)

//go:generate moq -out executor_mock.go . Executor

// Service allows access to the config file
type Service struct {
	Executor   Executor
	Port       string
	httpServer *http.Server
}

// Executor interface provides access to the function that executes commands from the config file for the requested task
type Executor interface {
	ExecuteTask(ctx context.Context, request task.Task) error
}

// JSON is a map alias, just for convenience
type JSON map[string]interface{}

// Run the listener and request's router, activates the rest server
func (s Service) Run(ctx context.Context) error {
	s.httpServer = &http.Server{
		Addr:         s.Port,
		Handler:      s.routes(),
		ReadTimeout:  time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		log.Printf("[DEBUG] termination requested")
		err := s.httpServer.Close()
		if err != nil {
			log.Printf("[WARN] can't close server: %v", err)
		}
		log.Printf("[INFO] server closed")
	}()

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("service failed to run, err:%v", err)
	}
	return nil
}

func (s Service) routes() chi.Router {
	mux := chi.NewRouter()
	mux.Post("/execute", s.postExecuteTask)

	return mux
}

// POST /execute?id={task_id}&key={key}
func (s Service) postExecuteTask(w http.ResponseWriter, r *http.Request) {
	request := task.Task{ID: r.URL.Query().Get("id"), Auth: r.URL.Query().Get("key")}
	ctx := r.Context()

	if err := s.Executor.ExecuteTask(ctx, request); err != nil {
		log.Printf("[WARN] can't execute %v: %v", request.ID, err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, JSON{"error": err.Error()})
		return
	}
	render.JSON(w, r, JSON{"status": "ok"})
}
