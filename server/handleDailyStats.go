package server

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/scottbrodersen/homegym/dailystats"
)

func DailyStatsApi(w http.ResponseWriter, r *http.Request) {
	rootpath := "/homegym/api/dailystats/"
	username, _, err := whoIsIt(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"message\": \"%s\"}", err.Error()), http.StatusForbidden)
		return
	}

	rxpStats := regexp.MustCompile(rootpath)

	if rxpStats.MatchString(r.URL.Path) {
		if r.Method == http.MethodPost {
			addDailyStats(*username, w, r)
			return
		} else if r.Method == http.MethodGet {
			getDailyStats(*username, w, r)
			return
		}
	}

	http.Error(w, "", http.StatusNotFound)

}

func addDailyStats(username string, w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, `{"message": "request body is required"}`, http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, `{"message": "could not parse URL query parameters"}`, http.StatusBadRequest)
		return
	}

	dateStr := r.Form.Get("date")
	if dateStr == "" {
		http.Error(w, `{"message": "no date parameter"}`, http.StatusBadRequest)
		return
	}

	date, err := strconv.Atoi(dateStr)
	if err != nil {
		http.Error(w, `{"message": "bad date parameter"}`, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"message": "error reading request body"}`, http.StatusBadRequest)
		return
	}

	if err := dailystats.DailyStatsManager.AddStats(username, int64(date), body); err != nil {
		slog.Error(err.Error())
		if errors.As(err, new(dailystats.ErrInvalidStats)) {
			http.Error(w, fmt.Sprintf("{\"message\": \"%s\"}", err.Error()), http.StatusBadRequest)
			return
		}
		http.Error(w, `{"message": "error adding stats"}`, http.StatusInternalServerError)
		return
	}
	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusOK)
}

func getDailyStats(username string, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, `{"message": "could not parse URL query parameters"}`, http.StatusBadRequest)
		return
	}

	startDate := 0
	start := r.Form.Get("start")
	if start != "" {
		startDate, err = strconv.Atoi(start)
		if err != nil {
			http.Error(w, `{"message": "bad start value"}`, http.StatusBadRequest)
			return
		}
	}
	endDate := 0
	end := r.Form.Get("end")
	if end != "" {
		endDate, err = strconv.Atoi(end)
		if err != nil {
			http.Error(w, `{"message": "bad end value"}`, http.StatusBadRequest)
			return
		}
	}

	page := 0
	pageSize := r.Form.Get("pagesize")
	if pageSize != "" {
		page, err = stringToInt(pageSize)
		if err != nil {
			http.Error(w, `{"message": "bad pagesize value"}`, http.StatusBadRequest)
			return
		}
	}

	stats, err := dailystats.DailyStatsManager.GetBioStatsPage(username, int64(startDate), int64(endDate), page)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, `{"message": "error getting stats"}`, http.StatusBadRequest)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(stats)
}
