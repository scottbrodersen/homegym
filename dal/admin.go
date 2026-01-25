package dal

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

func collectGarbage(db *badger.DB) error {
	return db.RunValueLogGC(0.5)
}

// InitHourlyGC schedules hourly garbage collection on the database.
func InitHourlyGC(db DBClient) {
	_, m, s := time.Now().Clock()
	secondsUntil := 3600 - m*60 - s

	slog.Info("Scheduling DB garbage collection in", fmt.Sprint(secondsUntil), "seconds")

	time.AfterFunc(time.Second*time.Duration(secondsUntil), func() {
		if err := collectGarbage(db.db); err != nil {
			slog.Info(err.Error())
		}

		slog.Info("GC complete")

		InitHourlyGC(db)
	})
}

func backup(db *badger.DB, filePath string) error {
	// Save a monthly backup on the first day of the month
	now := time.Now()
	if now.Day() == 1 {
		filePath = fmt.Sprintf("%s_%s", now.Month().String(), filePath)
	}

	// back up the existing backup
	if err := os.Rename(filePath, fmt.Sprintf("%s.old", filePath)); err != nil {
		slog.Error(err.Error())
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = db.Backup(file, 0)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

// InitDailyBackup schedules a daily backup of the database.
func InitDailyBackup(db DBClient, filePath string) {
	h, m, s := time.Now().Clock()
	secondsUntil := 24*60*60 - h*60*60 - m*60 - s

	slog.Info("Scheduling DB backup in", fmt.Sprint(secondsUntil), "seconds")

	time.AfterFunc(time.Second*time.Duration(secondsUntil), func() {
		if err := backup(db.db, filePath); err != nil {
			slog.Info(err.Error())
		}

		slog.Info("Backup complete")

		InitDailyBackup(db, filePath)
	})
}

func (c *DBClient) Restore(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		slog.Info(err.Error())
		return err
	}
	r := bufio.NewReader(f)
	err = c.db.Load(r, 100)
	if err != nil {
		slog.Info(err.Error())
		return err
	}
	return nil
}
