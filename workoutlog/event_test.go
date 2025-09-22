package workoutlog

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/scottbrodersen/homegym/dal"
)

var testEventMeta EventMeta = EventMeta{
	Overall: 4,
	Notes:   "test notes",
}

const testEventID = "test-event-id"

var exerciseType ExerciseType = testExerciseType()
var testWeight float32 = 10.5

const testDate int64 = int64(1719669151)

func newTestEvent() Event {
	return Event{
		Date:       testDate,
		ActivityID: testActivityID,
		EventMeta:  testEventMeta,
	}
}

func testSetsRepsSegmentMaker() ExerciseSegment {
	reps := []float32{}
	choices := []float32{0, 1}

	for i := 0; i < 3; i++ {
		reps = append(reps, choices[rand.Intn(len(choices))])
	}

	volume := [][]float32{reps}
	intensity := testWeight

	return ExerciseSegment{Intensity: intensity, Volume: volume}
}

func getEventPageReturnValuesMaker() ([][]byte, error) {
	//testDate := time.Now().Unix()
	events := make([][]byte, 10)

	for i := 0; i < 10; i++ {
		id := uuid.New().String()
		date := testDate + int64(i*100)
		event := Event{
			ID:         id,
			Date:       date,
			ActivityID: testActivityID,
			EventMeta:  testEventMeta,
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		events[i] = eventBytes
	}

	return events, nil
}

func TestEvents(t *testing.T) {
	exerciseTypeCache.Store(exerciseType.ID, exerciseType)

	Convey("Given a dal client and an exercise manager", t, func() {
		db := dal.NewMockDal()
		dal.DB = db

		eMgr := NewMockExerciseManager()
		ExerciseManager = eMgr

		Convey("when we create an event", func() {
			db.On("AddEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			newEvent := newTestEvent()
			eventID, err := EventManager.NewEvent(testUserID, newEvent)

			So(err, ShouldBeNil)
			So(eventID, ShouldNotBeEmpty)
		})

		Convey("when we update an event", func() {
			db.On("GetEvent", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test event"), nil)
			db.On("UpdateEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			newEvent := newTestEvent()
			newEvent.ID = testEventID
			err := EventManager.UpdateEvent(testUserID, testDate, newEvent)

			So(err, ShouldBeNil)
		})

		Convey("When we get an event's exercise instances", func() {
			exerciseBytes, err := json.Marshal(testExerciseInstance)
			if err != nil {
				t.Fatal()
			}
			instancesByte := [][]byte{exerciseBytes}

			db.On("GetEventExercises", mock.Anything, mock.Anything).Return(instancesByte, err)
			eMgr.On("GetExerciseType", mock.Anything, mock.Anything).Return(&exerciseType, nil)

			instances, err := EventManager.GetEventExercises(testUserID, testEventID)
			So(err, ShouldBeNil)
			So(instances, ShouldNotBeEmpty)
			So(len(instances), ShouldEqual, len(instancesByte))
		})

		Convey("When we get a page of events", func() {
			exerciseBytes, err := json.Marshal(testExerciseInstance)

			if err != nil {
				t.Fatal()
			}
			instancesByte := [][]byte{exerciseBytes}
			db.On("GetEventExercises", mock.Anything, mock.Anything).Return(instancesByte, err)

			eventsByte, err := getEventPageReturnValuesMaker()
			db.On("GetEventPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(eventsByte, err)

			events, err := EventManager.GetPageOfEvents(testUserID, Event{}, int(10))

			So(err, ShouldBeNil)
			So(len(events), ShouldEqual, 10)
		})

		Convey("When we delete an event", func() {
			db.On("DeleteEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			testEvent := newTestEvent()
			testEvent.ID = "testID"
			testEvent.Exercises = map[int]ExerciseInstance{0: testExerciseInstance}
			err := EventManager.DeleteEvent(testUserID, testEvent)

			So(err, ShouldBeNil)
		})
	})

	Convey("Given an event manager", t, func() {
		mockEM := NewMockEventAdmin()
		numEvents := 4
		testEvents = []Event{}
		for i := 0; i < numEvents; i++ {
			testEvents = append(testEvents, testEventMaker(testDate+int64(i)))
		}

		testEventExercises := testEvents[0].Exercises

		Convey("When we get a page of exercise instances, unfiltered, and the page is not full", func() {
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return(testEvents, nil).Once()
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return([]Event{}, nil).Once()
			mockEM.On("GetEventExercises", mock.Anything, mock.Anything).Return(testEventExercises, nil).Times(len(testEvents))
			EventManager = mockEM
			filter := ExerciseFilter{StartDate: time.Now().Unix()}
			dates, instances, err := getPageOfExercises(EventManager, testUserID, filter, 0)

			So(err, ShouldBeNil)
			So(len(dates), ShouldEqual, numEvents)
			So(len(instances), ShouldEqual, numEvents)
			for _, inst := range instances {
				So(len(inst), ShouldEqual, 3)
			}
		})

		Convey("When we get a full page of exercise instances, unfiltered", func() {
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return(testEvents, nil).Once()
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return([]Event{}, nil).Once()
			mockEM.On("GetEventExercises", mock.Anything, mock.Anything).Return(testEventExercises, nil).Times(len(testEvents))
			EventManager = mockEM
			filter := ExerciseFilter{StartDate: time.Now().Unix()}
			dates, instances, err := getPageOfExercises(EventManager, testUserID, filter, 4)

			So(err, ShouldBeNil)
			So(len(dates), ShouldEqual, 4)
			So(len(instances), ShouldEqual, 4)
			for _, inst := range instances {
				So(len(inst), ShouldEqual, 3)
			}
		})

		Convey("When we get instances from multiple pages of events, unfiltered", func() {
			// force event page size of 1
			for i := 0; i < numEvents; i++ {
				mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return([]Event{testEvents[i]}, nil).Once()
			}
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return([]Event{}, nil).Once()

			mockEM.On("GetEventExercises", mock.Anything, mock.Anything).Return(testEventExercises, nil).Times(len(testEvents))
			EventManager = mockEM
			filter := ExerciseFilter{StartDate: time.Now().Unix()}
			dates, instances, err := getPageOfExercises(EventManager, testUserID, filter, 0)

			So(err, ShouldBeNil)
			So(len(dates), ShouldEqual, numEvents)
			So(len(instances), ShouldEqual, numEvents)
			for _, inst := range instances {
				So(len(inst), ShouldEqual, 3)
			}
		})

		Convey("When we filter to return a single exercise type", func() {
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return(testEvents, nil).Once()
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return([]Event{}, nil).Once()
			mockEM.On("GetEventExercises", mock.Anything, mock.Anything).Return(testEventExercises, nil).Times(len(testEvents))
			EventManager = mockEM
			filter := ExerciseFilter{StartDate: time.Now().Unix(), ExerciseTypes: []string{"id1"}}
			dates, instances, err := getPageOfExercises(EventManager, testUserID, filter, 0)

			So(err, ShouldBeNil)
			So(len(dates), ShouldEqual, numEvents)
			So(len(instances), ShouldEqual, numEvents)
			for _, inst := range instances {
				So(len(inst), ShouldEqual, 1)
				So(inst[0].TypeID, ShouldEqual, "id1")
			}
		})

		Convey("When we filter on multiple exercise types", func() {
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return(testEvents, nil).Once()
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return([]Event{}, nil).Once()
			mockEM.On("GetEventExercises", mock.Anything, mock.Anything).Return(testEventExercises, nil).Times(len(testEvents))
			EventManager = mockEM
			filter := ExerciseFilter{StartDate: time.Now().Unix(), ExerciseTypes: []string{"id1"}}
			dates, instances, err := getPageOfExercises(EventManager, testUserID, filter, 0)

			So(err, ShouldBeNil)
			So(len(dates), ShouldEqual, numEvents)
			So(len(instances), ShouldEqual, numEvents)
			for _, inst := range instances {
				So(len(inst), ShouldEqual, 1)
				So(inst[0].TypeID, ShouldEqual, "id1")
			}
		})

		Convey("When we set the end date", func() {
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return(testEvents, nil).Once()
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return([]Event{}, nil).Once()
			mockEM.On("GetEventExercises", mock.Anything, mock.Anything).Return(testEventExercises, nil).Times(len(testEvents))
			EventManager = mockEM
			filter := ExerciseFilter{EndDate: testEvents[1].Date}
			dates, instances, err := getPageOfExercises(EventManager, testUserID, filter, 0)

			So(err, ShouldBeNil)
			So(len(dates), ShouldEqual, 3)
			So(len(instances), ShouldEqual, 3)
			for _, inst := range instances {
				So(len(inst), ShouldEqual, 3)
			}
		})

		Convey("When we set a start date", func() {
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return(testEvents, nil).Once()
			mockEM.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return([]Event{}, nil).Once()
			mockEM.On("GetEventExercises", mock.Anything, mock.Anything).Return(testEventExercises, nil).Times(len(testEvents))
			EventManager = mockEM
			filter := ExerciseFilter{StartDate: testEvents[1].Date}
			dates, instances, err := getPageOfExercises(EventManager, testUserID, filter, 0)

			So(err, ShouldBeNil)
			So(len(dates), ShouldEqual, 2)
			So(len(instances), ShouldEqual, 2)
			for _, inst := range instances {
				So(len(inst), ShouldEqual, 3)
			}

		})
	})
}
