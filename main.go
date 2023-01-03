package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	"plugin"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	_ "github.com/mattn/go-sqlite3"
)

const (
	TVBBS = "Television BBS"
)
const (
	TVBBS_VERSION = "0.0.1"
)

type User struct {
	Handle   string
	Level    int
	Password string
}

type MenuItem struct {
	Text    string
	Action  string
	Submenu string
}

type Menu struct {
	Title   string
	Items   []MenuItem
	Submenu map[string]Menu
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

func displayMenu(conn net.Conn, menu Menu) {
	// Print the menu title
	fmt.Fprintln(conn, menu.Title)

	// Print the menu items
	for i, item := range menu.Items {
		fmt.Fprintf(conn, "%d. %s\n", i+1, item.Text)
	}
}
func handleClient(conn net.Conn) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "bbs.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// check if client supports ansi - this needs to be developed
	fmt.Fprint(conn, "[Checking for ANSi support: Press ENTER | RETURN ]\r\n")
	fmt.Fprint(conn, "\x1b[31m")
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}
	if strings.Contains(response, "\x1b[31m") {
		fmt.Println("ANSI supported")
		//var hasANSI = true
	} else {
		fmt.Println("ANSI not supported")
		//var hasANSI = false
	}
	// Print a welcome message
	var banner = fmt.Sprintf("%s %s %s", TVBBS, TVBBS_VERSION, "\r\n")
	fmt.Fprintln(conn, banner)

	// Read input from the client
	reader := bufio.NewReader(conn)
	fmt.Fprint(conn, "Username: ")
	name, _ := reader.ReadString('\n')

	// Trim the leading and trailing whitespace from the name
	name = strings.TrimSpace(name)

	// Check if the user is already in the database
	var level int
	var password string
	err = db.QueryRow("SELECT password FROM users WHERE name=?", name).Scan(&password)
	if err != sql.ErrNoRows {
		level = 1
		// should add a way to allow user to say oops - i meant this other username
		fmt.Fprint(conn, "Enter your password: ")
		password, _ := reader.ReadString('\n')
		// Trim the leading and trailing whitespace from the password
		password = strings.TrimSpace(password)

		// Add the user to the database
		_, err = db.Exec("INSERT INTO users (name, level, password) VALUES (?, ?, ?)", name, level, password)
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
	fmt.Fprintf(conn, "Hello, %s! \n", name)
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
