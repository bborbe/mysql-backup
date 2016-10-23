package model

import "strconv"

type PostgresqlHost string

func (p PostgresqlHost) String() string {
	return string(p)
}

type PostgresqlPort int

func (p PostgresqlPort) Int() int {
	return int(p)
}

func (p PostgresqlPort) String() string {
	return strconv.Itoa(p.Int())
}

type PostgresqlUser string

func (p PostgresqlUser) String() string {
	return string(p)
}

type PostgresqlPassword string

func (p PostgresqlPassword) String() string {
	return string(p)
}

type PostgresqlDatabase string

func (p PostgresqlDatabase) String() string {
	return string(p)
}
