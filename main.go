package main

import (
	"flag"
	"fmt"

	"3-struct/storage"
)

func main() {
	// Define flags
	create := flag.Bool("create", false, "Create a new bin")
	update := flag.Bool("update", false, "Update an existing bin")
	delete := flag.Bool("delete", false, "Delete a bin")
	get := flag.Bool("get", false, "Get a bin by ID")
	list := flag.Bool("list", false, "List all stored bins")
	file := flag.String("file", "", "JSON file to use")
	name := flag.String("name", "", "Name for the bin")
	id := flag.String("id", "", "Bin ID")
	flag.Parse()

	// Load configuration
	// config := api.InitAPI()

	// fmt.Println("Key loaded", config.Key)

	// Execute commands
	// Create a bin from a file
	if *create && *file != "" && *name != "" {
		fmt.Printf("Create a bin, named %s from a file, named %s\n", *file, *name)
	}

	// Update a bin
	if *update && *file != "" && *id != "" {
		fmt.Printf("Update a bin, from a file, named %s with id %s\n", *file, *id)
	}

	// Delete a bin
	if *delete && *id != "" {
		fmt.Printf("Delete a bin with id %s\n", *id)
	}

	// Get a bin
	if *get && *id != "" {
		fmt.Printf("Get bin with id %s\n", *id)
	}

	// Get a list of bins
	if *list {
		fmt.Println("Get list of all bins from a file")
	}

	someString := "hello there"
	fileName := "file.txt"
	_, err := storage.SaveToFile([]byte(someString), fileName)
	if err != nil {
		printError(err)
	}
	// resultString, err := storage.ReadFile(fileName)
	// if err != nil {
	// 	printError(err)
	// }
	// fmt.Println(string(resultString))
}

func printError(err error) {
	fmt.Println(err.Error())
}
