package model

import (
	"fmt"
	"os"

	"github.com/bborbe/io/util"
)

type TargetDirectory string

func (b TargetDirectory) String() string {
	return string(b)
}

func (t TargetDirectory) Mkdir(perm os.FileMode) error {
	path, err := util.NormalizePath(t.String())
	if err != nil {
		return fmt.Errorf("normalize path '%s' failed: %v", t, err)
	}
	return os.MkdirAll(path, perm)
}
