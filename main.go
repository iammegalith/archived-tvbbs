package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name     string
	Age      int
	Password string
}

func main() {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "bbs.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Print a welcome message
	fmt.Println("Welcome to the Go BBS!")

	// Read input from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')

	// Trim the leading and trailing whitespace from the name
	name = strings.TrimSpace(name)

	// Check if the user is already in the database
	var age int
	var password string
	err = db.QueryRow("SELECT age, password FROM users WHERE name=?", name).Scan(&age, &password)
	if err != nil {
		// If the user is not in the database, ask for their age and password
		if err == sql.ErrNoRows {
			fmt.Print("Enter your age: ")
			ageStr, _ := reader.ReadString('\n')
			// Convert the age string to an int
			age, _ = strconv.Atoi(ageStr)

			fmt.Print("Enter your password: ")
			password, _ := reader.ReadString('\n')
			// Trim the leading and trailing whitespace from the password
			password = strings.TrimSpace(password)

			// Add the user to the database
			_, err = db.Exec("INSERT INTO users (name, age, password) VALUES (?, ?, ?)", name, age, password)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			// If there was an error other than sql.ErrNoRows, print the error
			fmt.Println(err)
			return
		}
	} else {
		// If the user is in the database, ask for their password
		fmt.Print("Enter your password: ")
		passwordInput, _ := reader.ReadString('\n')
		// Trim the leading and trailing whitespace from the password
		passwordInput = strings.TrimSpace(passwordInput)

		// Check if the password is correct
		if passwordInput != password {
			fmt.Println("Incorrect password!")
			return
		}
	}

	// Print a personalized greeting
	fmt.Printf("Hello, %s! You are %d years old.\n", name, age)
}
