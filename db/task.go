package db

import (
	"encoding/json"
	"errors"
	"net/http"
)

type TaskHandler interface {
	Handle(pause chan int) error
	OnPause() (state map[string]interface{}, err error)
	OnResume(state map[string]interface{}) error
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
}

func NewTask(name string, taskHandler TaskHandler) (*Task, error) {
	task := &Task{
		Name: name,
		Status: PendingStatus,
		Handler: taskHandler,
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

func (t *Task) Pause(db Database) error {
	var err error
	if db.Conn == nil {
		err = errors.New("no available connection to database")
	}
	query := `UPDATE vollect_tasks
		SET payload = $1, status = $2 WHERE id = $3 
		RETURNING id, name, payload, status`
	_, err = db.Conn.Exec(query, t.Payload, "paused", t.Id)

	return err
}

func PauseTask(db Database, taskId int) error {
	query := `UPDATE vollect_tasks
		SET status = $2 WHERE id = $3 
		RETURNING id, name, payload, status`
	_, err := db.Conn.Exec(query, "paused", taskId)
	return err
}

func FailTask(db Database, taskId int) {
	query := `UPDATE vollect_tasks SET status = $1 WHERE id = $2`
	_, _ = db.Conn.Exec(query, "failed", taskId)
	return
}

func DeleteTask(db Database, taskId int) error {
	query := `DELETE FROM vollect_tasks WHERE id = $1`
	_, err := db.Conn.Exec(query, taskId)
	return err
}

func (t *Task) Resume(db Database) error {
	query := `UPDATE vollect_tasks SET status = $1 WHERE id = $2`
	_, err := db.Conn.Exec(query, "pending", t.Id)
	return err
}

func GetPausedTask(database Database, taskId int) (*Task, error) {
	task := &Task{}
	query := `SELECT id, name, payload, status FROM vollect_tasks
				WHERE status = $1 AND id = $2`
	row := database.Conn.QueryRow(query, "paused", taskId)
	err := row.Scan(&task.Id, &task.Name, &task.Payload, &task.Status)
	return task, err
}

func (t *Task) Queue(db Database) error {
	if _, exists := JobHandlers[t.Name]; !exists {
		JobHandlers[t.Name] = t.Handler
	}

	var id int
	query := "INSERT INTO vollect_tasks (name, payload, status) VALUES ($1, $2, $3) RETURNING id"
	err := db.Conn.QueryRow(query, t.Name, t.Payload, t.Status).Scan(&id)
	if err != nil {
		return err
	}
	t.Id = id
	return nil
}

func (t *Task) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

