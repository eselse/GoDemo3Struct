package main

import (
	"flag"
	"fmt"

	"3-struct/api"
)

type flagsStruct struct {
	create  bool
	update  bool
	delete  bool
	get     bool
	list    bool
	file    string
	name    string
	id      string
	binName string
}

func defineFlags() flagsStruct {
	// Define flags
	create := flag.Bool("create", false, "Create a new bin")
	update := flag.Bool("update", false, "Update an existing bin")
	delete := flag.Bool("delete", false, "Delete a bin")
	get := flag.Bool("get", false, "Get a bin by ID")
	list := flag.Bool("list", false, "List all stored bins")
	file := flag.String("file", "", "JSON file to use")
	name := flag.String("name", "", "Name for the file with saved bins")
	binName := flag.String("binName", "binName", "Name for the bin")
	id := flag.String("id", "", "Bin ID")
	flag.Parse()
	return flagsStruct{
		create:  *create,
		update:  *update,
		delete:  *delete,
		get:     *get,
		list:    *list,
		file:    *file,
		name:    *name,
		id:      *id,
		binName: *binName,
	}
}

func main() {
	// Load configuration
	config := api.InitAPI()

	// Execute commands
	// Create a bin from a file
	flags := defineFlags()
	if flags.create && flags.file != "" && flags.name != "" {
		err := api.CreateBin(flags.file, flags.binName, flags.name, config.Key)
		if err != nil {
			printError(err)
		}
	}

	// Update a bin
	if flags.update && flags.file != "" && flags.id != "" {
		fmt.Printf("Update a bin, from a file, named %s with id %s\n", flags.file, flags.id)
	}

	// Delete a bin
	if flags.delete && flags.id != "" {
		fmt.Printf("Delete a bin with id %s\n", flags.id)
	}

	// Get a bin
	if flags.get && flags.id != "" {
		fmt.Printf("Get bin with id %s\n", flags.id)
		result := api.Get(flags.id, config.Key)
		fmt.Println(result)
	}

	// Get a list of bins
	if flags.list {
		// fmt.Println("Get list of all bins from a file")
		api.List()
	}
}

func printError(err error) {
	fmt.Println(err.Error())
}
