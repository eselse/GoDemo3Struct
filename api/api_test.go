// api_test.go
// Package api_test contains integration tests for the JSONBin.io client.
//
// These tests perform real HTTP requests using a free JSONBin.io account.
// To run them:
//
//  1. Sign up for free at https://jsonbin.io
//  2. Create a new bin → copy your "Master Key" (starts with $2a$)
//  3. Run: MASTER_KEY=your-key-here go test -v
//
// The tests cover full CRUD flow:
//   - Create a bin
//   - Read it back
//   - Update it
//   - Delete it
//
// This shows real-world API interaction, error handling, and cleanup —
// exactly what recruiters love to see in a portfolio.
package api_test

import (
	"os"
	"testing"

	"3-struct/api"
)

// TestCRUD performs full Create → Read → Update → Delete cycle
func TestCRUD(t *testing.T) {
	// Get master key from environment (safe + required)
	key := os.Getenv("MASTER_KEY")
	if key == "" {
		t.Skip("Skipping integration test: MASTER_KEY not set. Get free key at jsonbin.io")
	}

	client := api.NewClient(key)
	testFile := "testdata/sample.json"
	binName := "Go Client Test Bin"
	saveFile := "test-saved-bins.txt"

	// Clean up any previous test data
	_ = os.Remove(saveFile)

	// 1. CREATE
	t.Run("Create", func(t *testing.T) {
		id, err := client.CreateBin(testFile, binName, saveFile)
		if err != nil {
			t.Fatalf("CreateBin failed: %v", err)
		}
		if id == "" {
			t.Fatal("CreateBin returned empty ID")
		}
		t.Logf("Created bin with ID: %s", id)

		// Save ID for next steps
		err = os.Setenv("TEST_BIN_ID", id)
		if err != nil {
			t.Fatalf("CreateBin failed: %v", err)
		}
	})

	// 2. READ
	t.Run("Read", func(t *testing.T) {
		id := os.Getenv("TEST_BIN_ID")
		if id == "" {
			t.Skip("No bin ID from create step")
		}

		data, err := client.Get(id)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(data.Bins) == 0 {
			t.Fatal("Get returned empty data")
		}
		t.Logf("Successfully read bin: %d items", len(data.Bins))
	})

	// 3. UPDATE
	t.Run("Update", func(t *testing.T) {
		id := os.Getenv("TEST_BIN_ID")
		if id == "" {
			t.Skip("No bin ID from create step")
		}

		updateFile := "testdata/updated.json"
		if err := client.Update(updateFile, id); err != nil {
			t.Fatalf("Update failed: %v", err)
		}
		t.Log("Update successful")
	})

	// 4. DELETE
	t.Run("Delete", func(t *testing.T) {
		id := os.Getenv("TEST_BIN_ID")
		if id == "" {
			t.Skip("No bin ID from create step")
		}

		if err := client.Delete(id); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
		t.Log("Delete successful — cleanup complete")
	})

	// Final cleanup
	_ = os.Remove(saveFile)
}
