package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Print a welcome message
	fmt.Println("Welcome to the Go BBS!")

	// Read input from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')

	// Print a personalized greeting
	fmt.Printf("Hello, %s!\n", name)
}
