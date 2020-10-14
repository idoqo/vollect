package handler

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/idoko/vollect/db"
	"gitlab.com/idoko/vollect/response"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type counter struct {
	current int
	step int
}

func (c *counter) Handle(pause, terminate chan int) error {
	for c.current <= 1500 {
		select {
		case <-pause:
			log.Println("pausing counter :D")
			return nil
		case <-terminate:
			log.Println("stopping counter :D")
			return nil
		default:
			log.Println(c.current)
			c.current += c.step
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

func (c *counter) OnResume(state map[string]interface{}) error {
	current, exists := state["current"]
	log.Println(current)
	if exists {
		c.current = int(current.(float64))
	}
	return nil
}

func counters(r chi.Router) {
	r.Get("/", createCounter)
}

func createCounter(w http.ResponseWriter, r *http.Request) {
	c := &counter{current: 1, step: 1}
	rand.Seed(time.Now().UnixNano())
	name := randSeq(20)
	tsk, err := db.NewTask(name, c)
	if err != nil {
		render.Render(w, r, &response.ErrorResponse{Err: err, StatusCode: 500, Message: err.Error()})
		return
	}
	if err := tsk.Queue(dbInstance); err != nil {
		render.Render(w, r, &response.ErrorResponse{Err: err, StatusCode: 500, Message: err.Error()})
	}
	render.Render(w, r, &defaultMessage{Message: fmt.Sprintf("counter created. Task ID is %d", tsk.Id)})
	return
}

