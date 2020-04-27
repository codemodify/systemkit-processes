package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/codemodify/systemkit-processes/find"
)

func TestProcessByPID(t *testing.T) {
	rp, err := find.ProcessByPID(os.Getpid())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	runtimeProcess := rp.Details()
	if runtimeProcess.ProcessID != os.Getpid() {
		t.Fatalf("bad: %#v", runtimeProcess.ProcessID)
	}

	detailsAsBytes, _ := json.MarshalIndent(runtimeProcess, "", "\t")
	fmt.Println(string(detailsAsBytes))
}
