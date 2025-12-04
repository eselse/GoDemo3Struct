// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"3-struct/api"
)

// Command-line flags
var (
	create = flag.Bool("create", false, "Create a new bin from a JSON file")
	update = flag.Bool("update", false, "Update an existing bin")
	delete = flag.Bool("delete", false, "Delete a bin by ID")
	get    = flag.Bool("get", false, "Read a bin by ID")
	list   = flag.Bool("list", false, "List all saved bin IDs")

	file    = flag.String("file", "", "Path to JSON file (required for -create and -update)")
	id      = flag.String("id", "", "Bin ID (required for -get, -update, -delete)")
	binName = flag.String("binName", "My Bin", "Display name for the bin (used with -create)")
	saveAs  = flag.String("save", "saved-bins.txt", "File to save bin IDs (used with -create)")
)

func main() {
	flag.Parse()

	// === Get Master Key from environment (secure & standard practice) ===
	key := os.Getenv("JSONBIN_KEY")
	if key == "" {
		log.Fatal("Error: JSONBIN_KEY environment variable not set!\n" +
			"Get your free key at https://jsonbin.io → Create Account → Master Key")
	}

	// Create the client once
	client := api.NewClient(key)

	// Show help if no flags
	if flag.NFlag() == 0 {
		printUsage()
		return
	}

	// === CREATE ===
	if *create {
		if *file == "" {
			log.Fatal("Error: -file is required with -create")
		}
		id, err := client.CreateBin(*file, *binName, *saveAs)
		if err != nil {
			log.Fatalf("Failed to create bin: %v", err)
		}
		fmt.Printf("Success! Created bin → ID: %s\n", id)
		return
	}

	// === READ (Get) ===
	if *get {
		if *id == "" {
			log.Fatal("Error: -id is required with -get")
		}
		data, err := client.Get(*id)
		if err != nil {
			log.Fatalf("Failed to read bin: %v", err)
		}
		fmt.Printf("Bin %s content:\n", *id)
		for i, item := range data.Bins {
			fmt.Printf("  %d. %+v\n", i+1, item)
		}
		return
	}

	// === UPDATE ===
	if *update {
		if *file == "" || *id == "" {
			log.Fatal("Error: both -file and -id are required with -update")
		}
		if err := client.Update(*file, *id); err != nil {
			log.Fatalf("Failed to update bin: %v", err)
		}
		fmt.Printf("Success! Bin %s updated.\n", *id)
		return
	}

	// === DELETE ===
	if *delete {
		if *id == "" {
			log.Fatal("Error: -id is required with -delete")
		}
		if err := client.Delete(*id); err != nil {
			log.Fatalf("Failed to delete bin: %v", err)
		}
		fmt.Printf("Success! Bin %s deleted.\n", *id)
		return
	}

	// === LIST ===
	if *list {
		if err := client.List(*saveAs); err != nil {
			log.Fatalf("Failed to list bins: %v", err)
		}
		return
	}
}

func printUsage() {
	fmt.Printf(`JSONBin.io CLI Tool — Simple & Powerful

Usage:
  go run . -create -file data.json                  → create new bin
  go run . -get -id 67f1a2b3c4d5e6f7                 → read bin
  go run . -update -file new.json -id 67f1a2b3...   → update bin
  go run . -delete -id 67f1a2b3c4d5e6f7             → delete bin
  go run . -list                                     → show saved bins

Optional flags:
  -binName "My Data"     → name shown in JSONBin.io dashboard
  -save saved-bins.txt   → file to store bin IDs (default: saved-bins.txt)

Setup:
  export JSONBIN_KEY="your-master-key-here"
  Get free key: https://jsonbin.io

Examples:
  export JSONBIN_KEY=$2a$10$abc123...
  go run . -create -file students.json -binName "Class 2025"
  go run . -list
`)
}
