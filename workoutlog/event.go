package workoutlog

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
	"github.com/scottbrodersen/homegym/dal"
)

var DefaultPageSize = int(100)
var EventManager EventAdmin = new(eventManager)
var ErrInvalidEvent = fmt.Errorf("invalid event")

type EventAdmin interface {
	NewEvent(userID string, event Event) (*string, error)
	GetPageOfEvents(userID string, previousEvent Event, pageSize int) ([]Event, error)
	GetCachedExerciseType(exerciseTypeID string) *ExerciseType
	GetEventExercises(userID, eventID string) (map[int]ExerciseInstance, error)
	UpdateEvent(userID string, currentDate int64, event Event) error
	GetPageOfInstances(userID string, filter ExerciseFilter, pageSize int) ([]int64, [][]ExerciseInstance, error)
	DeleteEvent(userID string, event Event) error
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

func (em eventManager) GetCachedExerciseType(exerciseTypeID string) *ExerciseType {
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
func (em eventManager) NewEvent(userID string, event Event) (*string, error) {
	if userID == "" || event.ActivityID == "" || event.Date == 0 {
		return nil, ErrInvalidEvent
	}

	event.ID = uuid.New().String()

	eventJson, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// prepare exercise instances
	typeIDs, exercisesJSON, err := prepEventExercises(userID, event.ActivityID, event.Exercises)
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
func (em eventManager) UpdateEvent(userID string, currentDate int64, event Event) error {
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
	typeIDs, exercisesJSON, err := prepEventExercises(userID, event.ActivityID, event.Exercises)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	err = dal.DB.UpdateEvent(userID, event.ID, event.ActivityID, currentDate, event.Date, eventJson, typeIDs, exercisesJSON)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

func (em eventManager) DeleteEvent(userID string, event Event) error {
	if userID == "" || event.ID == "" || event.ActivityID == "" || event.Date == 0 {
		return ErrInvalidEvent
	}

	if err := dal.DB.DeleteEvent(userID, event.ID, event.ActivityID, event.Date); err != nil {
		log.WithError(err).Error("failed to delete event")
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func prepEventExercises(userID, activityID string, exerciseInstances map[int]ExerciseInstance) (map[int]string, map[int][]byte, error) {
	exInstances := map[int][]byte{}
	exTypeIDs := map[int]string{}

	for k, inst := range exerciseInstances {
		// check that the activity supports the exercise type
		err := checkActivityForExerciseType(userID, activityID, inst.TypeID)
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

// GetPageOfEvents gets a page of Event structs from the database.
// Set previousEvent to the last event returned in the previous page. Set it to an empty struct to get the first page.
// The maximum page size is 100. Exceeding the maximum returns an error.
func (em eventManager) GetPageOfEvents(userID string, previousEvent Event, pageSize int) ([]Event, error) {
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

func (em eventManager) GetEventExercises(userID, eventID string) (map[int]ExerciseInstance, error) {
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

// GetPageOfInstances gets a maximum of 3000 event instances
// The pageSize is slightly exceeded when the instances added from an event causes the total count to surpass the page limit.
// The returned map uses a date as the key, and the value is an array of instances that were performed on that date.
// The StartDate field of filter indicates the earliest date of the event to include in the page. Set to 0 to get the first page.
func (em eventManager) GetPageOfInstances(userID string, filter ExerciseFilter, pageSize int) ([]int64, [][]ExerciseInstance, error) {

	return getPageOfExercises(em, userID, filter, pageSize)
}

func getPageOfExercises(eventManager EventAdmin, userID string, filter ExerciseFilter, pageSize int) ([]int64, [][]ExerciseInstance, error) {
	if pageSize == 0 {
		pageSize = 3000 // exercise instances
	}

	dateStack := []int64{}
	instancesStack := [][]ExerciseInstance{}

	instanceCount := 0
	startEvent := Event{}

	if filter.StartDate == 0 {
		filter.StartDate = time.Now().Unix()
	}
	startEvent.Date = filter.StartDate

	for instanceCount < pageSize+1 {
		// get 100 events starting at the start date
		events, err := eventManager.GetPageOfEvents(userID, startEvent, DefaultPageSize)
		if err != nil {
			slog.Error("failed to get exercises", "user", userID, "error", err.Error())
			return nil, nil, err
		}

		if len(events) == 0 {
			break
		} else {
			// prepare for next page of events
			startEvent = events[len(events)-1]
		}

		// get the exercise instances for each event
		for _, evt := range events {
			if evt.Date <= filter.StartDate && evt.Date >= filter.EndDate {
				// TODO: limit the number of exercises in an event
				eventInstances, err := EventManager.GetEventExercises(userID, evt.ID)
				if err != nil {
					slog.Error("failed to get exercise instances", "user", userID, "error", err.Error())
					return nil, nil, err
				}
				// add them to the map if the pass through the filter
				instances := []ExerciseInstance{}
				for _, i := range eventInstances {
					if len(filter.ExerciseTypes) == 0 || slices.Index(filter.ExerciseTypes, i.TypeID) >= 0 {
						instances = append(instances, i)
						instanceCount++
					}
				}
				if len(instances) > 0 {
					dateStack = append(dateStack, evt.Date)
					instancesStack = append(instancesStack, instances)
				}
			}

		}
	}
	return dateStack, instancesStack, nil
}

func checkActivityForExerciseType(userID, activityID, exerciseTypeID string) error {
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

func NewMockEventAdmin() *MockEventAdmin {
	return new(MockEventAdmin)
}

type MockEventAdmin struct {
	mock.Mock
}

func (e *MockEventAdmin) GetCachedExerciseType(exerciseTypeID string) *ExerciseType {
	args := e.Called(exerciseTypeID)

	return args.Get(0).(*ExerciseType)

}

func (e *MockEventAdmin) NewEvent(userID string, event Event) (*string, error) {
	args := e.Called(userID, event)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*string), nil
}

func (e *MockEventAdmin) UpdateEvent(userID string, currentDate int64, event Event) error {
	args := e.Called(userID, currentDate, event)

	return args.Error(0)
}

func (e *MockEventAdmin) GetEventExercises(userID, eventID string) (map[int]ExerciseInstance, error) {
	args := e.Called(userID, eventID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[int]ExerciseInstance), nil
}

func (e *MockEventAdmin) GetPageOfEvents(userID string, previousEvent Event, pageSize int) ([]Event, error) {
	args := e.Called(userID, previousEvent, pageSize)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Event), nil
}

func (e *MockEventAdmin) GetPageOfInstances(userID string, filter ExerciseFilter, pageSize int) ([]int64, [][]ExerciseInstance, error) {
	args := e.Called(userID, filter, pageSize)

	if args.Error(2) != nil {
		return nil, nil, args.Error(1)

	}

	return args.Get(0).([]int64), args.Get(1).([][]ExerciseInstance), nil
}

func (e *MockEventAdmin) DeleteEvent(userID string, event Event) error {
	args := e.Called(userID, event)

	return args.Error(0)
}
