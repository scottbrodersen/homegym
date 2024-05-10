package workoutlog

import (
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/scottbrodersen/homegym/dal"
)

var DefaultPageSize = int(10)
var EventManager EventAdmin = new(eventManager)
var ErrInvalidEvent = fmt.Errorf("invalid event")

type EventAdmin interface {
	NewEvent(userID string, event Event) (*string, error)
	GetPageOfEvents(userID string, previousEvent Event, pageSize int) ([]Event, error)
	GetCachedExerciseType(exerciseTypeID string) *ExerciseType
	GetEventExercises(userID, eventID string) (map[int]ExerciseInstance, error)
	UpdateEvent(userID string, currentDate int64, event Event) error
}

type eventManager struct{}

type Event struct {
	ID         string `json:"id"`
	ActivityID string `json:"activityID"`
	Date       int64  `json:"date"`
	EventMeta
	Exercises map[int]ExerciseInstance `json:"exercises"` // key is the exercise index, ensures uniqueness
}

// zero value represents nil
type EventMeta struct {
	Mood       int    `json:"mood"`
	Motivation int    `json:"motivation"`
	Energy     int    `json:"energy"`
	Overall    int    `json:"overall"`
	Notes      string `json:"notes"`
}

var exerciseTypeCache sync.Map = sync.Map{}

func (em *eventManager) GetCachedExerciseType(exerciseTypeID string) *ExerciseType {
	cachedType, ok := exerciseTypeCache.Load(exerciseTypeID)
	if !ok {
		return nil
	}
	exerciseType := cachedType.(ExerciseType)

	return &exerciseType
}

// NewEvent adds a new event to the database
// The event data consists of everything but the exercises which are added subsequently in separate calls.
// If exercises are included in the event they are ignored.
// The event id is returned.
func (em *eventManager) NewEvent(userID string, event Event) (*string, error) {
	if userID == "" || event.ActivityID == "" || event.Date == 0 {
		return nil, ErrInvalidEvent
	}

	event.ID = uuid.New().String()

	eventJson, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// prepare exercise instances
	typeIDs, exercisesJSON, err := prepEventExercises(userID, event.ID, event.Date, event.Exercises)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	err = dal.DB.AddEvent(userID, event.ID, event.ActivityID, event.Date, eventJson, typeIDs, exercisesJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to add event: %w", err)
	}

	return &event.ID, nil
}

// UpdateEvent replaces a stored event.
func (em *eventManager) UpdateEvent(userID string, currentDate int64, event Event) error {
	if userID == "" || event.ID == "" || event.ActivityID == "" || event.Date == 0 {
		return ErrInvalidEvent
	}

	// Make sure the event exists
	existing, err := dal.DB.GetEvent(userID, event.ID, currentDate)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	if existing == nil {
		return ErrNotFound
	}

	eventJson, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// prepare exercise instances
	typeIDs, exercisesJSON, err := prepEventExercises(userID, event.ID, event.Date, event.Exercises)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	err = dal.DB.UpdateEvent(userID, event.ID, event.ActivityID, currentDate, event.Date, eventJson, typeIDs, exercisesJSON)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

func prepEventExercises(userID, eventID string, eventDate int64, exerciseInstances map[int]ExerciseInstance) (map[int]string, map[int][]byte, error) {
	exInstances := map[int][]byte{}
	exTypeIDs := map[int]string{}

	for k, inst := range exerciseInstances {
		// check that the activity supports the exercise type
		err := checkActivityForExerciseType(userID, eventID, inst.TypeID, eventDate)
		if err != nil {
			return nil, nil, err
		}
		exerciseType, err := ExerciseManager.GetExerciseType(userID, inst.TypeID)
		if err != nil {
			return nil, nil, err
		}

		err = exerciseType.validateInstance(&inst)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid exercise instance: %w", err)
		}

		instanceByte, err := json.Marshal(inst)
		if err != nil {
			log.WithError(err).Debug("failed to marshal exercise instance")
			return nil, nil, fmt.Errorf("failed to add exercise: %w", err)
		}
		exInstances[int(k)] = instanceByte
		exTypeIDs[int(k)] = inst.TypeID
	}

	return exTypeIDs, exInstances, nil
}

func (em *eventManager) GetPageOfEvents(userID string, previousEvent Event, pageSize int) ([]Event, error) {
	if pageSize > 100 {
		return nil, fmt.Errorf("page size cannot exceed 100")
	}

	page := pageSize
	if page == 0 {
		page = DefaultPageSize
	}

	eventsByte, err := dal.DB.GetEventPage(userID, previousEvent.ID, previousEvent.Date, page)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	events := []Event{}

	for _, v := range eventsByte {
		event := new(Event)
		if err := json.Unmarshal(v, event); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event %w", err)
		}

		events = append(events, *event)
	}

	for i := range events {
		events[i].Exercises, err = em.GetEventExercises(userID, events[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get event exercise instances: %w", err)
		}
	}

	return events, nil
}

func (em *eventManager) GetEventExercises(userID, eventID string) (map[int]ExerciseInstance, error) {
	storedInstances, err := dal.DB.GetEventExercises(userID, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to read event exercises: %w", err)
	}

	instances := map[int]ExerciseInstance{}

	for _, stored := range storedInstances {
		exerciseInstance := ExerciseInstance{}
		if err := json.Unmarshal(stored, &exerciseInstance); err != nil {
			return nil, fmt.Errorf("failed to unmarshal stored exercise instance: %w", err)
		}

		instances[exerciseInstance.Index] = exerciseInstance
	}

	return instances, nil
}

func checkActivityForExerciseType(userID, eventID, exerciseTypeID string, eventDate int64) error {

	eventByte, err := dal.DB.GetEvent(userID, eventID, eventDate)
	if err != nil {
		log.WithError(err).Debug("failed to get event")
		return fmt.Errorf("failed to check exercise belongs to event activity: %w", err)
	}
	if eventByte == nil {
		log.WithError(err).Debug("event not found")
		return fmt.Errorf("did not find the event: %w", err)
	}

	event := new(Event)
	if err := json.Unmarshal(eventByte, event); err != nil {
		log.WithError(err).Debug("failed to unmarshal stored event")
		return fmt.Errorf("failed to check exercise belongs to event activity: %w", err)
	}

	activityID := event.ActivityID
	_, exerciseIDs, err := dal.DB.ReadActivity(userID, activityID)
	if err != nil {
		log.WithError(err).Debug("failed to read activity")
		return fmt.Errorf("failed to check exercise belongs to event activity: %w", err)
	}

	found := false

	for _, id := range exerciseIDs {
		if id == exerciseTypeID {
			found = true

			break
		}
	}

	if !found {
		return fmt.Errorf("exercise type does not belong in the event's activity")
	}

	return nil
}
