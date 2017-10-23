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

	"gopkg.in/robfig/cron.v2"
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
	timezone := "America/Los_Angeles"

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
	tasksHandler := tasks.InitializeHandler(db, timezone)
	http.HandleFunc("/chores/", choresHandler.Handle)
	http.HandleFunc("/tasks/", tasksHandler.Handle)
	c := cron.New()
	c.AddFunc("TZ="+timezone+" 0 0 7,19 * * *", tasksHandler.Cron) // 7am and 7pm every day
	c.Start()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
