package tasks

import (
	"database/sql"
	//"encoding/json"
	//"io/ioutil"
	"log"
	"net/http"
)

type handler struct {
	setChoresNowStmt *sql.Stmt
	getAllTasksStmt  *sql.Stmt
	modifyTaskStmt   *sql.Stmt
}

func initStmt(f func(db *sql.DB) (*sql.Stmt, error), db *sql.DB) *sql.Stmt {
	stmt, err := f(db)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

func InitializeHandler(db *sql.DB) *handler {
	return &handler{
		setChoresNowStmt: initStmt(setChoresNowStmt, db),
		getAllTasksStmt:  initStmt(getAllTasksStmt, db),
		modifyTaskStmt:   initStmt(modifyTaskStmt, db),
	}
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("For the tasks\n")
	// fmt.Fprintf(w, "Heeyyyaaa.")
}
