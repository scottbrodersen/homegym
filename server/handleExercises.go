package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/scottbrodersen/homegym/workoutlog"
)

type returnedID = struct {
	ID string `json:"id"`
}

type returnedValue = struct {
	Value int `json:"value"`
}

// Handles requests for exercise types.
func ExerciseTypesApi(w http.ResponseWriter, r *http.Request) {
	rootpath := "/homegym/api/exercises/"

	username, _, err := whoIsIt(r.Context())
	if err != nil {
		slog.Debug(err.Error())
		http.Error(w, fmt.Sprintf("{\"message\": \"%s\"}", err.Error()), http.StatusForbidden)
		return
	}

	rxpNewType := regexp.MustCompile(fmt.Sprintf("^%s$", rootpath))
	rxpUpdateType := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]+)/?$", rootpath))
	rxpTypeOneRM := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]+)/onerm/?$", rootpath))
	rxpTypePR := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]+)/pr/?$", rootpath))

	if rxpNewType.MatchString(r.URL.Path) {
		if r.Method == http.MethodPost {
			newExerciseType(*username, w, r)
			return
		} else if r.Method == http.MethodGet {
			listExerciseTypes(*username, w)
			return
		}
	} else if rxpUpdateType.MatchString(r.URL.Path) {
		if r.Method == http.MethodPost {
			typeID := rxpUpdateType.FindStringSubmatch(r.URL.Path)[1]

			updateExerciseType(*username, typeID, w, r)
			return
		}
	} else if rxpTypeOneRM.MatchString(r.URL.Path) {
		typeID := rxpTypeOneRM.FindStringSubmatch(r.URL.Path)[1]

		if r.Method == http.MethodPost {
			newOneRM(*username, typeID, w, r)
			return
		} else if r.Method == http.MethodGet {
			getOneRM(*username, typeID, w)
			return
		}
	} else if rxpTypePR.MatchString(r.URL.Path) {
		typeID := rxpTypePR.FindStringSubmatch(r.URL.Path)[1]

		if r.Method == http.MethodPost {
			newPR(*username, typeID, w, r)
			return
		} else if r.Method == http.MethodGet {
			getPR(*username, typeID, w)
			return
		}
	}

	http.Error(w, "", http.StatusNotFound)
}

func newExerciseType(username string, w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		slog.Debug("No request body")
		http.Error(w, `{"message": "request body is required"}`, http.StatusBadRequest)
		return
	}
	var et *workoutlog.ExerciseType = &workoutlog.ExerciseType{}
	if err := json.NewDecoder(r.Body).Decode(et); err != nil {
		slog.Debug(err.Error())
		http.Error(w, `{"message": "problem with request body"}`, http.StatusBadRequest)
		return
	}

	id, err := workoutlog.ExerciseManager.NewExerciseType(username, et.Name, et.IntensityType, et.VolumeType, et.VolumeConstraint, et.Composition, et.Basis)

	if err != nil {
		nu := workoutlog.ErrNameNotUnique
		if errors.Is(err, nu) {
			slog.Debug(err.Error())
			http.Error(w, "name is not unique", http.StatusBadRequest)
		} else if errors.As(err, &workoutlog.ErrInvalidExercise{}) {
			slog.Debug(err.Error())
			http.Error(w, fmt.Sprintf("{\"message\": \"%s\"}", err.Error()), http.StatusBadRequest)
		} else {
			slog.Error(err.Error())
			http.Error(w, internalServerError, http.StatusInternalServerError)
		}
		return
	}

	body := returnedID{ID: *id}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(bodyJson)
}

func updateExerciseType(username, typeID string, w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		slog.Debug("No request body")
		http.Error(w, `"message": "request body is required"`, http.StatusBadRequest)
		return
	}
	var updated *workoutlog.ExerciseType = &workoutlog.ExerciseType{}
	if err := json.NewDecoder(r.Body).Decode(updated); err != nil {
		slog.Debug(err.Error())
		http.Error(w, `"message": "problem with request body"`, http.StatusBadRequest)
		return
	}
	if updated.ID != typeID {
		slog.Debug("typeID in path does not match typeID in body")
		http.Error(w, `"message": "typeID mismatch"`, http.StatusBadRequest)
		return
	}

	err := workoutlog.ExerciseManager.UpdateExerciseType(username, typeID, updated.Name, updated.IntensityType, updated.VolumeType, updated.VolumeConstraint, updated.Composition, updated.Basis)

	if err != nil {
		nu := workoutlog.ErrNameNotUnique
		if errors.Is(err, nu) {
			slog.Debug(err.Error())
			http.Error(w, `"message": "name is not unique"`, http.StatusBadRequest)
		} else if errors.As(err, &workoutlog.ErrInvalidExercise{}) {
			http.Error(w, fmt.Sprintf("{\"message\": \"%s\"}", err.Error()), http.StatusBadRequest)
		} else {
			slog.Error(err.Error())
			http.Error(w, internalServerError, http.StatusInternalServerError)
		}
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusNoContent)
}

func listExerciseTypes(username string, w http.ResponseWriter) {
	types, err := workoutlog.ExerciseManager.GetExerciseTypes(username)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	body, err := json.Marshal(types)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, `{"message": "failed to get exercise types"}`, http.StatusInternalServerError)
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func newPR(username, exID string, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Debug("Could not parse URL parameters")
		http.Error(w, `{"message": "could not parse URL query parameters"}`, http.StatusBadRequest)
		return
	}

	prValue := r.Form.Get("pr")
	if prValue == "" {
		slog.Debug("No pr parameter")
		http.Error(w, `{"message": "no pr parameter"}`, http.StatusBadRequest)
		return
	}

	prInt, err := strconv.Atoi(prValue)
	if err != nil {
		slog.Debug("cannot convert PR value to int")
		http.Error(w, `{"message": "bad pr value"}`, http.StatusBadRequest)
		return
	}

	err = workoutlog.ExerciseManager.SetPR(username, exID, prInt)

	if err != nil {
		slog.Debug("error setting PR")
		http.Error(w, `{"message": "error setting PR"}`, http.StatusBadRequest)
	}
	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusOK)
}

func getPR(username, exID string, w http.ResponseWriter) {
	pr, err := workoutlog.ExerciseManager.GetPR(username, exID)

	if err != nil {
		slog.Debug("error setting PR")
		http.Error(w, `{"message": "error setting PR"}`, http.StatusBadRequest)
		return
	}
	if pr == -1 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	body := returnedValue{Value: pr}
	bodyJSON, err := json.Marshal(body)

	if err != nil {
		slog.Debug("error marshalling response body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(bodyJSON)
}

func newOneRM(username, exID string, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Debug("Could not parse URL parameters")
		http.Error(w, `{"message": "could not parse URL query parameters"}`, http.StatusBadRequest)
		return
	}

	oneRMValue := r.Form.Get("onerm")
	if oneRMValue == "" {
		slog.Debug("No onerm parameter")
		http.Error(w, `{"message": "no onerm parameter"}`, http.StatusBadRequest)
		return
	}

	oneRMInt, err := strconv.Atoi(oneRMValue)
	if err != nil {
		slog.Debug("cannot convert onerm value to int")
		http.Error(w, `{"message": "bad onerm value"}`, http.StatusBadRequest)
	}

	err = workoutlog.ExerciseManager.Set1RM(username, exID, oneRMInt)

	if err != nil {
		slog.Debug("error setting 1RM")
		http.Error(w, `{"message": "error setting 1RM"}`, http.StatusBadRequest)
	}
	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusOK)
}

func getOneRM(username, exID string, w http.ResponseWriter) {
	oneRM, err := workoutlog.ExerciseManager.Get1RM(username, exID)

	if err != nil {
		slog.Debug("error setting 1RM")
		http.Error(w, `{"message": "error setting 1RM"}`, http.StatusBadRequest)
	}
	if oneRM == -1 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	body := returnedValue{Value: oneRM}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		slog.Debug("error marshalling response body")
		http.Error(w, "", http.StatusInternalServerError)
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(bodyJSON)
}
