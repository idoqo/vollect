package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/idoko/vollect/db"
	"gitlab.com/idoko/vollect/response"
	"log"
	"net/http"
)

type counter struct {
	current int
	step int
}

func (c *counter) Handle() error {
	c.current += c.step
	log.Println(c.current)
	return nil
}

func counters(r chi.Router) {
	r.Get("/", getCounter)
}

func getCounter(w http.ResponseWriter, r *http.Request) {
	c := &counter{current: 1, step: 1}
	tsk, err := db.NewTask("hello_world", dbInstance, c)
	if err != nil {
		render.Render(w, r, &response.ErrorResponse{Err: err, StatusCode: 500, Message: err.Error()})
		return
	}
	if err := tsk.Queue(); err != nil {
		render.Render(w, r, &response.ErrorResponse{Err: err, StatusCode: 500, Message: err.Error()})
	}
	return
}
