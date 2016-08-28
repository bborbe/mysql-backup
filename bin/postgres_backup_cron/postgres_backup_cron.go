package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"runtime"

	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/bborbe/postgres_backup_cron/backup_creator"
	"github.com/golang/glog"
)

const (
	LOCK_NAME                   = "/var/run/postgres_backup_cron.lock"
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
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	backupCreator := backup_creator.New()
	writer := os.Stdout
	err := do(writer, backupCreator.CreateBackup, *hostPtr, *portPtr, *userPtr, *passwordPtr, *databasePtr, *targetDirPtr, *waitPtr, *oneTimePtr, *lockPtr)
	if err != nil {
		glog.Exit(err)
	}
}

func do(writer io.Writer, createBackup CreateBackup, host string, port int, user string, pass string, database string, targetDir string, wait time.Duration, oneTime bool, lockName string) error {
	l := lock.NewLock(lockName)
	if err := l.Lock(); err != nil {
		return err
	}
	defer l.Unlock()
	glog.V(2).Info("backup cleanup cron started")
	defer glog.V(2).Info("backup cleanup cron finished")

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

	glog.V(2).Infof("host: %s, port: %d, user: %s, pass: %s, database: %s, targetDir: %s, wait: %v, oneTime: %v, lockName: %s", host, port, user, pass, database, targetDir, wait, oneTime, lockName)

	for {
		glog.V(2).Infof("backup started")
		if err := createBackup(host, port, user, pass, database, targetDir); err != nil {
			return err
		}
		glog.V(2).Infof("backup completed")

		if oneTime {
			return nil
		}

		glog.V(2).Infof("wait %v", wait)
		time.Sleep(wait)
		glog.V(2).Infof("sleep done")
	}
}
