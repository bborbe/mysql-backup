package model

type Name string

func (b Name) String() string {
	return string(b)
}
