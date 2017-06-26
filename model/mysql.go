package model

import (
	"strconv"
)

type MysqlHost string

func (p MysqlHost) String() string {
	return string(p)
}

type MysqlPort int

func (p MysqlPort) Int() int {
	return int(p)
}

func (p MysqlPort) String() string {
	return strconv.Itoa(p.Int())
}

type MysqlUser string

func (p MysqlUser) String() string {
	return string(p)
}

type MysqlPassword string

func (p MysqlPassword) String() string {
	return string(p)
}

type MysqlDatabase string

func (p MysqlDatabase) String() string {
	return string(p)
}
