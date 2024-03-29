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
	Mood:       1,
	Energy:     2,
	Motivation: 3,
	Overall:    4,
	Notes:      "test notes",
}

const testEventID = "test-event-id"

var exerciseType ExerciseType = testExerciseType()
var testWeight float32 = 10.5
var testDate int64 = time.Now().Unix()

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
	testDate := time.Now().Unix()
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
			db.On("AddEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			newEvent := newTestEvent()
			eventID, err := EventManager.NewEvent(testUserID, newEvent)

			So(err, ShouldBeNil)
			So(eventID, ShouldNotBeEmpty)
		})

		Convey("when we update an event", func() {
			db.On("GetEvent", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test event"), nil)
			db.On("AddEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			newEvent := newTestEvent()
			newEvent.ID = testEventID
			err := EventManager.UpdateEvent(testUserID, testDate, newEvent)

			So(err, ShouldBeNil)
		})

		Convey("when we add an exercise to the event", func() {
			db.On("AddExercisesToEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			testEvent := newTestEvent()
			jsonEvent, err := json.Marshal(testEvent)
			if err != nil {
				t.Fatal()
			}
			db.On("GetEvent", mock.Anything, mock.Anything, mock.Anything).Return(jsonEvent, nil)
			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testActivityName, []string{testExerciseID}, nil)

			testEtype := testExerciseType()
			eMgr.On("GetExerciseType", mock.Anything, mock.Anything).Return(&testEtype, nil)

			exerciseInstances := []ExerciseInstance{
				{
					TypeID:   testExerciseID,
					Index:    0,
					Segments: []ExerciseSegment{testSetsRepsSegmentMaker()},
				},
			}

			err = EventManager.AddExercisesToEvent(testUserID, testEvent.ID, time.Now().Unix(), exerciseInstances)

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
			//So(instances[testIngestedExerciseInstance.Index], ShouldResemble, testExportedExerciseInstance)
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
	})
}
