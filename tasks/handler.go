package tasks

import (
  "fmt"
	"net/http"
	"database/sql"
  //"encoding/json"
)

type handler struct{}

func InitializeHandler(db *sql.DB) *handler {
  // initialize struct and statements
  return &handler{}
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("For the tasks\n")
  fmt.Fprintf(w, "Heeyyyaaa.")
}
