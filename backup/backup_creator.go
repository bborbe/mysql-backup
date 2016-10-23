package backup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"io/ioutil"

	"github.com/bborbe/io/util"
	"github.com/bborbe/postgres_backup_cron/model"
	"github.com/golang/glog"
)

// Create backup
func Create(
	host model.PostgresqlHost,
	port model.PostgresqlPort,
	user model.PostgresqlUser,
	pass model.PostgresqlPassword,
	database model.PostgresqlDatabase,
	targetDirectory model.TargetDirectory,
) error {
	//pg_dump -Z 9 -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER} -F c -b -v -f ${BACKUP_NAME} ${POSTGRES_DB}
	backupfile := model.BuildBackupfileName(targetDirectory, database, time.Now())

	if backupfile.Exists() {
		glog.V(1).Infof("backup %s already exists => skip", backupfile)
		return nil
	}

	if err := writePasswordFile(host, port, user, pass); err != nil {
		return err
	}

	glog.V(1).Infof("pg_dump started")
	if err := runCommand("pg_dump", targetDirectory, "-Z", "9", "-h", host.String(), "-p", port.String(), "-U", user.String(), "-F", "c", "-b", "-v", "-f", backupfile.String(), database.String()); err != nil {
		glog.V(2).Infof("pg_dump failed, delete incomplete backup: %v", err)
		if err := backupfile.Delete(); err != nil {
			glog.Warningf("delete incomplete backup failed: %v", err)
		}
		return err
	}
	glog.V(1).Infof("pg_dump finshed")
	return nil
}

func writePasswordFile(host model.PostgresqlHost, port model.PostgresqlPort, user model.PostgresqlUser, pass model.PostgresqlPassword) error {
	content := fmt.Sprintf("%s:%d:*:%s:%s\n", host, port, user, pass)
	path, err := util.NormalizePath("~/.pgpass")
	if err != nil {
		return err
	}
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
