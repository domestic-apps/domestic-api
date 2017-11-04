package main

import (
	"bufio"
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

func addUser() {
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
	// TODO
	fmt.Println("This would be the part where we add username and passHash to the database: " + username)
}

func authUser(username string, password string, r *http.Request) bool {
	if password == "" {
		// TODO log
		return false
	}
	// TODO
	fmt.Println("This would be the part where we fetch the passHash from the database using this username: " + username)
	//passHash := []byte("wowo")
	//err := bcrypt.CompareHashAndPassword(passHash, []byte(password))
	//if err != nil {
	// TODO temporary
	if password != "wowo" {
		// TODO log bad password
		return false
	}
	return true
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
