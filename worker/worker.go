package worker

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/idoko/vollect/db"
	"log"
	"sync"
	"time"
)

type Worker struct {
	db db.Database
	pollInterval time.Duration
	RunningJobs map[int]db.Task
	StopChan chan bool

	pause chan int // used to signal to handlers to pause their running tasks
	terminate chan int //used to signal to handlers to stop their running tasks
}

func NewWorker(database db.Database) *Worker {
	w := &Worker{
		db: database,
		pollInterval: 5 * time.Second,
		RunningJobs: make(map[int]db.Task),
		pause: make(chan int),
		terminate: make(chan int),
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
			} else if !run {
				// no job was found, take a nap
			}
			time.Sleep(w.pollInterval)
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
		// mark as failed so it doesn't get executed again
		db.FailTask(w.db, task.Id)
		return false, fmt.Errorf("no handler defined for %s", task.Name)
	}
	w.RunningJobs[task.Id] = task
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = handler.Handle(w.pause, w.terminate)
	}()
	wg.Wait()
	return true, err
}

func (w *Worker) Pause(taskId int) error {
	t, exists := w.RunningJobs[taskId]
	if !exists {
		// if the job is not running, simply update the state so it doesn't get executed in the next run
		//return db.PauseTask(w.db, taskId)
		return errors.New(fmt.Sprintf("no running job found with ID: %d", taskId))
	}
	t.Handler, _ = db.JobHandlers[t.Name]
	data, err := t.Handler.OnPause()
	if err != nil {
		return err
	}
	w.pause <- 1
	if len(data) > 0 {
		payload := make(map[string]interface{})
		err = json.Unmarshal(t.Payload, &payload)
		if err != nil {
			return err
		}
		// update the payload to include the *current* state of the job
		for k, v := range data {
			payload[k] = v
		}
		t.Payload, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}
	err = t.Pause(w.db)
	return err
}

func (w *Worker) Stop(taskId int) error {
	_, exists := w.RunningJobs[taskId]
	if exists {
		w.terminate <- 1
	}
	// if the job is not running, simply update the state so it doesn't get executed in the next run
	return db.DeleteTask(w.db, taskId)
}

func (w *Worker) Resume(taskId int) error {
	var exists bool
	taskInfo, err := db.GetPausedTask(w.db, taskId)
	if err != nil {
		return err
	}
	taskInfo.Handler, exists = db.JobHandlers[taskInfo.Name]
	if !exists {
		return errors.New("handler not found. you may have to restart the task")
	}
	mp := make(map[string]interface{})
	json.Unmarshal(taskInfo.Payload, &mp)
	err = taskInfo.Handler.OnResume(mp)
	if err != nil {
		return err
	}
	return taskInfo.Resume(w.db)
}
