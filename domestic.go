package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/domestic-apps/domestic-api/chores"
	"github.com/domestic-apps/domestic-api/tasks"
	"time"
)

type secrets struct {
	Uname string `json:"username"`
	Pwd   string `json:"password"`
}

func main() {
	// get mysql username and password from configuration
	file, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		log.Fatal("File error: %v\n", err)
	}

	var s secrets
	json.Unmarshal(file, &s)

	// Set up Database
	db, err := sql.Open("mysql",
		s.Uname+":"+s.Pwd+"@tcp(localhost:3306)/domestic?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare statements in application handlers.
	choresHandler := chores.InitializeHandler(db)
	tasksHandler := tasks.InitializeHandler(db)
	http.HandleFunc("/chores/", choresHandler.Handle)
	http.HandleFunc("/tasks/", tasksHandler.Handle)
	// log.Fatal(http.ListenAndServe(":8080", nil))

	// Let's try doing a database changey thing!
	currentTime := time.Now() // We'll query the chores for things at this time.
	//
	stmt, err := db.Prepare("INSERT INTO chores(chore_id, c_time) SELECT (chore_id, NULL) from chores where ? = true AND ((dwm = 'd') OR (dwm = 'w' AND day = ?) OR (dwm = 'm' AND date = ?))")

	if err != nil {
	log.Fatal(err)
	}
	stmt.Exec("morning", currentTime.Weekday, currentTime.Day)
}
