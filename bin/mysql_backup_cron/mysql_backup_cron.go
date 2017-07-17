package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/bborbe/mysql_backup_cron/model"
	"github.com/bborbe/mysql_backup_cron/mysql"
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
	mysqlHostPtr     = flag.String(parameterMysqlHost, "", "host")
	mysqlPortPtr     = flag.Int(parameterMysqlPort, 5432, "port")
	mysqlDatabasePtr = flag.String(parameterMysqlDatabase, "", "database")
	mysqlUserPtr     = flag.String(parameterMysqlUser, "", "username")
	mysqlPasswordPtr = flag.String(parameterMysqlPassword, "", "password")
	waitPtr          = flag.Duration(parameterWait, time.Minute*60, "wait")
	oneTimePtr       = flag.Bool(parameterOneTime, false, "exit after first backup")
	targetDirPtr     = flag.String(parameterTargetDir, "", "target directory")
	lockPtr          = flag.String(parameterLock, defaultLockName, "lock")
	namePtr          = flag.String(parameterName, defaultName, "name")
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

	oneTime := *oneTimePtr
	wait := *waitPtr

	mysqlDumper := mysql.NewDumper(
		model.MysqlDatabase(*mysqlDatabasePtr),
		false,
		model.MysqlHost(*mysqlHostPtr),
		model.MysqlPort(*mysqlPortPtr),
		model.MysqlUser(*mysqlUserPtr),
		model.MysqlPassword(*mysqlPasswordPtr),
		model.Name(*namePtr),
		model.TargetDirectory(*targetDirPtr),
	)

	glog.V(1).Infof("name: %s, host: %s, port: %d, user: %s, password-length: %d, database: %s, targetDir: %s, wait: %v, oneTime: %v, lockName: %s", mysqlDumper.Name, mysqlDumper.Host, mysqlDumper.Port, mysqlDumper.User, len(mysqlDumper.Password), mysqlDumper.Database, mysqlDumper.TargetDirectory, wait, oneTime)

	if err := mysqlDumper.Validate(); err != nil {
		return fmt.Errorf("validate mysql parameter failed: %v", err)
	}

	var c cron.Cron
	if oneTime {
		c = cron.NewOneTimeCron(mysqlDumper.Run)
	} else {
		c = cron.NewWaitCron(
			wait,
			mysqlDumper.Run,
		)
	}
	return c.Run(context.Background())
}
