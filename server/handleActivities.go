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
	rxpPrograms := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]*)/programs/?$", rootpath))
	rxpProgramsID := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]*)/programs/([a-zA-Z0-9-]*)/?$", rootpath))
	rxpProgramInstances := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]*)/programs/([a-zA-Z0-9-]*)/instances/?$", rootpath))
	rxpProgramInstancesID := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]*)/programs/([a-zA-Z0-9-]*)/instances/([a-zA-Z0-9-]{7,})/?$", rootpath))
	rxpProgramInstancesActive := regexp.MustCompile(fmt.Sprintf("^%s([a-zA-Z0-9-]*)/programs/instances/active/?$", rootpath))

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
			updateExercises(*username, w, r)
			return
		}
	} else if rxpActivity.MatchString(r.URL.Path) {
		if r.Method == http.MethodPost {
			updateActivity(*username, w, r)
			return
		}
	} else if rxpPrograms.MatchString(r.URL.Path) {
		activityID := rxpPrograms.FindStringSubmatch(r.URL.Path)[1]
		if r.Method == http.MethodPost {
			newProgram(*username, activityID, w, r)
			return
		} else if r.Method == http.MethodGet {
			getProgramPage(*username, activityID, w, r)
			return
		}
	} else if rxpProgramsID.MatchString(r.URL.Path) {
		ids := rxpProgramsID.FindStringSubmatch(r.URL.Path)
		activityID := ids[1]
		programID := ids[2]
		if r.Method == http.MethodPost {
			updateProgram(*username, activityID, programID, w, r)
			return
		} else if r.Method == http.MethodGet {
			getProgram(*username, activityID, programID, w)
			return
		}
	} else if rxpProgramInstances.MatchString(r.URL.Path) {
		ids := rxpProgramInstances.FindStringSubmatch(r.URL.Path)
		activityID := ids[1]
		programID := ids[2]

		if r.Method == http.MethodPost {
			addProgramInstance(*username, activityID, programID, w, r)
			return
		} else if r.Method == http.MethodGet {
			getProgramInstancePage(*username, programID, w, r)
			return
		}
	} else if rxpProgramInstancesID.MatchString(r.URL.Path) {
		ids := rxpProgramInstancesID.FindStringSubmatch(r.URL.Path)
		activityID := ids[1]
		programID := ids[2]
		instanceID := ids[3]

		if r.Method == http.MethodPost {
			updateProgramInstance(*username, activityID, programID, instanceID, w, r)
			return
		} else if r.Method == http.MethodGet {
			getProgramInstance(*username, programID, instanceID, w)
			return
		}
	} else if rxpProgramInstancesActive.MatchString(r.URL.Path) {
		ids := rxpProgramInstancesActive.FindStringSubmatch(r.URL.Path)
		activityID := ids[1]

		if r.Method == http.MethodPost {
			programID := r.URL.Query().Get("programid")
			if programID == "" {
				http.Error(w, `{"message":"missing programid query parameter"}`, http.StatusBadRequest)
				return
			}

			instanceID := r.URL.Query().Get("instanceid")
			if instanceID == "" {
				http.Error(w, `{"message":"missing instanceid query parameter"}`, http.StatusBadRequest)
				return
			}

			setActiveProgramInstance(*username, activityID, programID, instanceID, w, r)

			return
		} else if r.Method == http.MethodGet {
			getActiveProgramInstance(*username, activityID, w)

			return
		}
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func listActivities(username string, w http.ResponseWriter) {
	activities, err := workoutlog.ActivityManager.GetActivityNames(username)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	body, err := json.Marshal(activities)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
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
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(activity)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
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
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(activity.ExerciseIDs)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func updateExercises(username string, w http.ResponseWriter, r *http.Request) {
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
func updateActivity(username string, w http.ResponseWriter, r *http.Request) {
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
				http.Error(w, internalServerError, http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
