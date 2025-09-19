/*
homegym initializes a database client and starts the Home Gym server.
The program runs in test mode or production mode according to a flag.
When in test mode, the path to the database is provided as a flag.
In production mode, the path is obtained from the HOMEGYM_DB_PATH environment variable.

Usage:

	homegym -testmode true|false [-path path]

The flags are:

	-testmode
		Set to true to run in test mode and false otherwise.
		When true, you must provide a path.
		When false, the path must be stored in the environment variable.
	-path
		For testmode, the relative path to the database.
*/
package main

import (
	"flag"
	"log/slog"
	"os"

	"path/filepath"

	"log"

	"github.com/scottbrodersen/homegym/auth"
	"github.com/scottbrodersen/homegym/dal"
	"github.com/scottbrodersen/homegym/server"
)

const (
	dbPathEnv  = "HOMEGYM_DB_PATH"
	testDBPath = "./../server/testDB"
)

func main() {
	var currentLogLevel = slog.SetLogLoggerLevel(slog.LevelInfo)
	defer slog.SetLogLoggerLevel(currentLogLevel)

	pathFlag := flag.String("path", "", "path to the homegym database")
	testModeFlag := flag.Bool("testmode", false, "run in test mode")

	flag.Parse()

	dbPath := ""
	port := 0

	switch *testModeFlag {
	case false:
		dbPath = *pathFlag
		if dbPath == "" {
			dbPath = os.Getenv(dbPathEnv)
		}
		port = 80
	case true:
		dbPath = testDBPath
		port = 3000
	}

	if dbPath == "" {
		log.Fatal("Database path not configured")
	}
	slog.Debug("using database", "path", dbPath)

	db, err := dal.InitClient(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Destroy()
	dal.DB = db
	if err = auth.InitiateKeyRotation(auth.KeyTypes.Token); err != nil {
		log.Fatal(err)
	}
	auth.CleanupSessions()

	dal.InitHourlyGC(*db)

	dbBackupDir := filepath.Join(dbPath, "backups")
	if err := os.Mkdir(dbBackupDir, 0750); err != nil && !os.IsExist(err) {
		// backups are important
		log.Fatal(err)
	}

	dbBackupFile := filepath.Join(dbBackupDir, "backup.bak")
	dal.InitDailyBackup(*db, dbBackupFile)

	if dbPath == testDBPath {
		if err := AddData(); err != nil {
			slog.Warn("error adding data.", "error", err.Error())
			dal.DB.Iter8er()
		}
	}

	// StartUnsafe for dev purposes only!
	//server.StartUnsafe(server.DefaultShutdown, port)
	server.StartSafe(server.DefaultShutdown, port)
}
