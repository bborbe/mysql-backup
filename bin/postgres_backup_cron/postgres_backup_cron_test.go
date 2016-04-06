package main

import (
	"testing"

	"bytes"

	. "github.com/bborbe/assert"
)

func TestDoFail(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, "", 0, "", "", "")
	if err = AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestDoSuccess(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, "host", 5432, "user", "pass", "db")
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}
