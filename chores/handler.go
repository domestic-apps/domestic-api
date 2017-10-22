package chores

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type handler struct {
	addChoreStmt        *sql.Stmt
	getAllChoresStmt    *sql.Stmt
	getRecentChoresStmt *sql.Stmt
	modifyChoreStmt     *sql.Stmt
	deleteChoreStmt     *sql.Stmt
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
		addChoreStmt:        initStmt(addChoreStmt, db),
		getAllChoresStmt:    initStmt(getAllChoresStmt, db),
		getRecentChoresStmt: initStmt(getRecentChoresStmt, db),
		modifyChoreStmt:     initStmt(modifyChoreStmt, db),
		deleteChoreStmt:     initStmt(deleteChoreStmt, db),
	}
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var chore Chore
	switch r.Method {
	case http.MethodPost:
		// TODO more input validation
		err := json.Unmarshal(body, &chore)
		if err != nil {
			log.Println("error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.addChore(&chore)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		choreBytes, _ := json.Marshal(chore)
		w.Write(choreBytes)
	case http.MethodGet:
		choreList, err := h.getAllChores()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		choreBytes, _ := json.Marshal(choreList)
		w.Header().Set("Content-Type", "application/json")
		w.Write(choreBytes)
	case http.MethodPut:
		// TODO more input validation
		err := json.Unmarshal(body, &chore)
		if err != nil {
			log.Println("error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.modifyChore(&chore)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK) // TODO modified?
	case http.MethodDelete:
		// TODO more input validation
		err := json.Unmarshal(body, &chore)
		if err != nil {
			log.Println("error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.deleteChore(chore.Key)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
