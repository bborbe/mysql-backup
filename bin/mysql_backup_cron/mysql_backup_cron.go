package main

import (
	"fmt"
	"time"

	"runtime"

	"context"
	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/bborbe/mysql_backup_cron/backup"
	"github.com/bborbe/mysql_backup_cron/model"
	"github.com/golang/glog"
)

const (
	defaultLockName        = "/var/run/mysql_backup_cron.lock"
	defaultName            = "mysql"
	parameterMysqlHost     = "host"
	parameterMysqlPort     = "port"
	parameterMysqlDatabase = "database"
	parameterMysqlUser     = "username"
	parameterMysqlPassword = "password"
	parameterTargetDir     = "targetdir"
	parameterWait          = "wait"
	parameterOneTime       = "one-time"
	parameterLock          = "lock"
	parameterName          = "name"
)

var (
	hostPtr      = flag.String(parameterMysqlHost, "", "host")
	portPtr      = flag.Int(parameterMysqlPort, 5432, "port")
	databasePtr  = flag.String(parameterMysqlDatabase, "", "database")
	userPtr      = flag.String(parameterMysqlUser, "", "username")
	passwordPtr  = flag.String(parameterMysqlPassword, "", "password")
	waitPtr      = flag.Duration(parameterWait, time.Minute*60, "wait")
	oneTimePtr   = flag.Bool(parameterOneTime, false, "exit after first backup")
	targetDirPtr = flag.String(parameterTargetDir, "", "target directory")
	lockPtr      = flag.String(parameterLock, defaultLockName, "lock")
	namePtr      = flag.String(parameterName, defaultName, "name")
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

	glog.V(1).Info("backup mysql cron started")
	defer glog.V(1).Info("backup mysql cron finished")

	return exec()
}

func exec() error {
	host := model.MysqlHost(*hostPtr)
	if len(host) == 0 {
		return fmt.Errorf("parameter %s missing", parameterMysqlHost)
	}
	port := model.MysqlPort(*portPtr)
	if port <= 0 {
		return fmt.Errorf("parameter %s missing", parameterMysqlPort)
	}
	user := model.MysqlUser(*userPtr)
	if len(user) == 0 {
		return fmt.Errorf("parameter %s missing", parameterMysqlUser)
	}
	pass := model.MysqlPassword(*passwordPtr)
	if len(pass) == 0 {
		return fmt.Errorf("parameter %s missing", parameterMysqlPassword)
	}
	database := model.MysqlDatabase(*databasePtr)
	if len(database) == 0 {
		return fmt.Errorf("parameter %s missing", parameterMysqlDatabase)
	}
	targetDir := model.TargetDirectory(*targetDirPtr)
	if len(targetDir) == 0 {
		return fmt.Errorf("parameter %s missing", parameterTargetDir)
	}
	name := model.Name(*namePtr)
	if len(name) == 0 {
		return fmt.Errorf("parameter %s missing", parameterName)
	}

	oneTime := *oneTimePtr
	wait := *waitPtr
	lockName := *lockPtr

	glog.V(1).Infof("name: %s, host: %s, port: %d, user: %s, password-length: %d, database: %s, targetDir: %s, wait: %v, oneTime: %v, lockName: %s", name, host, port, user, len(pass), database, targetDir, wait, oneTime, lockName)

	action := func(ctx context.Context) error {
		return backup.Create(name, host, port, user, pass, database, targetDir)
	}

	cron := cron.New(
		oneTime,
		wait,
		action,
	)
	return cron.Run(context.Background())
}
