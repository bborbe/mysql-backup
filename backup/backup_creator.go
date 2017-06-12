package backup

import (
	"fmt"
	"github.com/bborbe/io/util"
	"github.com/bborbe/mysql_backup_cron/model"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Create backup
func Create(
	name model.Name,
	host model.MysqlHost,
	port model.MysqlPort,
	user model.MysqlUser,
	pass model.MysqlPassword,
	database model.MysqlDatabase,
	targetDirectory model.TargetDirectory,
) error {
	// mysqldump --lock-tables=false -u [USER] -h [HOSTNAME] -p [DATABASENAME]
	backupfile := model.BuildBackupfileName(name, targetDirectory, database, time.Now())
	if backupfile.Exists() {
		glog.V(1).Infof("backup %s already exists => skip", backupfile)
		return nil
	}
	path, err := util.NormalizePath("~/.my.cnf")
	if err != nil {
		return err
	}
	if err := writePasswordFile(path, user, pass); err != nil {
		return err
	}
	glog.V(1).Infof("mysqldump started")
	if err := runCommand("mysqldump", targetDirectory, "--defaults-file="+path, "--lock-tables=false", "--user", user.String(), "--host", host.String(), "--port", port.String(), "--result-file", backupfile.String(), database.String()); err != nil {
		glog.V(2).Infof("mysqldump failed, delete incomplete backup: %v", err)
		if err := backupfile.Delete(); err != nil {
			glog.Warningf("delete incomplete backup failed: %v", err)
		}
		return err
	}
	glog.V(1).Infof("mysqldump finshed")
	return nil
}

func writePasswordFile(path string, user model.MysqlUser, pass model.MysqlPassword) error {
	content := fmt.Sprintf("[mysqldump]\nuser=%s\npassword=%s\n\n[mysql]\nuser=%s\npassword=%s\n", user.String(), pass.String(), user.String(), pass.String())
	return ioutil.WriteFile(path, []byte(content), 0600)
}

func runCommand(command string, cwd model.TargetDirectory, args ...string) error {
	debug := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	glog.V(2).Infof("execute %s", debug)
	cmd := exec.Command(command, args...)
	if glog.V(4) {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	if cwd != "" {
		cmd.Dir = cwd.String()
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	glog.V(2).Infof("%s started", debug)
	if err := cmd.Wait(); err != nil {
		glog.Warningf("%s failed: %v", debug, err)
		return fmt.Errorf("%s failed: %v", debug, err)
	}
	glog.V(2).Infof("%s finished", command)
	return nil
}
