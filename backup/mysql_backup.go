package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/bborbe/mysql_backup_cron/dumper"
	"github.com/bborbe/mysql_backup_cron/model"
	"github.com/golang/glog"
)

type Backuper struct {
	Database        model.MysqlDatabase
	AllDatabases    bool
	Host            model.MysqlHost
	Port            model.MysqlPort
	User            model.MysqlUser
	Password        model.MysqlPassword
	Name            model.Name
	TargetDirectory model.TargetDirectory
	DateFilename    bool
	OverwriteBackup bool
}

func New(
	database model.MysqlDatabase,
	allDatabases bool,
	host model.MysqlHost,
	port model.MysqlPort,
	user model.MysqlUser,
	password model.MysqlPassword,
	name model.Name,
	targetDirectory model.TargetDirectory,
	dateFilename bool,
	overwriteBackup bool,
) *Backuper {
	d := new(Backuper)
	d.Database = database
	d.Host = host
	d.Port = port
	d.User = user
	d.Password = password
	d.Name = name
	d.TargetDirectory = targetDirectory
	d.AllDatabases = allDatabases
	d.DateFilename = dateFilename
	d.OverwriteBackup = overwriteBackup
	return d
}

func (b *Backuper) Validate() error {
	if len(b.Host) == 0 {
		return fmt.Errorf("mysql host missing")
	}
	if b.Port <= 0 {
		return fmt.Errorf("mysql port missing")
	}
	if len(b.User) == 0 {
		return fmt.Errorf("mysql user missing")
	}
	if len(b.Password) == 0 {
		return fmt.Errorf("mysql password missing")
	}
	if len(b.Name) == 0 {
		return fmt.Errorf("mysql name missing")
	}
	if len(b.TargetDirectory) == 0 {
		return fmt.Errorf("mysql target dir missing")
	}
	if b.AllDatabases == false && len(b.Database) == 0 {
		return fmt.Errorf("mysql database missing")
	}
	if b.AllDatabases == true && len(b.Database) > 0 {
		return fmt.Errorf("mysql database is not allowed with all enabled")
	}
	return nil
}

func (b *Backuper) checkBackupShouldBeSkipped(backupFile model.BackupFilename) bool {
	return !b.OverwriteBackup && backupFile.Exists()
}

func (b *Backuper) Run(ctx context.Context) error {
	if err := b.TargetDirectory.Mkdir(0700); err != nil {
		return fmt.Errorf("create targetdirectory %v failed: %v", b.TargetDirectory, err)
	}
	backupFile := b.backupFile()
	if b.checkBackupShouldBeSkipped(backupFile) {
		glog.V(1).Infof("backup %s already exists => skip", backupFile)
		return nil
	}
	dumper := dumper.New(
		b.Name,
		b.Host,
		b.Port,
		b.User,
		b.Password,
		b.TargetDirectory,
	)
	if b.AllDatabases {
		return dumper.All(backupFile)
	}
	return dumper.Database(b.Database, backupFile)
}

func (b *Backuper) backupFile() model.BackupFilename {
	var database model.MysqlDatabase
	if b.AllDatabases {
		database = "all"
	} else {
		database = b.Database
	}
	if b.DateFilename {
		return model.BuildBackupfileNameWithDate(b.Name, b.TargetDirectory, database, time.Now())
	}
	return model.BuildBackupfileName(b.Name, b.TargetDirectory, database)
}
