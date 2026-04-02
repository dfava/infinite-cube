package topology

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGoldenParametric(t *testing.T) {
	tests := []struct {
		name string
		top  any
	}{
		{"htree-2", HTree(2)},
		{"book-3", Book(3)},
		{"snake-4", SnakeChain(4)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := json.MarshalIndent(tt.top, "", "  ")
			if err != nil {
				t.Fatalf("failed to marshal topology: %v", err)
			}

			goldenPath := filepath.Join("testdata", tt.name+".golden")

			// If golden file doesn't exist, create it (first run or new test)
			if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
				if err := os.WriteFile(goldenPath, actual, 0644); err != nil {
					t.Fatalf("failed to write golden file: %v", err)
				}
				t.Logf("Created golden file %s", goldenPath)
				return
			}

			expected, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("failed to read golden file: %v", err)
			}

			if string(actual) != string(expected) {
				t.Errorf("Golden mismatch for %s. If this change is intentional, delete the golden file and re-run the test.", tt.name)
			}
		})
	}
}
