package topology

import (
	"fmt"
	"infinite-cube/internal/validate"
	"testing"
)

func TestBook(t *testing.T) {
	book := Book(3)
	report := validate.AnalyzeTopology(book)
	for _, issue := range report.Issues {
		fmt.Println(issue)
	}
	if len(report.Issues) != 0 {
		t.Fatalf("Expected no issues, but found %d", len(report.Issues))
	}
}
