package model

import (
	"fmt"
	"github.com/bborbe/io/util"
	"os"
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
