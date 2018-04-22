package lock

import (
	"os"
	"syscall"

	"github.com/golang/glog"
)

type Lock interface {
	Lock() error
	Unlock() error
}

type lock struct {
	lockName string
	file     *os.File
}

func NewLock(lockName string) *lock {
	l := new(lock)
	l.lockName = lockName
	return l
}

func (l *lock) Lock() error {
	glog.V(2).Info("try lock")
	var err error
	l.file, _ = os.Open(l.lockName)
	if l.file == nil {
		l.file, err = os.Create(l.lockName)
		if err != nil {
			glog.V(2).Info("create lock file failed")
			return err
		}
	}
	err = syscall.Flock(int(l.file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		glog.V(2).Info("lock fail, already locked")
		return err
	}
	glog.V(2).Info("locked")
	return nil
}

func (l *lock) Unlock() error {
	glog.V(2).Info("try unlock")
	var err error
	err = syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	if err != nil {
		glog.V(2).Info("unlock failed")
		return err
	}
	glog.V(2).Info("unlocked")
	return l.file.Close()
}
