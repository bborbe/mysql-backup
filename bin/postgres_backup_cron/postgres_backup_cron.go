package main

import (
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/log"
	"io"
	"os"
	"fmt"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"

	PARAMETER_POSTGRES_HOST = "host"
	PARAMETER_POSTGRES_PORT = "port"
	PARAMETER_POSTGRES_DATABASE = "database"
	PARAMETER_POSTGRES_USER = "user"
	PARAMETER_POSTGRES_PASSWORD = "password"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	hostPtr := flag.String(PARAMETER_POSTGRES_HOST, "", "host")
	portPtr := flag.Int(PARAMETER_POSTGRES_PORT, 5432, "port")
	databasePtr := flag.String(PARAMETER_POSTGRES_DATABASE, "", "database")
	userPtr := flag.String(PARAMETER_POSTGRES_USER, "", "user")
	passwordPtr := flag.String(PARAMETER_POSTGRES_PASSWORD, "", "password")

	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	writer := os.Stdout
	err := do(writer, *hostPtr, *portPtr, *userPtr, *passwordPtr, *databasePtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, host string, port int, user string, pass string, database string) error {
	logger.Debug("start")
	if len(host) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_POSTGRES_HOST)
	}
	if port <= 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_POSTGRES_PORT)
	}
	if len(user) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_POSTGRES_USER)
	}
	if len(pass) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_POSTGRES_PASSWORD)
	}
	if len(database) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_POSTGRES_DATABASE)
	}
	logger.Debug("done")
	return nil
}
