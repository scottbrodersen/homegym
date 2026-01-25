package admin

import (
	"log/slog"

	"github.com/scottbrodersen/homegym/dal"
)

var DBPath string
var DBBackupDir = "backups"
var DBDailyBackupFile = "backup.bak"

// DatabaseAdmin defines functions for performing database admin tasks.
type DatabaseAdmin interface {
	RestoreBackup(userID, filepath string) error
}

// DatabaseManager implements DatabaseAdmin.
type DatabaseUtil struct{}

var DatabaseManager DatabaseUtil = DatabaseUtil{}

// RestoreBackup restores the database using a backup file.
func (*DatabaseUtil) RestoreBackup(filepath string) error {
	err := dal.DB.Restore(filepath)
	if err != nil {
		slog.Error("error restoring backup", "error", err.Error())
		return err
	}

	slog.Info("database restored")
	return nil
}
