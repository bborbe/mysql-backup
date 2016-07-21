package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"runtime"

	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/bborbe/log"
	"github.com/bborbe/postgres_backup_cron/backup_creator"
)

const (
	LOCK_NAME                   = "/var/run/postgres_backup_cron.lock"
	PARAMETER_LOGLEVEL          = "loglevel"
	PARAMETER_POSTGRES_HOST     = "host"
	PARAMETER_POSTGRES_PORT     = "port"
	PARAMETER_POSTGRES_DATABASE = "database"
	PARAMETER_POSTGRES_USER     = "username"
	PARAMETER_POSTGRES_PASSWORD = "password"
	PARAMETER_TARGET_DIR        = "targetdir"
	PARAMETER_WAIT              = "wait"
	PARAMETER_ONE_TIME          = "one-time"
	PARAMETER_LOCK              = "lock"
)

var (
	logger       = log.DefaultLogger
	logLevelPtr  = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	hostPtr      = flag.String(PARAMETER_POSTGRES_HOST, "", "host")
	portPtr      = flag.Int(PARAMETER_POSTGRES_PORT, 5432, "port")
	databasePtr  = flag.String(PARAMETER_POSTGRES_DATABASE, "", "database")
	userPtr      = flag.String(PARAMETER_POSTGRES_USER, "", "username")
	passwordPtr  = flag.String(PARAMETER_POSTGRES_PASSWORD, "", "password")
	waitPtr      = flag.Duration(PARAMETER_WAIT, time.Minute*60, "wait")
	oneTimePtr   = flag.Bool(PARAMETER_ONE_TIME, false, "exit after first backup")
	targetDirPtr = flag.String(PARAMETER_TARGET_DIR, "", "target directory")
	lockPtr      = flag.String(PARAMETER_LOCK, LOCK_NAME, "lock")
)

type CreateBackup func(host string, port int, user string, pass string, database string, targetDir string) error

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	backupCreator := backup_creator.New()
	writer := os.Stdout
	err := do(writer, backupCreator.CreateBackup, *hostPtr, *portPtr, *userPtr, *passwordPtr, *databasePtr, *targetDirPtr, *waitPtr, *oneTimePtr, *lockPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, createBackup CreateBackup, host string, port int, user string, pass string, database string, targetDir string, wait time.Duration, oneTime bool, lockName string) error {
	l := lock.NewLock(lockName)
	if err := l.Lock(); err != nil {
		return err
	}
	defer l.Unlock()
	logger.Debug("backup cleanup cron started")
	defer logger.Debug("backup cleanup cron finished")

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
	if len(targetDir) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_TARGET_DIR)
	}

	logger.Debugf("host: %s, port: %d, user: %s, pass: %s, database: %s, targetDir: %s, wait: %v, oneTime: %v, lockName: %s", host, port, user, pass, database, targetDir, wait, oneTime, lockName)

	for {
		logger.Debugf("backup started")
		if err := createBackup(host, port, user, pass, database, targetDir); err != nil {
			return err
		}
		logger.Debugf("backup completed")

		if oneTime {
			return nil
		}

		logger.Debugf("wait %v", wait)
		time.Sleep(wait)
		logger.Debugf("sleep done")
	}
}
