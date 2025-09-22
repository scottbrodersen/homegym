package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/scottbrodersen/homegym/workoutlog"
)

var testExType workoutlog.ExerciseType = testExerciseTypeNonComposite()

const (
	testTime           = 1719824281
	testIntensityValue = 5.0
	numSets            = 4
	numReps            = 1
	url                = "/homegym/api/events/"
)

func testEvents(number int) []workoutlog.Event {
	testEventMeta := workoutlog.EventMeta{
		Overall: 4,
		Notes:   "test note",
	}

	testInstance := workoutlog.ExerciseInstance{
		TypeID:   testExType.ID,
		Segments: []workoutlog.ExerciseSegment{},
	}

	sets := []float32{}
	for i := 0; i < numSets; i++ {
		sets = append(sets, float32(numReps))
	}
	vol := [][]float32{sets}

	testInstance.Segments = append(testInstance.Segments, workoutlog.ExerciseSegment{
		Intensity: testIntensityValue,
		Volume:    vol,
	})

	testExercises := map[int]workoutlog.ExerciseInstance{1: testInstance}

	events := []workoutlog.Event{}

	for i := 0; i < number; i++ {
		events = append(events, workoutlog.Event{
			ID:         fmt.Sprintf("test-event-%d", i),
			ActivityID: "test-activity-id",
			Date:       testTime + int64(i),
			EventMeta:  testEventMeta,
			Exercises:  testExercises,
		})
	}

	return events
}

func TestHandleEvents(t *testing.T) {

	Convey("Given an event manager and an exercise manager", t, func() {

		mockEventManager := workoutlog.NewMockEventAdmin()
		workoutlog.EventManager = mockEventManager

		mockExerciseManager := workoutlog.NewMockExerciseManager()
		workoutlog.ExerciseManager = mockExerciseManager

		Convey("When we add an event", func() {
			eventID := "test-event"
			mockEventManager.On("NewEvent", mock.Anything, mock.Anything).Return(&eventID, nil)

			testEvent := `{"date": 1234, "activityID": "test-activity-id","mood": 1}`

			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(testEvent)))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			EventsApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)

			expectedEvent := &workoutlog.Event{}
			_ = json.Unmarshal([]byte(testEvent), expectedEvent)

			So(mockEventManager.AssertCalled(t, "NewEvent", GymContextValue(testContext(), usernameKey), *expectedEvent), ShouldBeTrue)
		})

		Convey("When we update an event", func() {
			eventID := "test-event-id"
			currentDate := "1234"
			mockEventManager.On("UpdateEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			testEvent := `{"date": 1234, "id": "` + eventID + `", "activityID": "test-activity-id","mood": 1}`

			reqUrl := fmt.Sprintf("%s%s/%s", url, currentDate, eventID)

			req := httptest.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer([]byte(testEvent)))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			EventsApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusNoContent)

			eventStruct := &workoutlog.Event{}
			_ = json.Unmarshal([]byte(testEvent), eventStruct)

			So(mockEventManager.Calls[0].Method, ShouldEqual, "UpdateEvent")
			updateEventCalledArgs := mockEventManager.Calls[0].Arguments
			So(updateEventCalledArgs.Get(0), ShouldEqual, GymContextValue(testContext(), usernameKey))
			So(updateEventCalledArgs.Get(1), ShouldEqual, 1234)
			So(updateEventCalledArgs.Get(2), ShouldResemble, *eventStruct)
		})

		Convey("When we get event exercises", func() {
			event := testEvents(1)[0]

			mockEventManager.On("GetEventExercises", mock.Anything, mock.Anything, mock.Anything).Return(event.Exercises, nil)

			reqUrl := fmt.Sprintf("%s13456/test-event-id/exercises", url)

			req := httptest.NewRequest(http.MethodGet, reqUrl, nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			EventsApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)

			returnedExercises := new(map[int]workoutlog.ExerciseInstance)
			mockEventManager.On("GetCachedExerciseType", mock.Anything).Return(&testExType, nil)
			So(json.NewDecoder(w.Result().Body).Decode(returnedExercises), ShouldBeNil)
			So(*returnedExercises, ShouldResemble, event.Exercises)
		})

		Convey("When we delete an event", func() {
			event := testEvents(1)[0]
			eventJSON, err := json.Marshal(event)
			if err != nil {
				t.Fail()
			}

			mockEventManager.On("DeleteEvent", mock.Anything, mock.Anything).Return(nil)

			reqUrl := fmt.Sprintf("%s%d/%s/", url, event.Date, event.ID)

			req := httptest.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(eventJSON))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			EventsApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusNoContent)
		})

		Convey("When we get a page of events", func() {

			events := testEvents(2)
			mockEventManager.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return(events, nil)
			mockExerciseManager.On("GetExerciseType", mock.Anything, mock.Anything).Return(&testExType, nil)
			mockEventManager.On("GetCachedExerciseType", mock.Anything).Return(&testExType, nil).Times(len(events))

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			EventsApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			returnedEvents := []workoutlog.Event{}

			if err := json.NewDecoder(w.Result().Body).Decode(&returnedEvents); err != nil {
				t.Fail()
			}

			So(returnedEvents, ShouldResemble, events)
		})

		Convey("When we get a page of metrics", func() {
			numEvents := 5
			events := testEvents(numEvents)
			testInstances := [][]workoutlog.ExerciseInstance{}
			testDates := []int64{}
			for _, evt := range events {
				testDates = append(testDates, evt.Date)
				testInstances = append(testInstances, []workoutlog.ExerciseInstance{evt.Exercises[1]})
			}
			mockExerciseManager.On("GetExerciseType", mock.Anything, mock.Anything).Return(&testExType, nil)

			mockEventManager.On("GetPageOfInstances", mock.Anything, mock.Anything, mock.Anything).Return(testDates, testInstances, nil)

			//  /api/events/metrics?type=blah&startdate=blah&enddate=blah
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%smetrics", url), nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			EventsApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			returnedMetrics := metrics{}
			err := json.NewDecoder(w.Result().Body).Decode(&returnedMetrics)

			So(err, ShouldBeNil)
			So(len(returnedMetrics.Dates), ShouldEqual, numEvents)
			So(len(returnedMetrics.Volume), ShouldEqual, numEvents)
			So(len(returnedMetrics.Load), ShouldEqual, numEvents)
			for i := 0; i < numEvents; i++ {
				So(returnedMetrics.Dates[i], ShouldEqual, testTime+int64(i))
				So(returnedMetrics.Volume[i], ShouldEqual, numSets*numReps)
				So(returnedMetrics.Load[i], ShouldEqual, numSets*numReps*testIntensityValue)
			}
		})
	})
}
