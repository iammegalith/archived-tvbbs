package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name     string
	Age      int
	Password string
}

type MenuItem struct {
	Text string
	Action string
}

type Menu struct {
	Title string
	Items []MenuItem
}

func handleClient(conn net.Conn) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "bbs.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Check if the terminal supports ANSI escape codes
	if terminal.IsTerminal(int(conn.Fd())) {
		fmt.Println("ANSI escape codes are supported.")
	} else {
		fmt.Println("ANSI escape codes are not supported.")
	}

	// Print a welcome message
	fmt.Fprintln(conn, "Welcome to the Go BBS!")

	// Read input from the client
	reader := bufio.NewReader(conn)
	fmt.Fprint(conn, "Enter your name: ")
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
			fmt.Fprint(conn, "Enter your age: ")
			ageStr, _ := reader.ReadString('\n')
			// Convert the age string to an int
			age, _ = strconv.Atoi(ageStr)

			fmt.Fprint(conn, "Enter your password: ")
			password, _ := reader.ReadString('\n')
			// Trim the leading and trailing whitespace from the password
			password = strings.TrimSpace(password)

			// Add the user to the database
			_, err = db.Exec("INSERT INTO users (name, age, password) VALUES (?, ?, ?)", name, age, password)				// Load the main menu from a YAML file
				menuData, err := ioutil.ReadFile("menu.yaml")
				if err != nil {
					fmt.Fprintln(conn, err)
					return
				}
			
				var menu Menu
				err = yaml.Unmarshal(menuData, &menu)
				if err != nil {
					fmt.Fprintln(conn, err)
					return
				}
			
				// Print the main menu
				fmt.Fprintln(conn, menu.Title)
				for i, item := range menu.Items {
					fmt.Fprintf(conn, "%d. %s\n", i+1, item.Text)
				}
			
				// Read the user's selection
				fmt.Fprint(conn, "Enter your selection: ")
				selectionStr, _ := reader.ReadString('\n')
				// Convert the selection string to an int
				selection, _ := strconv.Atoi(selectionStr)
			
				// Perform the selected action
				if selection > 0 && selection <= len(menu.Items) {
					item := menu.Items[selection-1]
					fmt.Fprintf(conn, "Performing action: %s\n", item.Action)
				} else {
					fmt.Fprintln(conn, "Invalid selection!")
				}
			
				// Print a personalized greeting
				fmt.Fprintf(conn, "Hello, %s! You are %d years old.\n", name, age)
			}
			
			func main() {
				// Create a Telnet server
				ln, err := net.Listen("tcp", ":23")
				if err != nil {
					fmt.Println(err)
					return
				}
				defer ln.Close()
			
				fmt.Println("BBS Telnet server listening on port 23...")
			
				// Accept incoming connections and handle them in a new goroutine
				for {
					conn, err := ln.Accept()
					if err != nil {
						fmt.Println(err)
						continue
					}
					go handleClient(conn)
				}
			}
			