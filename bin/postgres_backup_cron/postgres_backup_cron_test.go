package main

import (
	"testing"

	"bytes"

	"time"

	. "github.com/bborbe/assert"
)

func TestDoFail(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, func(host string, port int, user string, pass string, database string) error {
		return nil
	}, "", 0, "", "", "", time.Minute, false)
	if err = AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestDoSuccess(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, func(host string, port int, user string, pass string, database string) error {
		return nil
	}, "host", 5432, "user", "pass", "db", time.Minute, true)
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}
