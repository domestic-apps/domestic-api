package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/goji/httpauth"

	"github.com/domestic-apps/domestic-api/chores"
	"github.com/domestic-apps/domestic-api/tasks"

	"gopkg.in/robfig/cron.v2"
)

type secrets struct {
	Uname string `json:"username"`
	Pwd   string `json:"password"`
	Cert  string `json:"cert"`
	Key   string `json:"key"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "add-user" {
		addUser()
	} else {
		runServer()
	}
}

func runServer() {
	// get mysql username and password from configuration
	file, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		log.Fatalf("File error: %v\n", err)
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
	r := mux.NewRouter()
	r.HandleFunc("/chores", choresHandler.HandleCreate).Methods("POST")
	r.HandleFunc("/chores", choresHandler.HandleReadList).Methods("GET")
	r.HandleFunc("/chores", choresHandler.HandleUpdate).Methods("PUT")
	r.HandleFunc("/chores", choresHandler.HandleDelete).Methods("DELETE")
	r.HandleFunc("/tasks", tasksHandler.HandleReadList).Methods("GET")
	r.HandleFunc("/tasks", tasksHandler.HandleSetDone).Methods("PUT")
	c := cron.New()
	c.AddFunc("TZ="+timezone+" 0 0 7,19 * * *", tasksHandler.Cron) // 7am and 7pm every day
	c.Start()

	// Add middleware (auth + logging)
	o := httpauth.AuthOptions{
		AuthFunc: authUser,
	}
	h := httpauth.BasicAuth(o)(r)
	h = handlers.LoggingHandler(os.Stdout, h)

	log.Fatal(http.ListenAndServeTLS(":443", s.Cert, s.Key, h))
}
