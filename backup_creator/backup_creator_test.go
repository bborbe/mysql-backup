package backup_creator

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsBackupCreator(t *testing.T) {
	c := New()
	var i *BackupCreator
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
