package db

import (
	"encoding/json"
	"net/http"
)

type TaskHandler interface {
	Handle() error
}

var JobHandlers = make(map[string]TaskHandler)

type Task struct {
	Id      int
	Name 	string
	Status  string `json:"status"`
	Handler TaskHandler
	Payload []byte
	db Database
}

func NewTask(name string, database Database, taskHandler TaskHandler) (*Task, error) {
	task := &Task{
		Name: name,
		Status: "pending",
		Handler: taskHandler,
		db: database,
	}
	p := make(map[string]interface{})
	p["handler"] = name
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	task.Payload = payload
	return task, nil
}

func (t *Task) Queue() error {
	if _, exists := JobHandlers[t.Name]; !exists {
		JobHandlers[t.Name] = t.Handler
	}

	var id int
	query := "INSERT INTO vollect_tasks (name, payload, status) VALUES ($1, $2, $3) RETURNING id"
	err := t.db.Conn.QueryRow(query, t.Name, t.Payload, t.Status).Scan(&id)
	if err != nil {
		return err
	}
	t.Id = id
	return nil
}

func (t *Task) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

