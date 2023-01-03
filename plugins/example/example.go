package example

import "fmt"

type Module interface {
	Execute() error
}

type ExampleModule struct{}

/*
title: Main Menu
items:
- text: View messages
  action: view_messages
- text: Execute example module
  action: module:/path/to/example_module.so
*/

func (m *ExampleModule) Execute() error {
	fmt.Println("Executing example module!")
	return nil
}

var Module Module

func main() {
	Module = &ExampleModule{}
}
