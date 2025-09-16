// Package dailystats provides access to daily statistics on the database.
// The funcs are meant to be called by the server.
// Generally, payloads to be stored are handled as JSON strings
// which are generated and consumable by the server client.
package dailystats

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/scottbrodersen/homegym/dal"
)

// A DailyStats stores statistics that a user tracks every day.
type DailyStats struct {
	Date          int64   `json:"date"`
	BloodGlucose  float32 `json:"bg,omitempty"`
	BloodPressure []int   `json:"bp,omitempty"`
	Sleep         float32 `json:"sleep,omitempty"`
	Food          Food    `json:"food,omitempty"`
	BodyWeight    int     `json:"bodyweight,omitempty"`
	Mood          int     `json:"mood,omitempty"`
	Stress        int     `json:"stress,omitempty"`
	Energy        int     `json:"energy,omitempty"`
}

type Food struct {
	Protein     int    `json:"protein"`
	Carbs       int    `json:"carbs"`
	Fat         int    `json:"fat"`
	Fiber       int    `json:"fiber"`
	Description string `json:"description,omitempty"`
}

// An ErrInvalidStats generates an error for when a daily stat is invalid.
type ErrInvalidStats struct {
	Message string
}

func (e ErrInvalidStats) Error() string {
	return fmt.Sprintf("invalid stats: %s", e.Message)
}

// The DailyStatsAdmin interface defines the utility funcs for daily stats.
type DailyStatsAdmin interface {
	AddStats(userID string, date int64, statsJSON []byte) error
	GetBioStatsPage(userID string, startDate, endDate int64, pageSize int) ([]byte, error)
}

// A DailyStatsUtil implements DailyStatsAdmin.
type DailyStatsUtil struct{}

var DailyStatsManager DailyStatsAdmin = DailyStatsUtil{}

// AddStats stores a set of daily stats for a specific date in the database.
// The provided JSON data is validated before being stored.
func (dsu DailyStatsUtil) AddStats(userID string, date int64, statsJSON []byte) error {

	stats := DailyStats{}

	if err := json.Unmarshal(statsJSON, &stats); err != nil {
		slog.Debug(err.Error())
		return ErrInvalidStats{Message: "could not unmarshal stats"}
	}

	if err := stats.validate(); err != nil {
		slog.Debug(err.Error())
		return ErrInvalidStats{Message: err.Error()}
	}

	if err := dal.DB.AddBioStats(userID, date, statsJSON); err != nil {
		return fmt.Errorf("could not add stats: %w", err)
	}

	return nil
}

// GetBioStatsPage retrieves a page of daily stats from the database.
// The page size is limited, and defaults to, 3000.
func (dsu DailyStatsUtil) GetBioStatsPage(userID string, startDate, endDate int64, pageSize int) ([]byte, error) {
	if pageSize == 0 || pageSize > 3000 {
		pageSize = 3000 // stat instances
	}

	statsByte, err := dal.DB.GetBioStatsPage(userID, startDate, pageSize)
	if err != nil {
		return nil, err
	}

	stats := []DailyStats{}

	// return only the stats that are within the date range
	for _, statByte := range statsByte {
		stat := DailyStats{}
		err := json.Unmarshal(statByte, &stat)
		if err != nil {
			return nil, ErrInvalidStats{Message: fmt.Sprintf("error unmarshaling a daily stat, %s", err.Error())}
		}

		if stat.Date > endDate {
			stats = append(stats, stat)
		}
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return nil, fmt.Errorf("could not marshal daily stats: %w", err)
	}

	return statsJSON, nil
}

func (ds DailyStats) validate() error {
	if ds.Date == 0 {
		return fmt.Errorf("date is a required field")
	}

	if ds.BloodPressure != nil && len(ds.BloodPressure) != 2 {
		return fmt.Errorf("blood pressure slice must contain 2 values")
	}

	if (ds.BloodGlucose == 0 && ds.BloodPressure == nil && ds.Sleep == 0 && ds.Food == Food{} && ds.BodyWeight == 0 && ds.Mood == 0 && ds.Stress == 0 && ds.Energy == 0) {
		return fmt.Errorf("at least one daily stat is required")
	}

	return nil
}
