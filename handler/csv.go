package handler

import (
	"encoding/csv"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/idoko/vollect/db"
	"gitlab.com/idoko/vollect/response"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type csvParser struct {
	filename string
	currentRow int
}

func (c *csvParser) Handle(pause, terminate chan int) error {
	file, err := os.Open(c.filename)
	if err != nil {
		return err
	}
	r := csv.NewReader(file)
	for {
		select {
		case <-pause:
			log.Println("pausing csv parser :D")
			return nil
		case <-terminate:
			log.Println("terminating csv parser :D")
			return nil
		default:
			//the built-in csv doesn't support jumping to a specific row,
			// and I can't think of a good way to skip the first N rows (where N is the
			// c.currentRow value that was saved)
			row, err := r.Read()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			c.currentRow++
			log.Println(row[0])
			time.Sleep(3 * time.Second)
		}
	}
}

func (c *csvParser) OnPause() (state map[string]interface{}, err error) {
	state = make(map[string]interface{})
	state["current_row"] = c.currentRow
	return state, nil
}

func (c *csvParser) OnResume(state map[string]interface{}) error {
	currentRow, exists := state["current_row"]
	log.Println(currentRow)
	if exists {
		c.currentRow = int(currentRow.(float64))
	}
	return nil
}

func csvRoute(r chi.Router) {
	r.Post("/upload", upload)
}

func upload(w http.ResponseWriter, r *http.Request) {
	// file should be <= 10MB
	r.ParseMultipartForm(10  * 1024 * 1024)

	file, _, err := r.FormFile("csv")
	if err != nil {
		render.Render(w, r, response.ErrBadRequest(err))
	}
	defer file.Close()
	rand.Seed(time.Now().UnixNano())
	fileName := randSeq(10) + ".csv"
	dst, err := os.Create(fileName)
	defer dst.Close()
	if err != nil {
		render.Render(w, r, response.ErrServerError(err))
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		render.Render(w, r, response.ErrServerError(err))
		return
	}
	task, err := createImportTask(fileName)

	if err != nil {
		render.Render(w, r, response.ErrServerError(err))
	}

	render.Render(w, r, &defaultMessage{Message: fmt.Sprintf("file upload complete. Task ID is %d", task.Id)})
}

func createImportTask(csvFile string) (*db.Task, error) {
	rand.Seed(time.Now().UnixNano())
	taskName := randSeq(20)
	parser := &csvParser{
		filename: csvFile,
		currentRow: 0,
	}
	task, err := db.NewTask(taskName, parser)
	if err != nil {
		return task, err
	}
	err = task.Queue(dbInstance)
	return task, err
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}