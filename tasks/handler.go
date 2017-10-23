package tasks

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type handler struct {
	setChoresNowStmt *sql.Stmt
	getAllTasksStmt  *sql.Stmt
	modifyTaskStmt   *sql.Stmt
	timezone         *time.Location
}

func initStmt(f func(db *sql.DB) (*sql.Stmt, error), db *sql.DB) *sql.Stmt {
	stmt, err := f(db)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

func InitializeHandler(db *sql.DB, timezone string) *handler {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatal(err)
	}
	return &handler{
		setChoresNowStmt: initStmt(setChoresNowStmt, db),
		getAllTasksStmt:  initStmt(getAllTasksStmt, db),
		modifyTaskStmt:   initStmt(modifyTaskStmt, db),
		timezone:         l,
	}
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		taskList, err := h.getAllTasks()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		taskBytes, _ := json.Marshal(taskList)
		w.Header().Set("Content-Type", "application/json")
		w.Write(taskBytes)
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO more input validation
		var task Task
		err = json.Unmarshal(body, &task)
		if err != nil {
			log.Println("error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.modifyTask(&task)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK) // TODO modified?
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Cron is run periodically, whenever the cron job in main triggers the function.
func (h *handler) Cron() {
	log.Println("Time: ", time.Now().In(h.timezone))
	err := h.setChoresNow(time.Now().In(h.timezone))
	if err != nil {
		log.Println(err)
	}
}
