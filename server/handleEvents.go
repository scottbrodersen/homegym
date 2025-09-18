package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/scottbrodersen/homegym/workoutlog"
)

type metrics struct {
	Dates  []int64   `json:"dates"`
	Volume []float32 `json:"volume"`
	Load   []float32 `json:"load"`
}

// EventsAi handles requests for events.
func EventsApi(w http.ResponseWriter, r *http.Request) {
	rootPath := "/homegym/api/events/"

	username, _, err := whoIsIt(r.Context())
	if err != nil {
		slog.Debug(err.Error())
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// path to all events
	rxpRootPath := regexp.MustCompile(fmt.Sprintf("^%s$", rootPath))
	//  path to an event
	rxpEventPath := regexp.MustCompile(fmt.Sprintf("^%s(\\d+)/([a-zA-Z0-9-]+)/?$", rootPath))
	//  path to the exercise instances of an event
	rxpExercisesPath := regexp.MustCompile(fmt.Sprintf("^%s(\\d+)/([a-zA-Z0-9-]+)/exercises/?$", rootPath))
	//  path to metrics for a range of events e.g. /api/events/metrics?type=blah&start=blah&end=blah
	rxpMetrics := regexp.MustCompile(fmt.Sprintf("^%smetrics(\\?[a-z]+=[a-zA-Z0-9-]+((&[a-z]+=[a-zA-Z0-9-]+)*)?)?$", rootPath))

	slog.Debug("parsing path", "path", r.URL.Path)

	if rxpRootPath.MatchString(r.URL.Path) {
		if r.Method == http.MethodGet {
			getPageOfEvents(*username, w, r)
			body, _ := io.ReadAll(r.Body)
			r.Body.Close()
			slog.Debug("got page of events: ", "events", string(body))
		} else if r.Method == http.MethodPost {
			addEvent(*username, w, r)
		}
		return
	} else if rxpExercisesPath.MatchString(r.URL.Path) {
		eventID := rxpExercisesPath.FindStringSubmatch(r.URL.Path)[2]

		if r.Method == http.MethodGet {
			getExercises(*username, eventID, w)
			return
		}
	} else if rxpEventPath.MatchString(r.URL.Path) {
		currentDate := rxpEventPath.FindStringSubmatch(r.URL.Path)[1]
		eventID := rxpEventPath.FindStringSubmatch(r.URL.Path)[2]
		if r.Method == http.MethodPost {
			updateEvent(*username, currentDate, w, r)
			return
		} else if r.Method == http.MethodDelete {
			deleteEvent(*username, eventID, currentDate, w, r)
			return
		}
	} else if rxpMetrics.MatchString(r.URL.Path) {
		if r.Method == http.MethodGet {

			getMetrics(*username, w, r)
			return
		}
	}

	http.Error(w, `{"message": "unsupported request type"}`, http.StatusBadRequest)
}

func addEvent(username string, w http.ResponseWriter, r *http.Request) {
	newEvent := new(workoutlog.Event)

	if err := json.NewDecoder(r.Body).Decode(newEvent); err != nil {
		slog.Debug(err.Error())
		http.Error(w, `{"message": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	eventID, err := workoutlog.EventManager.NewEvent(username, *newEvent)

	if err != nil {
		if errors.Is(err, workoutlog.ErrInvalidEvent) {
			slog.Debug(err.Error())
			http.Error(w, `{"message":"invalid event"}`, http.StatusBadRequest)
			return
		}
		slog.Error(err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	body := returnedID{ID: *eventID}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(bodyJson)
	w.WriteHeader(http.StatusOK)
}

func updateEvent(username, currentDate string, w http.ResponseWriter, r *http.Request) {
	updatedEvent := new(workoutlog.Event)

	currentDateInt, err := stringToInt64(currentDate)
	if err != nil {
		slog.Debug(err.Error())
		http.Error(w, "invalid date format", http.StatusBadRequest)
	}

	if err := json.NewDecoder(r.Body).Decode(updatedEvent); err != nil {
		slog.Debug(err.Error())
		http.Error(w, `{"message": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := workoutlog.EventManager.UpdateEvent(username, currentDateInt, *updatedEvent); err != nil {
		if errors.Is(err, workoutlog.ErrNotFound) {
			slog.Debug(err.Error())
			http.Error(w, `{"message":"event not found"}`, http.StatusNotFound)
			return
		} else if errors.Is(err, workoutlog.ErrInvalidEvent) {
			slog.Debug(err.Error())
			http.Error(w, `{"message":"invalid event"}`, http.StatusBadRequest)
			return
		} else {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusNoContent)
}

func deleteEvent(username, eventID, eventDate string, w http.ResponseWriter, r *http.Request) {
	event := new(workoutlog.Event)

	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		slog.Debug(err.Error())
		http.Error(w, `{"message": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	dateInt, err := stringToInt64(eventDate)
	if err != nil {
		slog.Debug(err.Error())
		http.Error(w, `{"message": "invalid data parameter"}`, http.StatusBadRequest)
		return
	}

	if dateInt != event.Date {
		slog.Debug("incorrect date")
		http.Error(w, `{"message": "incorrect date"}`, http.StatusBadRequest)
		return
	}

	if eventID != event.ID {
		slog.Debug("incorrect ID")
		http.Error(w, `{"message": "incorrect ID"}`, http.StatusBadRequest)
		return
	}

	if err := workoutlog.EventManager.DeleteEvent(username, *event); err != nil {
		slog.Debug(err.Error())
		http.Error(w, `{"message":"failed to delete event"}`, http.StatusNotFound)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusNoContent)
}

func getExercises(username, eventID string, w http.ResponseWriter) {
	exercises, err := workoutlog.EventManager.GetEventExercises(username, eventID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	exercisesJson, err := json.Marshal(exercises)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(exercisesJson)
}

func getPageOfEvents(username string, w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		http.Error(w, "failed to read url query parameters", http.StatusBadRequest)
		return
	}

	pageSize, err := stringToInt(r.Form.Get("count"))
	if err != nil {
		slog.Debug(err.Error())
		http.Error(w, "failed to read count param", http.StatusBadRequest)
		return
	}

	date, err := stringToInt64(r.Form.Get("date"))
	if err != nil {
		slog.Debug(err.Error())
		http.Error(w, "failed to read date param", http.StatusBadRequest)
		return
	}

	event := workoutlog.Event{
		ID:   r.Form.Get("previousID"),
		Date: date,
	}

	events, err := workoutlog.EventManager.GetPageOfEvents(username, event, pageSize)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	eventsJson, err := json.Marshal(events)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(eventsJson)
}

// getMetrics returns a page of metrics.
func getMetrics(username string, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Debug(err.Error())
		http.Error(w, `{"message": "could not parse URL query parameters"}`, http.StatusBadRequest)
		return
	}
	startDate := 0
	start := r.Form.Get("start")
	if start != "" {
		startDate, err = strconv.Atoi(start)
		if err != nil {
			slog.Debug(err.Error())
			http.Error(w, `{"message": "bad start value"}`, http.StatusBadRequest)
			return
		}
	}
	endDate := 0
	end := r.Form.Get("end")
	if end != "" {
		endDate, err = strconv.Atoi(end)
		if err != nil {
			slog.Debug(err.Error())
			http.Error(w, `{"message": "bad end value"}`, http.StatusBadRequest)
			return
		}
	}

	filter := workoutlog.ExerciseFilter{
		StartDate:     int64(startDate),
		EndDate:       int64(endDate),
		ExerciseTypes: []string(r.Form["type"]),
	}

	dateStack, instancesStack, err := workoutlog.EventManager.GetPageOfInstances(username, filter, 0)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, fmt.Sprintf("{\"message\": \"%s\"}", err.Error()), http.StatusInternalServerError)
		return
	}

	// Calculate the metrics
	// For each date, calculate the total volume and load
	totalVol := []float32{}
	totalLoad := []float32{}

	for _, instances := range instancesStack {
		volume := float32(0)
		load := float32(0)
		for _, inst := range instances {
			exerciseType, err := workoutlog.ExerciseManager.GetExerciseType(username, inst.TypeID)
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, `{"message": "could not find exercise type"}`, http.StatusInternalServerError)
				return
			}
			instLoad, instVolume := exerciseType.CalculateMetrics(&inst)
			volume += instVolume
			load += instLoad
		}
		totalVol = append(totalVol, volume)
		totalLoad = append(totalLoad, load)
	}

	dateMetrics := metrics{
		Dates:  dateStack,
		Volume: totalVol,
		Load:   totalLoad,
	}

	body, err := json.Marshal(dateMetrics)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, `{"message": "failed to marshal response body"}`, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)

}

// TODO: centralize these helper functions
func stringToInt64(str string) (int64, error) {
	if str == "" {
		return int64(0), nil
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64(i), nil
}

func stringToInt(str string) (int, error) {
	if str == "" {
		return int(0), nil
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}
