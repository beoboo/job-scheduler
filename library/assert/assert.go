package assert

import (
	"fmt"
	"github.com/beoboo/job-scheduler/library/stream"
	"testing"
)

func AssertStatus(t *testing.T, st string, expected string) {
	if st != expected {
		t.Fatalf("Job status should be \"%s\", got \"%s\"", expected, st)
	}
}

func AssertOutput(t *testing.T, o *stream.Stream, expected []string) {
	fmt.Println("Checking output")

	for _, e := range expected {
		l, err := o.Read()
		if err != nil {
			t.Fatal("Expected output line")
		}

		if l.Text != e {
			t.Fatalf("Job output should contain \"%s\", got \"%s\"", e, l.Text)
		}
	}
}
