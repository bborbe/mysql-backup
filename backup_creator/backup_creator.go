package backup_creator

type backupCreator struct {
}

type BackupCreator interface {
	CreateBackup(host string, port int, user string, pass string, database string) error
}

func New() *backupCreator {
	return new(backupCreator)
}

func (b *backupCreator) CreateBackup(host string, port int, user string, pass string, database string) error {
	return nil
}
