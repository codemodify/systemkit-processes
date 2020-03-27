package tests

import (
	"os"
	"testing"

	"github.com/codemodify/systemkit-processes/find"
)

func TestFindProcess(t *testing.T) {
	p, err := find.ProcessByPID(os.Getpid())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if p == nil {
		t.Fatal("should have process")
	}

	if p.PID() != os.Getpid() {
		t.Fatalf("bad: %#v", p.PID())
	}
}

func TestProcesses(t *testing.T) {
	// This test works because there will always be SOME processes
	// running.
	p, err := find.AllProcesses()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(p) <= 0 {
		t.Fatal("should have processes")
	}

	found := false
	for _, p1 := range p {
		if p1.Details().Executable == "go" || p1.Details().Executable == "go.exe" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("should have Go")
	}
}
