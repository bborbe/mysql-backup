package main

import (
	"fmt"
	"time"

	"runtime"

	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/bborbe/postgres_backup_cron/backup_creator"
	"github.com/golang/glog"
)

const (
	lockName                  = "/var/run/postgres_backup_cron.lock"
	parameterPostgresHost     = "host"
	parameterPostgresPort     = "port"
	parameterPostgresDatabase = "database"
	parameterPostgresUser     = "username"
	parameterPostgresPassword = "password"
	parameterTargetDir        = "targetdir"
	parameterWait             = "wait"
	parameterOneTime          = "one-time"
	parameterLock             = "lock"
)

var (
	hostPtr      = flag.String(parameterPostgresHost, "", "host")
	portPtr      = flag.Int(parameterPostgresPort, 5432, "port")
	databasePtr  = flag.String(parameterPostgresDatabase, "", "database")
	userPtr      = flag.String(parameterPostgresUser, "", "username")
	passwordPtr  = flag.String(parameterPostgresPassword, "", "password")
	waitPtr      = flag.Duration(parameterWait, time.Minute*60, "wait")
	oneTimePtr   = flag.Bool(parameterOneTime, false, "exit after first backup")
	targetDirPtr = flag.String(parameterTargetDir, "", "target directory")
	lockPtr      = flag.String(parameterLock, lockName, "lock")
)

type CreateBackup func(host string, port int, user string, pass string, database string, targetDir string) error

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	backupCreator := backup_creator.New()
	err := do(
		backupCreator.CreateBackup,
		*hostPtr,
		*portPtr,
		*userPtr,
		*passwordPtr,
		*databasePtr,
		*targetDirPtr,
		*waitPtr,
		*oneTimePtr,
		*lockPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	createBackup CreateBackup,
	host string,
	port int,
	user string,
	pass string,
	database string,
	targetDir string,
	wait time.Duration,
	oneTime bool,
	lockName string,
) error {
	l := lock.NewLock(lockName)
	if err := l.Lock(); err != nil {
		return err
	}
	defer l.Unlock()
	glog.V(2).Info("backup cleanup cron started")
	defer glog.V(2).Info("backup cleanup cron finished")

	if len(host) == 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresHost)
	}
	if port <= 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresPort)
	}
	if len(user) == 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresUser)
	}
	if len(pass) == 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresPassword)
	}
	if len(database) == 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresDatabase)
	}
	if len(targetDir) == 0 {
		return fmt.Errorf("parameter %s missing", parameterTargetDir)
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
