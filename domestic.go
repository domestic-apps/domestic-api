package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"

	"github.com/domestic-apps/domestic-api/chores"
	"github.com/domestic-apps/domestic-api/tasks"

	"gopkg.in/robfig/cron.v2"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
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

func addUser() {
	username, password := credentials()
	passHash, err := bcrypt.GenerateFromPassword(password, 6)
	fmt.Print("Enter Password Again: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Fatal("Could not get the password, oh no!")
	}
	err = bcrypt.CompareHashAndPassword(passHash, bytePassword)
	if err != nil {
		log.Fatal("Passwords did not match, or some other thing went wrong")
	}
	fmt.Println("This would be the part where we add username and passHash to the database: " + username)
}

// credentials gets username and password credentials from the command line.
// This was totally stolen from https://play.golang.org/p/l-9IP1mrhA
func credentials() (string, []byte) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Fatal("Could not get the password, oh no!")
	}

	return strings.TrimSpace(username), bytePassword
}

func runServer() {
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
	r := mux.NewRouter()
	r.HandleFunc("/chores", choresHandler.HandleCreate).Methods("POST").Schemes("https")
	r.HandleFunc("/chores", choresHandler.HandleReadList).Methods("GET").Schemes("https")
	r.HandleFunc("/chores", choresHandler.HandleUpdate).Methods("PUT").Schemes("https")
	r.HandleFunc("/chores", choresHandler.HandleDelete).Methods("DELETE").Schemes("https")
	r.HandleFunc("/tasks/", tasksHandler.Handle)
	c := cron.New()
	c.AddFunc("TZ="+timezone+" 0 0 7,19 * * *", tasksHandler.Cron) // 7am and 7pm every day
	c.Start()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServeTLS(":443", s.Cert, s.Key, nil))
}
