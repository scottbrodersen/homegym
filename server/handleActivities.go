package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/scottbrodersen/homegym/workoutlog"
)

func ActivitiesApi(w http.ResponseWriter, r *http.Request) {
	rootpath := "/homegym/api/activities/"

	username, _, err := whoIsIt(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	rxpRootPath := regexp.MustCompile(fmt.Sprintf("^%s$", rootpath))
	rxpExercises := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]*)/exercises/?$", rootpath))
	rxpActivity := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]*)/?$", rootpath))

	if rxpRootPath.MatchString(r.URL.Path) {
		if r.Method == http.MethodPost {
			newActivity(*username, w, r)
			return
		} else if r.Method == http.MethodGet {
			listActivities(*username, w)
			return
		}
	} else if rxpExercises.MatchString(r.URL.Path) {
		activityID := rxpExercises.FindStringSubmatch(r.URL.Path)[1]

		if r.Method == http.MethodGet {
			listExercises(*username, activityID, w)
			return
		} else if r.Method == http.MethodPost {
			updateExercises(*username, activityID, w, r)
			return
		}
	} else if rxpActivity.MatchString(r.URL.Path) {
		activityID := rxpActivity.FindStringSubmatch(r.URL.Path)[1]
		if r.Method == http.MethodPost {
			updateActivity(*username, activityID, w, r)
			return
		}
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func listActivities(username string, w http.ResponseWriter) {
	activities, err := workoutlog.ActivityManager.GetActivityNames(username)
	if err != nil {
		http.Error(w, `{"message":"failed to get activities"}`, http.StatusInternalServerError)
	}

	body, err := json.Marshal(activities)
	if err != nil {
		http.Error(w, `{"message":"failed to get activities"}`, http.StatusInternalServerError)
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func newActivity(username string, w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var activity *workoutlog.Activity = &workoutlog.Activity{}
	if err := json.NewDecoder(r.Body).Decode(activity); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	activity, err = workoutlog.ActivityManager.NewActivity(username, activity.Name)
	if err != nil {
		if errors.Is(err, workoutlog.ErrActivityNameTaken) {
			http.Error(w, `{"message":"name is not unique"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"message":"failed to add activity"}`, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(activity)
	if err != nil {
		http.Error(w, `{"message":"failed to add activity"}`, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func listExercises(username, activityID string, w http.ResponseWriter) {
	activity := workoutlog.Activity{ID: activityID}
	err := activity.GetActivityExercises(username)
	if err != nil {
		http.Error(w, "failed to get exercises", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(activity.ExerciseIDs)
	if err != nil {
		http.Error(w, "failed to get exercises", http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func updateExercises(username, activityID string, w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var activity *workoutlog.Activity = &workoutlog.Activity{}
	if err := json.NewDecoder(r.Body).Decode(activity); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := workoutlog.ActivityManager.UpdateActivityExercises(username, *activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// only supports renaming right now
func updateActivity(username, activityID string, w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var activity *workoutlog.Activity = &workoutlog.Activity{}
	if err := json.NewDecoder(r.Body).Decode(activity); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if activity.Name != "" {
		err := workoutlog.ActivityManager.RenameActivity(username, *activity)
		//errTaken := workoutlog.ErrActivityNameTaken
		if err != nil {
			if errors.Is(err, workoutlog.ErrActivityNameTaken) {
				http.Error(w, `{"message":"name is not unique"}`, http.StatusBadRequest)
				return
			} else {
				http.Error(w, `{"message":"failed to rename activity"}`, http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
