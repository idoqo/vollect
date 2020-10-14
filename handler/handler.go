package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/idoko/vollect/db"
	"gitlab.com/idoko/vollect/response"
	"gitlab.com/idoko/vollect/worker"
	"net/http"
)

type defaultMessage struct {
	Message string `json:"message"`
}

var (
	dbInstance db.Database
	queueWorker *worker.Worker
)

func (dm *defaultMessage) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (dm *defaultMessage) Bind(r *http.Request) error {
	return nil
}

func NewHandler(database db.Database, wk *worker.Worker) http.Handler {
	dbInstance = database
	queueWorker = wk
	r := chi.NewRouter()
	r.MethodNotAllowed(notAllowedHandler)
	r.NotFound(notFoundHandler)
	r.Route("/tasks", tasks)
	r.Route("/counter", counters)
	r.Route("/csv", csvRoute)
	r.Get("/", index)
	return r
}

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, response.ErrMethodNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(404)
	render.Render(w, r, response.ErrNotFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	render.Render(w, r, &defaultMessage{Message: "hi there"})
}

