package programs

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/scottbrodersen/homegym/dal"
)

// ProgramInstance contains details of the performance of a program
// StartEvent is the ID of the event that performed the first workout in the program
// Fulfillment maps program workouts to events
// The key of the fulmillment map is {blockIndex}#{microCycleIndex}#{workoutIndex}
type ProgramInstance struct {
	ID          string            `json:"id,omitempty"`
	ProgramID   string            `json:"programID"`
	Title       string            `json:"title"`
	StartEvent  string            `json:"startEvent,omitempty"`
	ActivityID  string            `json:"activityID"`
	Active      bool              `json:"-"`
	Complete    bool              `json:"complete,omitempty"`
	Fulfillment map[string]string `json:"fulfillment,omitempty"`
}

var ErrInvalidProgram = errors.New("invalid program")
var ErrInvalidProgramInstance = errors.New("invalid program instance")

func (pi ProgramInstance) validate() error {
	if pi.ID == "" {
		return errors.Join(ErrInvalidProgramInstance, fmt.Errorf("missing ID"))
	}

	if pi.ProgramID == "" {
		return errors.Join(ErrInvalidProgramInstance, fmt.Errorf("missing program ID"))
	}

	if pi.Title == "" {
		return errors.Join(ErrInvalidProgramInstance, fmt.Errorf("missing title"))
	}
	if pi.ActivityID == "" {
		return errors.Join(ErrInvalidProgramInstance, fmt.Errorf("missing activity ID"))
	}
	return nil
}

type WorkoutSegment struct {
	ExerciseTypeID string `json:"exerciseTypeID"`
	Prescription   string `json:"prescription"`
}

func (ws WorkoutSegment) validate() error {
	if ws.ExerciseTypeID == "" {
		return fmt.Errorf("missing exercise type ID")
	}

	if ws.Prescription == "" {
		return fmt.Errorf("missing prescription")
	}

	return nil
}

type Workout struct {
	Title     string           `json:"title"`
	Segments  []WorkoutSegment `json:"segments"`
	Intensity string           `json:"intensity,omitempty"`
}

func (w Workout) validate() error {
	if w.Title == "" {
		return fmt.Errorf("missing title")
	}
	if w.Segments == nil || len(w.Segments) == 0 {
		return fmt.Errorf("missing segments")
	}

	for _, s := range w.Segments {
		if err := s.validate(); err != nil {
			return fmt.Errorf("invalid segment: %w", err)
		}
	}

	return nil
}

type MicroCycle struct {
	Title     string    `json:"title"`
	Span      int       `json:"span"`
	Intensity string    `json:"intensity,omitempty"`
	Workouts  []Workout `json:"workouts,omitempty"`
}

func (mc MicroCycle) validate() error {
	if mc.Title == "" {
		return fmt.Errorf("missing title")
	}
	if mc.Span == 0 {
		return fmt.Errorf("missing span")
	}

	for _, w := range mc.Workouts {
		if err := w.validate(); err != nil {
			return fmt.Errorf("invalid workout: %w", err)
		}
	}
	return nil
}

type Block struct {
	Title       string       `json:"title"`
	MicroCycles []MicroCycle `json:"microCycles,omitempty"`
	Intensity   string       `json:"intensity,omitempty"`
}

func (b Block) validate() error {
	if b.Title == "" {
		return fmt.Errorf("missing title")
	}
	for _, m := range b.MicroCycles {
		if err := m.validate(); err != nil {
			return fmt.Errorf("invalid microcycle: %w", err)
		}
	}

	return nil
}

// Program defines a training program.
// Programs must be associated with an activity.
// The intent is to define the structure of the program in blocks, microcycles (weeks), and workouts.
// The exercises are explicitly specified but the intensity and volume are descriptive.
// Intensity can be provided at each sub-phase of a program to enable progressively precise descriptiions.
type Program struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	ActivityID string  `json:"activityID"`
	Blocks     []Block `json:"blocks,omitempty"`
}

func (p Program) validate() error {
	if p.ID == "" {
		return errors.Join(ErrInvalidProgram, fmt.Errorf("missing ID"))
	}

	if p.Title == "" {
		return errors.Join(ErrInvalidProgram, fmt.Errorf("missing title"))
	}

	if p.ActivityID == "" {
		return errors.Join(ErrInvalidProgram, fmt.Errorf("missing activity ID"))
	}
	for _, b := range p.Blocks {
		if err := b.validate(); err != nil {
			return errors.Join(ErrInvalidProgram, err)
		}
	}

	return nil
}

type ProgramAdmin interface {
	AddProgram(userID string, program Program) (*string, error)
	UpdateProgram(userID string, program Program) error
	GetProgramsPageForActivity(userID, activityID, previousProgramID string, pageSize int) ([]Program, error)
	AddProgramInstance(userID string, instance ProgramInstance) (*string, error)
	UpdateProgramInstance(userID string, instance ProgramInstance) error
	GetProgramInstancesPage(userID, activityID, programID, previousProgramInstanceID string, pageSize int) ([]ProgramInstance, error)
	SetActiveProgramInstance(userID, activityID, programID, instanceID string) error
	GetActiveProgramInstance(userID, activityID, programID string) (*ProgramInstance, error)
}

type ProgramUtil struct{}

var ProgramManager ProgramAdmin = ProgramUtil{}

func (pu ProgramUtil) AddProgram(userID string, program Program) (*string, error) {
	if program.ID != "" {
		return nil, fmt.Errorf("new programs cannot have an ID")
	}

	program.ID = uuid.New().String()

	if err := program.validate(); err != nil {
		return nil, err
	}

	// make sure the activity exists
	activityName, _, err := dal.DB.ReadActivity(userID, program.ActivityID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate activity ID: %w", err)
	}

	if activityName == nil {
		return nil, fmt.Errorf("activity does not exist")
	}

	programJSON, err := json.Marshal(program)
	if err != nil {
		return nil, fmt.Errorf("failed to parse program: %w", err)
	}

	if err := dal.DB.AddProgram(userID, program.ActivityID, program.ID, programJSON); err != nil {
		return nil, fmt.Errorf("failed to add program: %w", err)
	}

	return &program.ID, nil
}

func (pu ProgramUtil) UpdateProgram(userID string, program Program) error {
	if err := program.validate(); err != nil {
		return err
	}

	// Make sure the activity exists
	activityName, _, err := dal.DB.ReadActivity(userID, program.ActivityID)
	if err != nil {
		return fmt.Errorf("failed to validate activity ID: %w", err)
	}

	if activityName == nil {
		return fmt.Errorf("activity does not exist")
	}

	// Make sure the program exists
	existing, err := dal.DB.GetProgramPage(userID, program.ActivityID, program.ID, 1)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	if len(existing) == 0 {
		return fmt.Errorf("program not found")
	}

	programJSON, err := json.Marshal(program)
	if err != nil {
		return fmt.Errorf("failed to update program: %w", err)
	}

	if err := dal.DB.AddProgram(userID, program.ActivityID, program.ID, programJSON); err != nil {
		return fmt.Errorf("failed to update program: %w", err)
	}

	return nil
}

func (pu ProgramUtil) GetProgramsPageForActivity(userID, activityID, previousProgramID string, pageSize int) ([]Program, error) {
	numToGet := pageSize
	if numToGet > 100 {
		numToGet = 100
	}

	programsByte, err := dal.DB.GetProgramPage(userID, activityID, previousProgramID, numToGet)
	if err != nil {
		return nil, fmt.Errorf("failed to get programs: %w", err)
	}

	if len(programsByte) == 0 {
		return nil, nil
	}

	programs := []Program{}

	for _, p := range programsByte {
		program := new(Program)
		err := json.Unmarshal(p, program)
		if err != nil {
			return nil, fmt.Errorf("failed to parse stored program: %w", err)
		}
		programs = append(programs, *program)
	}

	return programs, nil
}

func (pu ProgramUtil) AddProgramInstance(userID string, instance ProgramInstance) (*string, error) {
	if instance.ID != "" {
		return nil, fmt.Errorf("new program instances cannot have an ID")
	}

	// Make sure the activity exists
	activityName, _, err := dal.DB.ReadActivity(userID, instance.ActivityID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate activity ID: %w", err)
	}

	if activityName == nil {
		return nil, fmt.Errorf("activity does not exist")
	}

	// Make sure the program exists
	existing, err := dal.DB.GetProgramPage(userID, instance.ActivityID, instance.ID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to validate programID: %w", err)
	}

	if len(existing) == 0 {
		return nil, fmt.Errorf("program not found")
	}

	instance.ID = uuid.New().String()

	if err := instance.validate(); err != nil {
		return nil, err
	}

	instanceJSON, err := json.Marshal(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to parse program instance: %w", err)
	}

	if err := dal.DB.AddProgramInstance(userID, instance.ActivityID, instance.ProgramID, instance.ID, instanceJSON); err != nil {
		return nil, fmt.Errorf("failed to add program: %w", err)
	}

	return &instance.ID, nil
}

func (pu ProgramUtil) UpdateProgramInstance(userID string, instance ProgramInstance) error {
	if err := instance.validate(); err != nil {
		return err
	}

	existing, err := dal.DB.GetProgramInstancePage(userID, instance.ActivityID, instance.ProgramID, instance.ID, 1)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("program instance not found")
	}

	programJSON, err := json.Marshal(instance)
	if err != nil {
		return fmt.Errorf("failed to parse program instance: %w", err)
	}

	if err := dal.DB.AddProgramInstance(userID, instance.ActivityID, instance.ProgramID, instance.ID, programJSON); err != nil {
		return fmt.Errorf("failed to update program: %w", err)
	}

	return nil

}

func (pu ProgramUtil) GetProgramInstancesPage(userID, activityID, programID, previousProgramInstanceID string, pageSize int) ([]ProgramInstance, error) {
	numToGet := pageSize
	if numToGet > 100 {
		numToGet = 100
	}

	instancesByte, err := dal.DB.GetProgramInstancePage(userID, activityID, programID, previousProgramInstanceID, numToGet)
	if err != nil {
		return nil, fmt.Errorf("failed to get program instances: %w", err)
	}

	if len(instancesByte) == 0 {
		return nil, nil
	}
	instances := []ProgramInstance{}
	for _, p := range instancesByte {
		instance := new(ProgramInstance)
		err := json.Unmarshal(p, instance)
		if err != nil {
			return nil, fmt.Errorf("failed to parse program instance: %w", err)
		}
		instances = append(instances, *instance)
	}

	return instances, nil
}

func (pu ProgramUtil) SetActiveProgramInstance(userID, activityID, programID, instanceID string) error {
	if err := dal.DB.SetActiveProgramInstance(userID, activityID, programID, instanceID); err != nil {
		return err
	}

	return nil
}

func (pu ProgramUtil) GetActiveProgramInstance(userID, activityID, programID string) (*ProgramInstance, error) {
	instanceBytes, err := dal.DB.GetActiveProgramInstance(userID, activityID, programID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active program instance: %w", err)
	}

	instance := new(ProgramInstance)

	// TODO: unmarshal here, then marshal in handler. really?
	if err := json.Unmarshal(instanceBytes, instance); err != nil {
		return nil, fmt.Errorf("failed to parse program instance: %w", err)
	}

	return instance, nil
}
