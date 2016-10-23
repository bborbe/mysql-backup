package main

import (
	"fmt"
	"time"

	"runtime"

	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/bborbe/postgres_backup_cron/backup"
	"github.com/bborbe/postgres_backup_cron/model"
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

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {
	lockName := *lockPtr
	l := lock.NewLock(lockName)
	if err := l.Lock(); err != nil {
		return err
	}
	defer func() {
		if err := l.Unlock(); err != nil {
			glog.Warningf("unlock failed: %v", err)
		}
	}()

	glog.V(1).Info("backup cleanup cron started")
	defer glog.V(1).Info("backup cleanup cron finished")

	return exec()
}

func exec() error {
	host := model.PostgresqlHost(*hostPtr)
	if len(host) == 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresHost)
	}
	port := model.PostgresqlPort(*portPtr)
	if port <= 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresPort)
	}
	user := model.PostgresqlUser(*userPtr)
	if len(user) == 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresUser)
	}
	pass := model.PostgresqlPassword(*passwordPtr)
	if len(pass) == 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresPassword)
	}
	database := model.PostgresqlDatabase(*databasePtr)
	if len(database) == 0 {
		return fmt.Errorf("parameter %s missing", parameterPostgresDatabase)
	}
	targetDir := model.TargetDirectory(*targetDirPtr)
	if len(targetDir) == 0 {
		return fmt.Errorf("parameter %s missing", parameterTargetDir)
	}

	oneTime := *oneTimePtr
	wait := *waitPtr

	glog.V(1).Infof("host: %s, port: %d, user: %s, password-length: %d, database: %s, targetDir: %s, wait: %v, oneTime: %v, lockName: %s", host, port, user, len(pass), database, targetDir, wait, oneTime, lockName)

	for {
		glog.V(1).Infof("backup started")
		if err := backup.Create(host, port, user, pass, database, targetDir); err != nil {
			glog.Warningf("backup failed: %v", err)
		} else {
			glog.V(1).Infof("backup completed")
		}

		if oneTime {
			glog.V(2).Infof("one time => exit")
			return nil
		}

		glog.V(2).Infof("wait %v", wait)
		time.Sleep(wait)
		glog.V(2).Infof("sleep done")
	}
}
