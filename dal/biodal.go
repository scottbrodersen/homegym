package dal

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

const (
	bioKey = "bio"
)

// AddBioStats stores daily statistics for a user.
func (c *DBClient) AddBioStats(userID string, date int64, stats []byte) error {
	prefix := []string{userKey, userID, bioKey, fmt.Sprint(date)}

	bgEntry := badger.NewEntry(key(prefix), stats)

	if err := writeUpdates(c, []*badger.Entry{bgEntry}); err != nil {
		return fmt.Errorf("failed to update bio stats: %w", err)
	}

	return nil
}

// GetBioStatsPage gets a page of daily statistics for a user.
// To get stats for a specific day set pageSize to 1.
// The date is used to determine the first item to include.
func (c *DBClient) GetBioStatsPage(userID string, startDate int64, pageSize int) ([][]byte, error) {
	prefix := []string{userKey, userID, bioKey}
	var startKey []byte = nil
	firstPage := true
	if startDate != 0 {
		startKey = key(append(prefix, fmt.Sprint(startDate)))
		if pageSize > 1 {
			firstPage = false
		}
	}
	entries, err := readKeyPrefixPage(c, startKey, key(prefix), pageSize, "", firstPage, true)
	if err != nil {
		return nil, fmt.Errorf("failed to read bio stats: %w", err)
	}

	stats := [][]byte{}

	for _, entry := range entries {

		stats = append(stats, entry.Value)
	}

	return stats, nil

}
