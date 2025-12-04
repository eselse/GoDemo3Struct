package env

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err // file not found is acceptable in many cases
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			// Log it
			log.Printf("warning: failed to close file: %v", closeErr)
			// And optionally propagate if no other error occurred
			if err == nil {
				err = closeErr
			}
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle inline comments: KEY=value # comment
		if strings.Contains(line, "#") {
			line = strings.Split(line, "#")[0]
			line = strings.TrimSpace(line)
		}

		// Split on first = only
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // invalid line
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		value = strings.Trim(value, `"'`)

		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}
