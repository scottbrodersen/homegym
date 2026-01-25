package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/scottbrodersen/homegym/admin"
	"github.com/scottbrodersen/homegym/auth"
)

// AdminAPI handles requests for admin tasks.
func AdminApi(w http.ResponseWriter, r *http.Request) {
	rootpath := "/homegym/api/admin/"

	_, role, err := whoIsIt(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// allow admin roles only
	if *role != string(auth.Admin) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// path to restore daily backup
	rxpRestoreDaily := regexp.MustCompile(fmt.Sprintf("^%srestoredaily/?$", rootpath))

	if rxpRestoreDaily.MatchString(r.URL.Path) {
		if r.Method == http.MethodGet {
			restoreDailyBackup(w)
			return
		}
	}
	http.Error(w, "not found admin", http.StatusNotFound)

}

func restoreDailyBackup(w http.ResponseWriter) {
	backupPath := filepath.Join(admin.DBPath, admin.DBBackupDir, admin.DBDailyBackupFile)

	if err := admin.DatabaseManager.RestoreBackup(backupPath); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}
