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
	Name string
	Age  int
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
	err = db.QueryRow("SELECT age FROM users WHERE name=?", name).Scan(&age)
	if err != nil {
		// If the user is not in the database, ask for their age
		if err == sql.ErrNoRows {
			fmt.Print("Enter your age: ")
			ageStr, _ := reader.ReadString('\n')
			// Convert the age string to an int
			age, _ = strconv.Atoi(ageStr)

			// Add the user to the database
			_, err = db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", name, age)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			// If there was an error other than sql.ErrNoRows, print the error
			fmt.Println(err)
			return
		}
	}

	// Print a personalized greeting
	fmt.Printf("Hello, %s! You are %d years old.\n", name, age)
}
