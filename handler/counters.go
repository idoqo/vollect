package handler

import (
	"github.com/go-chi/chi"
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

}
