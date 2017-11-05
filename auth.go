package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

const hashNum = 6

func addUser(db *sql.DB) {
	username, password := credentials()
	passHash, err := bcrypt.GenerateFromPassword(password, hashNum)
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
	stmt, _ := db.Prepare("INSERT INTO users(username, password) VALUES(?,?)")
	_, err = stmt.Exec(username, passHash)
	if err == nil {
		fmt.Println("Successfully added user!")
	}
}

func authUser(db *sql.DB) func(string, string, *http.Request) bool {
	stmt, err := db.Prepare("SELECT password from users where username = ?")
	if err != nil {
		log.Fatalf("Stmt generation failed, here's the error: %v", err)
		return func(username string, password string, r *http.Request) bool {
			return false
		}
	}
	return func(username string, password string, r *http.Request) bool {
		if password == "" {
			// TODO log no password
			return false
		}
		rows, err := stmt.Query(username)
		defer rows.Close()
		var passHash []byte

		for rows.Next() {
			err := rows.Scan(&passHash)
			if err != nil {
				return false
			}
		}
		err = bcrypt.CompareHashAndPassword(passHash, []byte(password))
		if err != nil {
			// TODO log bad password
			return false
		}
		return true
	}
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
