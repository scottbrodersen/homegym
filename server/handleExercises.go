package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/scottbrodersen/homegym/workoutlog"
	log "github.com/sirupsen/logrus"
)

type returnedID = struct {
	ID string `json:"id"`
}

func ExerciseTypesApi(w http.ResponseWriter, r *http.Request) {
	rootpath := "/homegym/api/exercises/"

	username, _, err := whoIsIt(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	rxpNewType := regexp.MustCompile(fmt.Sprintf("^%s$", rootpath))
	rxpUpdateType := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]+)/?$", rootpath))

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
	}

	http.Error(w, "", http.StatusNotFound)
}

func newExerciseType(username string, w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, "request body is required", http.StatusBadRequest)
		return
	}
	var et *workoutlog.ExerciseType = &workoutlog.ExerciseType{}
	if err := json.NewDecoder(r.Body).Decode(et); err != nil {
		log.Error(err)
		http.Error(w, "problem with request body", http.StatusBadRequest)
		return
	}

	id, err := workoutlog.ExerciseManager.NewExerciseType(username, et.Name, et.IntensityType, et.VolumeType, et.VolumeConstraint, et.Composition, et.Basis)

	if err != nil {
		nu := workoutlog.ErrNameNotUnique
		if errors.Is(err, nu) {
			http.Error(w, "name is not unique", http.StatusBadRequest)
		} else if errors.As(err, &workoutlog.ErrInvalidExercise{}) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, internalServerError, http.StatusInternalServerError)
		}
		return
	}

	body := returnedID{ID: *id}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(bodyJson)
}

func updateExerciseType(username, typeID string, w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, `"message": "request body is required"`, http.StatusBadRequest)
		return
	}
	var updated *workoutlog.ExerciseType = &workoutlog.ExerciseType{}
	if err := json.NewDecoder(r.Body).Decode(updated); err != nil {
		log.Error(err)
		http.Error(w, `"message": "problem with request body"`, http.StatusBadRequest)
		return
	}
	if updated.ID != typeID {
		http.Error(w, `"message": "typeID mismatch"`, http.StatusBadRequest)
		return
	}

	err := workoutlog.ExerciseManager.UpdateExerciseType(username, typeID, updated.Name, updated.IntensityType, updated.VolumeType, updated.VolumeConstraint, updated.Composition, updated.Basis)

	if err != nil {
		nu := workoutlog.ErrNameNotUnique
		if errors.Is(err, nu) {
			http.Error(w, `"message": "name is not unique"`, http.StatusBadRequest)
		} else if errors.As(err, &workoutlog.ErrInvalidExercise{}) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
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
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	body, err := json.Marshal(types)
	if err != nil {
		message := struct{ Message string }{Message: "failed to get exercise types"}
		body, _ = json.Marshal(message)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}
