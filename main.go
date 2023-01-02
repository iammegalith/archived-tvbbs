package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Name string
	Age  int
}

func main() {
	// Create a map to store the users' information
	users := make(map[string]User)

	// Print a welcome message
	fmt.Println("Welcome to the Go BBS!")

	// Read input from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')

	// Trim the leading and trailing whitespace from the name
	name = strings.TrimSpace(name)

	// Check if the user is already in the map
	user, ok := users[name]
	if !ok {
		// If the user is not in the map, ask for their age
		fmt.Print("Enter your age: ")
		age, _ := reader.ReadString('\n')
		// Convert the age string to an int
		ageInt, _ := strconv.Atoi(age)

		// Add the user to the map
		users[name] = User{Name: name, Age: ageInt}
	} else {
		// If the user is in the map, use the information from the map
		ageInt = user.Age
	}

	// Print a personalized greeting
	fmt.Printf("Hello, %s! You are %d years old.\n", name, ageInt)
}
