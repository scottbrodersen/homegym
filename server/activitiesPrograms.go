package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/scottbrodersen/homegym/programs"
)

func newProgram(username, activityID string, w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var program *programs.Program = &programs.Program{}
	if err := json.NewDecoder(r.Body).Decode(program); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if program.ActivityID != activityID {
		log.Error("wrong activity ID in path")
		http.Error(w, `{"message":"wrong activity ID in path"}`, http.StatusBadRequest)
		return
	}

	programID, err := programs.ProgramManager.AddProgram(username, *program)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errors.As(err, new(programs.ErrInvalidProgram)) {
			http.Error(w, fmt.Sprintf(`{"message":"%s"}`, jsonSafeError(err)), http.StatusBadRequest)
			return
		}
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(struct {
		ID string `json:"id"`
	}{ID: *programID})

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func updateProgram(username, activityID, programID string, w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var program *programs.Program = &programs.Program{}
	if err := json.NewDecoder(r.Body).Decode(program); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if program.ID != programID {
		log.Error("wrong program ID in path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if program.ActivityID != activityID {
		log.Error("wrong activity ID in path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := programs.ProgramManager.UpdateProgram(username, *program)
	if err != nil {
		if errors.As(err, new(programs.ErrInvalidProgram)) {
			http.Error(w, `{"message":"invalid program"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusOK)
}

func getProgram(username, activityID, programID string, w http.ResponseWriter) {
	page, err := programs.ProgramManager.GetProgramsPageForActivity(username, activityID, programID, int(1))

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	body := []byte{}

	if len(page) > 0 {
		body, err = json.Marshal(page[0])
	}

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func getProgramPage(username, activityID string, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	previousID := ""
	previousProgram, ok := r.Form["previous"]
	if ok {
		previousID = previousProgram[0]
	}
	pageSize, ok := r.Form["size"]
	if !ok {
		http.Error(w, `{"message":"missing size query parameter"}`, http.StatusBadRequest)
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize[0])
	if err != nil {
		http.Error(w, `{"message":"invalid page size"}`, http.StatusBadRequest)
		return
	}

	page, err := programs.ProgramManager.GetProgramsPageForActivity(username, activityID, previousID, int(pageSizeInt))

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(page)

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func addProgramInstance(username, activityID, programID string, w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, `{"message":"no body"}`, http.StatusBadRequest)
		return
	}

	var programInstance *programs.ProgramInstance = &programs.ProgramInstance{}

	if err := json.NewDecoder(r.Body).Decode(programInstance); err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprintf(`{"message":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if programInstance.ProgramID != programID {
		http.Error(w, `{"message":"wrong program ID in path"}`, http.StatusBadRequest)
		return
	}

	if programInstance.StartTime == 0 {
		http.Error(w, `{"message":"no start date"}`, http.StatusBadRequest)
		return
	}

	err := programs.ProgramManager.AddProgramInstance(username, programInstance)
	if err != nil {
		if errors.Is(err, programs.ErrInvalidProgramInstance) {
			http.Error(w, fmt.Sprintf(`{"message":"invalid program: %s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(programInstance)

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func updateProgramInstance(username, activityID, programID, instanceID string, w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var programInstance *programs.ProgramInstance = &programs.ProgramInstance{}
	if err := json.NewDecoder(r.Body).Decode(programInstance); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if programInstance.ID != instanceID {
		log.Error("wrong program instance ID in path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if programInstance.ActivityID != activityID {
		log.Error("wrong activity ID in path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if programInstance.ProgramID != programID {
		log.Error("wrong program ID in path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := programs.ProgramManager.UpdateProgramInstance(username, *programInstance)
	if err != nil {
		if errors.Is(err, programs.ErrInvalidProgramInstance) {
			http.Error(w, `{"message":"invalid program instance"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusOK)
}

func getProgramInstance(username, programID, instanceID string, w http.ResponseWriter) {
	page, err := programs.ProgramManager.GetProgramInstancesPage(username, programID, instanceID, int(1))

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	body := []byte{}

	if len(page) > 0 {
		body, err = json.Marshal(page[0])
	}

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func getProgramInstancePage(username, programID string, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	previousID := ""
	previousProgramInstance, ok := r.Form["previous"]
	if ok {
		previousID = previousProgramInstance[0]
	}

	pageSize, ok := r.Form["size"]
	if !ok {
		http.Error(w, `{"message":"missing size query parameter"}`, http.StatusBadRequest)
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize[0])
	if err != nil {
		http.Error(w, `{"message":"invalid page size"}`, http.StatusBadRequest)
		return
	}

	page, err := programs.ProgramManager.GetProgramInstancesPage(username, programID, previousID, int(pageSizeInt))

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var body []byte = []byte{}

	if len(page) > 0 {
		body, err = json.Marshal(page)

		if err != nil {
			http.Error(w, internalServerError, http.StatusInternalServerError)
			return
		}
	} else {
		body = []byte("[]")
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func setActiveProgramInstance(username, activityID, programID, instanceID string, w http.ResponseWriter, r *http.Request) {

	err := programs.ProgramManager.SetActiveProgramInstance(username, activityID, programID, instanceID)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusOK)
}

func getActiveProgramInstance(username, activityID string, w http.ResponseWriter) {
	activeInstance, err := programs.ProgramManager.GetActiveProgramInstance(username, activityID)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	if activeInstance == nil {
		h := w.Header()
		standardHeaders(&h)
		w.Write([]byte("{}"))
		return
	}

	body, err := json.Marshal(activeInstance)

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}
