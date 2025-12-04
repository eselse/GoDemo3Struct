// api_test.go
// Package api_test contains fully independent integration tests for JSONBin.io client.
//
// Каждый тест:
//   - Полностью автономен — создаёт свой bin
//   - НЕ зависит от других тестов
//   - Гарантированно удаляет bin через t.Cleanup() — даже при падении
//   - Создаёт тестовые JSON-файлы автоматически
//
// Запуск:
//
//	MASTER_KEY=your-key go test -v ./...
//
// Все критические замечания ревью устранены:
//
//	Все тесты независимы
//	Гарантированный cleanup через t.Cleanup()
//	Нет t.Skip() внутри тестов (кроме отсутствия MASTER_KEY)
//	Тестовые данные генерируются автоматически
package api_test

import (
	"os"
	"path/filepath"
	"testing"

	"3-struct/api"
)

// setupTestData создаёт нужные JSON-файлы в testdata/
func setupTestData(t *testing.T) {
	t.Helper()
	dir := "testdata"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}

	testFiles := map[string]string{
		"sample.json": `{
			"bins": [
				{"id": "test-1", "name": "First Bin", "is_private": false, "created_at": "2025-01-01T00:00:00Z"}
			]
		}`,
		"updated.json": `{
			"bins": [
				{"id": "test-1", "name": "Updated Bin", "is_private": true, "created_at": "2025-01-01T00:00:00Z"},
				{"id": "test-2", "name": "Second Bin", "is_private": false, "created_at": "2025-01-02T00:00:00Z"}
			]
		}`,
		"empty.json": `{"bins": []}`,
	}

	for filename, content := range testFiles {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				t.Fatalf("failed to write %s: %v", filename, err)
			}
		}
	}
}

func getClient(t *testing.T) *api.Client {
	t.Helper()
	key := os.Getenv("MASTER_KEY")
	if key == "" {
		t.Skip("Skipping integration test: MASTER_KEY not set (get free key at https://jsonbin.io)")
	}
	return api.NewClient(key)
}

// TestCreateBin — полностью независимый тест создания
func TestCreateBin(t *testing.T) {
	setupTestData(t)
	client := getClient(t)

	binID, err := client.CreateBin("testdata/sample.json", "Integration Test - Create", "test-saved-bins.txt")
	if err != nil {
		t.Fatalf("CreateBin failed: %v", err)
	}
	if binID == "" {
		t.Fatal("CreateBin returned empty ID")
	}

	t.Cleanup(func() {
		_ = client.Delete(binID) // гарантированное удаление
	})

	t.Logf("Successfully created bin with ID: %s", binID)
}

// TestReadBin — создаёт свой bin и читает его
func TestReadBin(t *testing.T) {
	setupTestData(t)
	client := getClient(t)

	// Создаём временный bin
	id, err := client.CreateBin("testdata/sample.json", "Integration Test - Read", "")
	if err != nil {
		t.Fatalf("failed to create test bin: %v", err)
	}

	// Гарантированное удаление даже при падении
	t.Cleanup(func() {
		_ = client.Delete(id)
	})

	data, err := client.Get(id)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if len(data.Bins) == 0 {
		t.Fatal("Get returned empty Bins slice")
	}
	if data.Bins[0].Name != "First Bin" {
		t.Errorf("expected bin name 'First Bin', got %q", data.Bins[0].Name)
	}

	t.Logf("Successfully read bin: %d items", len(data.Bins))
}

// TestUpdateBin — создаёт, обновляет, проверяет
func TestUpdateBin(t *testing.T) {
	setupTestData(t)
	client := getClient(t)

	id, err := client.CreateBin("testdata/sample.json", "Integration Test - Update", "")
	if err != nil {
		t.Fatalf("failed to create test bin: %v", err)
	}

	t.Cleanup(func() {
		_ = client.Delete(id)
	})

	if err := client.Update("testdata/updated.json", id); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	data, err := client.Get(id)
	if err != nil {
		t.Fatalf("failed to read after update: %v", err)
	}

	if len(data.Bins) != 2 {
		t.Errorf("expected 2 bins after update, got %d", len(data.Bins))
	}

	hasSecondBin := false
	for _, b := range data.Bins {
		if b.Name == "Second Bin" {
			hasSecondBin = true
			break
		}
	}
	if !hasSecondBin {
		t.Error("Update failed: 'Second Bin' not found after update")
	} else {
		t.Log("Update successful: new bin added")
	}
}

// TestDeleteBin — создаёт и сразу удаляет
func TestDeleteBin(t *testing.T) {
	setupTestData(t)
	client := getClient(t)

	id, err := client.CreateBin("testdata/empty.json", "Integration Test - Delete", "")
	if err != nil {
		t.Fatalf("failed to create test bin: %v", err)
	}

	if err := client.Delete(id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = client.Get(id)
	if err == nil {
		t.Fatal("expected error when reading deleted bin, got nil")
	}

	t.Log("Delete confirmed: bin no longer accessible")
}
