package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	"plugin"
	"strconv"
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
	Text   string
	Action string
}

type Menu struct {
	Title string
	Items []MenuItem
}

type Message struct {
	ID      int
	Thread  int
	Author  string
	Subject string
	Body    string
	Date    string
}

type Module interface {
	Execute() error
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
	// Load the main menu from a YAML file
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
		if item.Action == "view_messages" {
			// Retrieve the message threads from the database
			rows, err := db.Query("SELECT id, subject FROM threads")
			if err != nil {
				fmt.Fprintln(conn, err)
				return
			}
			defer rows.Close()

			// Print the message threads
			fmt.Fprintln(conn, "Message threads:")
			for rows.Next() {
				var id int
				var subject string
				if err := rows.Scan(&id, &subject); err != nil {
					fmt.Fprintln(conn, err)
					return
				}
				fmt.Fprintf(conn, "%d. %s\n", id, subject)
			}
			if err := rows.Err(); err != nil {
				fmt.Fprintln(conn, err)
				return
			}

			// Read the user's selection
			fmt.Fprint(conn, "Enter your selection: ")
			threadStr, _ := reader.ReadString('\n')
			// Convert the thread string to an int
			thread, _ := strconv.Atoi(threadStr)

			// Retrieve the messages in the selected thread from the database
			rows, err = db.Query("SELECT id, author, subject, body, date FROM messages WHERE thread=?", thread)
			if err != nil {
				fmt.Fprintln(conn, err)
				return
			}
			defer rows.Close()

			// Print the messages
			fmt.Fprintln(conn, "Messages:")
			for rows.Next() {
				var message Message
				if err := rows.Scan(&message.ID, &message.Author, &message.Subject, &message.Body, &message.Date); err != nil {
					fmt.Fprintln(conn, err)
					return
				}
				fmt.Fprintf(conn, "%d. %s - %s\n%s\n", message.ID, message.Author, message.Subject, message.Body)
			}
			if err := rows.Err(); err != nil {
				fmt.Fprintln(conn, err)
				return
			}
		} else if strings.HasPrefix(item.Action, "module:") {
			// Load the module from a file
			modulePath := strings.TrimPrefix(item.Action, "module:")
			p, err := plugin.Open(modulePath)
			if err != nil {
				fmt.Fprintln(conn, err)
				return
			}

			// Look up the Module symbol in the plugin
			symModule, err := p.Lookup("Module")
			if err != nil {
				fmt.Fprintln(conn, err)
				return
			}

			// Assert that the symbol is of the correct type
			module, ok := symModule.(Module)
			if !ok {
				fmt.Fprintln(conn, "invalid module type")
				return
			}

			// Execute the module
			if err := module.Execute(); err != nil {
				fmt.Fprintln(conn, err)
				return
			}
		} else {
			fmt.Fprintf(conn, "Performing action: %s\n", item.Action)
		}
	} else {
		fmt.Fprintln(conn, "Invalid selection!")
	}
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
