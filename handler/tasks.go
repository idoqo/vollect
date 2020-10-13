package handler

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/idoko/vollect/response"
	"net/http"
	"strconv"
)

var taskIdKey = "taskId"

func tasks(r chi.Router) {
	r.Get("/", getAllTasks)
	r.Post("/pause", pauseTask)
	r.Post("/resume", resumeTask)
	r.Post("/terminate", terminateTask)

	r.Route("/{taskId}", func (r chi.Router) {
		r.Use(TaskContext)
		r.Get("/", getTask)
	})
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {}

func getTask(w http.ResponseWriter, r *http.Request) {}

func pauseTask(w http.ResponseWriter, r *http.Request) {}

func resumeTask(w http.ResponseWriter, r *http.Request) {}

func terminateTask(w http.ResponseWriter, r *http.Request) {}

func TaskContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskId := chi.URLParam(r, "taskId")
		if taskId == "" {
			render.Render(w, r, response.ErrBadRequest(fmt.Errorf("task ID is required")))
			return
		}
		id, err := strconv.Atoi(taskId)
		if err != nil {
			render.Render(w, r, response.ErrBadRequest(fmt.Errorf("invalid task ID")))
			return
		}
		ctx := context.WithValue(r.Context(), taskIdKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
