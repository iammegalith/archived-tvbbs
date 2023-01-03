package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name     string
	Age      int
	Password string
}

func handleClient(conn net.Conn) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "bbs.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// get file descriptor from net.TCPConn
	file, err := conn.(*net.TCPConn).File()
	if err != nil {
		fmt.Println(err)
		return
	}
	fd := file.Fd()
	// Check if the terminal supports ANSI escape codes
	if terminal.IsTerminal(int(fd)) {
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
			_, err = db.Exec("INSERT INTO users (name, age, password) VALUES (?, ?, ?)", name, age, password)
			if err != nil {
				fmt.Fprintln(conn, err)
				return
			}
		} else {
			// If there was an error other than sql.ErrNoRows, print the error
			fmt.Fprintln(conn, err)
			return
		}
	} else {
		// If the user is in the database, ask for their password
		fmt.Fprint(conn, "Enter your password: ")
		passwordInput, _ := reader.ReadString('\n')
		// Trim the leading and trailing whitespace from the password
		passwordInput = strings.TrimSpace(passwordInput)

		// Check if the password is correct
		if passwordInput != password {
			fmt.Fprintln(conn, "Incorrect password!")
			return
		}
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
