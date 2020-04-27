package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/codemodify/systemkit-processes/find"
)

func TestAllProcesses(t *testing.T) {
	allProcesses, err := find.AllProcesses()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(allProcesses) <= 0 {
		t.Fatal("should have processes")
	}

	for _, proc := range allProcesses {
		detailsAsBytes, _ := json.Marshal(proc.Details())
		fmt.Println(string(detailsAsBytes))
	}
}
