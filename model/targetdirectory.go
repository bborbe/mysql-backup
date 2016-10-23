package model

type TargetDirectory string

func (b TargetDirectory) String() string {
	return string(b)
}
