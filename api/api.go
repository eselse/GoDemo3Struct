// Package api provides a simple, easy-to-use client for JSONBin.io (v3 API).
//
// This package allows you to:
//   - Create new bins from local JSON files
//   - Read existing bins by ID
//   - Update bin contents
//   - Delete bins
//   - Save and list created bin IDs locally for quick reuse
//
// It is designed to be beginner-friendly while following Go best practices:
// clean error handling, clear naming, and separation of concerns.
//
// The client works with your existing file.DB interface (from package file)
// and uses the JSONBin.io Master Key for authentication.
//
// Example:
//
//	client := api.NewClient("your-master-key")
//	id, err := client.CreateBin("data.json", "My Backup", "saved-bins.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Created bin ID:", id)
//
// This package is perfect for small tools, backups, CLI apps,
// or learning how to interact with REST APIs in Go.
//
// See the examples folder or the main function in the README for full usage.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"3-struct/bins"
	"3-struct/file"
)

// Client is your main tool to work with JSONBin.io
type Client struct {
	MasterKey string
	fileDB    file.DB // works perfectly with file.NewFileDB()
}

// NewClient creates a new JSONBin.io client
func NewClient(key string) *Client {
	return &Client{
		MasterKey: key,
		fileDB:    file.NewFileDB(),
	}
}

// closeBody drains the response body and closes it safely.
// This allows HTTP connection reuse and satisfies all linters.
func (c *Client) closeBody(rc io.ReadCloser) {
	if rc != nil {
		_, _ = io.Copy(io.Discard, rc) // drain for connection reuse
		_ = rc.Close()                 // acknowledge error return — no warnings!
	}
}

// CreateBin creates a new bin and optionally saves its ID locally
func (c *Client) CreateBin(localFile, binName, saveAs string) (string, error) {
	data, err := c.fileDB.ReadJSON(localFile)
	if err != nil {
		return "", fmt.Errorf("cannot read file %s: %w", localFile, err)
	}

	req, err := http.NewRequest("POST", "https://api.jsonbin.io/v3/b", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", c.MasterKey)
	req.Header.Set("X-Bin-Name", binName)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer c.closeBody(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create bin: status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Metadata struct {
			ID string `json:"id"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	binID := result.Metadata.ID
	fmt.Printf("Bin created! Name: %s → ID: %s\n", binName, binID)

	if saveAs != "" {
		if err := c.fileDB.Append([]byte(binID+"\n"), saveAs); err != nil {
			fmt.Printf("Warning: bin created but failed to save ID: %v\n", err)
		} else {
			fmt.Printf("ID saved to %s\n", saveAs)
		}
	}

	return binID, nil
}

// Get reads a bin by its ID and returns the data
func (c *Client) Get(id string) (bins.BinList, error) {
	req, err := http.NewRequest("GET", "https://api.jsonbin.io/v3/b/"+id, nil)
	if err != nil {
		return bins.BinList{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Master-Key", c.MasterKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return bins.BinList{}, fmt.Errorf("request failed: %w", err)
	}
	defer c.closeBody(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return bins.BinList{}, fmt.Errorf("failed to read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return bins.BinList{}, fmt.Errorf("failed to read bin (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Record bins.BinList `json:"record"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return bins.BinList{}, fmt.Errorf("failed to unmarshal bin data: %w", err)
	}

	return result.Record, nil
}

// Update replaces the content of an existing bin
func (c *Client) Update(localFile, id string) error {
	data, err := c.fileDB.ReadJSON(localFile)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	req, err := http.NewRequest("PUT", "https://api.jsonbin.io/v3/b/"+id, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", c.MasterKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer c.closeBody(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Bin %s updated successfully!\n", id)
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("update failed (status %d): %s", resp.StatusCode, string(body))
}

// Delete removes a bin permanently
func (c *Client) Delete(id string) error {
	req, err := http.NewRequest("DELETE", "https://api.jsonbin.io/v3/b/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Master-Key", c.MasterKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer c.closeBody(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Bin %s deleted successfully!\n", id)
		return nil
	}

	return fmt.Errorf("delete failed (status %d)", resp.StatusCode)
}

// List shows all saved bin IDs from your local file
func (c *Client) List(saveFile string) error {
	data, err := c.fileDB.ReadPlain(saveFile)
	if err != nil {
		return fmt.Errorf("cannot read saved bins: %w", err)
	}

	lines := bytes.Split(bytes.TrimSpace(data), []byte("\n"))
	if len(lines) == 0 || (len(lines) == 1 && len(bytes.TrimSpace(lines[0])) == 0) {
		fmt.Println("No saved bins yet.")
		return nil
	}

	fmt.Println("Your saved bins:")
	for i, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			fmt.Printf("  %d. %s\n", i+1, string(line))
		}
	}
	return nil
}
