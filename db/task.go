package db

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type TaskHandler interface {
	Handle(pause chan int) error
	OnPause() (state map[string]interface{}, err error)
}
const (
	PausedStatus = "paused"
	PendingStatus = "pending"
)
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
		Status: PendingStatus,
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

func (t *Task) Pause() error {
	var err error
	if t.db.Conn == nil {
		err = errors.New("no available connection to database")
	}
	log.Printf("%v\n", t)
	query := `UPDATE vollect_tasks
		SET payload = $1, status = $2 WHERE id = $3 
		RETURNING id, name, payload, status`
	_, err = t.db.Conn.Exec(query, t.Payload, "paused", t.Id)

	return err
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

func (t *Task) UseDB(database Database) {
	t.db = database
}

func (t *Task) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

