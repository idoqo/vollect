package worker

import (
	"database/sql"
	"fmt"
	"gitlab.com/idoko/vollect/db"
	"log"
	"time"
)

type Worker struct {
	db db.Database
	pollInterval time.Duration
	StopChan chan bool
}

func NewWorker(database db.Database) *Worker {
	w := &Worker{
		db: database,
		pollInterval: 5 * time.Second,
	}
	return w
}

func (w *Worker) Run() error {
	for {
		select {
		case <- w.StopChan:
			return nil
		default:
			if run, err := w.NextJob(); err != nil {
				log.Println(err.Error())
				//todo mark job as error so it skips the next time
			} else if !run {
				// no job was found, take a nap
				time.Sleep(w.pollInterval)
			}
		}
	}
}

func (w *Worker) NextJob() (run bool, err error) {
	log.Println("calling next job...")
	task := db.Task{}
	query := `SELECT id, name, payload, status FROM vollect_tasks
		WHERE status = $1 ORDER BY id 
		FOR UPDATE SKIP LOCKED
		LIMIT 1`
	row := w.db.Conn.QueryRow(query, "pending")
	err = row.Scan(&task.Id, &task.Name, &task.Payload, &task.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	// deconstruct the payload and see if we can execute the handler
	log.Println(string(task.Payload))
	//defer task.Done()
	handler, exists := db.JobHandlers[task.Name]
	if !exists {
		return false, fmt.Errorf("no handler defined for %s", task.Name)
	}
	err = handler.Handle()
	return true, err
}