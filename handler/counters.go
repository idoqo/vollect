package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/idoko/vollect/db"
	"gitlab.com/idoko/vollect/response"
	"log"
	"net/http"
	"time"
)

type counter struct {
	current int
	step int
}

func (c *counter) Handle(pause chan int) error {
	for c.current <= 1500 {
		select {
		case <-pause:
			log.Println("pausing :D")
			return nil
		default:
			c.current += c.step
			log.Println(c.current)
			time.Sleep(3 * time.Second)
		}
	}
	return nil
}

func (c *counter) OnPause() (state map[string]interface{}, err error) {
	state = make(map[string]interface{})
	state["current"] = c.current
	return state, nil
}

func counters(r chi.Router) {
	r.Get("/", createCounter)
}

func createCounter(w http.ResponseWriter, r *http.Request) {
	c := &counter{current: 1, step: 1}
	tsk, err := db.NewTask("counter", dbInstance, c)
	if err != nil {
		render.Render(w, r, &response.ErrorResponse{Err: err, StatusCode: 500, Message: err.Error()})
		return
	}
	if err := tsk.Queue(); err != nil {
		render.Render(w, r, &response.ErrorResponse{Err: err, StatusCode: 500, Message: err.Error()})
	}
	return
}

