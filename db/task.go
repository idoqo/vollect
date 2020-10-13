package db

import (
	"encoding/json"
	"net/http"
)

type TaskHandler interface {
	Handle() error
}

type Task struct {
	id      int
	status  string `json:"status"`
	Handler TaskHandler

	payload []byte
	db Database
}

func NewTask(name string, database Database, taskHandler TaskHandler) (*Task, error) {
	task := &Task{
		status: "pending",
		Handler: taskHandler,
		db: database,
	}
	p := make(map[string]interface{})
	p["handler"] = name
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	task.payload = payload
	return task, nil
}

func (t *Task) Queue() error {
	var id int
	query := "INSERT INTO tasks (payload, status) VALUES ($1, $2) RETURNING id"
	err := t.db.Conn.QueryRow(query, t.payload, t.status).Scan(&id)
	if err != nil {
		return err
	}
	t.id = id
	return nil
}

func (t *Task) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

