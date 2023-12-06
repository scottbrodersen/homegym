package server

import (
	"encoding/json"
	"errors"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	programID, err := programs.ProgramManager.AddProgram(username, *program)
	if err != nil {
		if errors.Is(err, programs.ErrInvalidProgram) {
			http.Error(w, `{"message":"invalid program"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"message":"failed to add program"}`, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(struct {
		ID string `json:"id"`
	}{ID: *programID})

	if err != nil {
		http.Error(w, "failed to add program", http.StatusInternalServerError)
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
		if errors.Is(err, programs.ErrInvalidProgram) {
			http.Error(w, `{"message":"invalid program"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"message":"failed to add program"}`, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusOK)
}

func getProgram(username, activityID, programID string, w http.ResponseWriter, r *http.Request) {
	page, err := programs.ProgramManager.GetProgramsPageForActivity(username, activityID, programID, uint64(1))

	if err != nil {
		http.Error(w, `{"message":"failed to get programs"}`, http.StatusInternalServerError)
		return
	}

	body := []byte{}

	if len(page) > 0 {
		body, err = json.Marshal(page[0])
	}

	if err != nil {
		http.Error(w, `{"message":"failed to marshal program"}`, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func getProgramPage(username, activityID string, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	previousProgram, ok := r.Form["previous"]
	if !ok {
		http.Error(w, `{"message":"missing previous query parameter"}`, http.StatusBadRequest)
		return
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

	page, err := programs.ProgramManager.GetProgramsPageForActivity(username, activityID, previousProgram[0], uint64(pageSizeInt))

	if err != nil {
		http.Error(w, `{"message":"failed to get programs"}`, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(page)

	if err != nil {
		http.Error(w, `{"message":"failed to marshal programs"}`, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func addProgramInstance(username, activityID, programID string, w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var programInstance *programs.ProgramInstance = &programs.ProgramInstance{}
	if err := json.NewDecoder(r.Body).Decode(programInstance); err != nil {
		log.Error(err)
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

	programInstanceID, err := programs.ProgramManager.AddProgramInstance(username, *programInstance)
	if err != nil {
		if errors.Is(err, programs.ErrInvalidProgramInstance) {
			http.Error(w, `{"message":"invalid program"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"message":"failed to add program"}`, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(struct {
		ID string `json:"id"`
	}{ID: *programInstanceID})

	if err != nil {
		http.Error(w, "failed to add programn instance", http.StatusInternalServerError)
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
		http.Error(w, `{"message":"failed to add program instance"}`, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusOK)
}

func getProgramInstance(username, activityID, programID, instanceID string, w http.ResponseWriter, r *http.Request) {
	page, err := programs.ProgramManager.GetProgramInstancesPage(username, activityID, programID, instanceID, uint64(1))

	if err != nil {
		http.Error(w, `{"message":"failed to get program instance"}`, http.StatusInternalServerError)
		return
	}

	body := []byte{}

	if len(page) > 0 {
		body, err = json.Marshal(page[0])
	}

	if err != nil {
		http.Error(w, `{"message":"failed to marshal program instance"}`, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func getProgramInstancePage(username, activityID, programID string, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	previousProgramInstance, ok := r.Form["previous"]
	if !ok {
		http.Error(w, `{"message":"missing previous query parameter"}`, http.StatusBadRequest)
		return
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

	page, err := programs.ProgramManager.GetProgramInstancesPage(username, activityID, programID, previousProgramInstance[0], uint64(pageSizeInt))

	if err != nil {
		http.Error(w, `{"message":"failed to get program instances"}`, http.StatusInternalServerError)
		return
	}

	var body []byte = []byte{}

	if len(page) > 0 {
		body, err = json.Marshal(page)

		if err != nil {
			http.Error(w, `{"message":"failed to marshal program instances"}`, http.StatusInternalServerError)
			return
		}
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}

func setActiveProgramInstance(username, activityID, programID string, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	instanceID, ok := r.Form["id"]
	if !ok {
		http.Error(w, `{"message":"missing id query parameter"}`, http.StatusBadRequest)
		return
	}

	err := programs.ProgramManager.SetActiveProgramInstance(username, activityID, programID, instanceID[0])
	if err != nil {
		http.Error(w, `{"message":"failed set active program instance"}`, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.WriteHeader(http.StatusOK)
}

func getActiveProgramInstance(username, activityID, programID string, w http.ResponseWriter, r *http.Request) {
	activeInstance, err := programs.ProgramManager.GetActiveProgramInstance(username, activityID, programID)
	if err != nil {
		http.Error(w, `{"message":"failed to get active program instance"}`, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(activeInstance)

	if err != nil {
		http.Error(w, `{"message":"failed to marshal program instance"}`, http.StatusInternalServerError)
		return
	}

	h := w.Header()
	standardHeaders(&h)
	w.Write(body)
}
