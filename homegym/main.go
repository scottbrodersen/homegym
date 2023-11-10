package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/scottbrodersen/homegym/auth"
	"github.com/scottbrodersen/homegym/server"

	"github.com/scottbrodersen/homegym/dal"
)

const (
	dbPathEnv  = "HOMEGYM_DB_PATH"
	testDBPath = "./../server/testDB"
)

func main() {

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
	log.Debug("using database at path ", dbPath)
	var db dal.Dstore

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

	if dbPath == testDBPath {
		if err := AddData(); err != nil {
			log.WithError(err).Warn("error adding data")
		}
	}

	server.StartUnsafe(server.DefaultShutdown, port)

}
