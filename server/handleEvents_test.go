package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/scottbrodersen/homegym/workoutlog"
)

func TestHandleEvents(t *testing.T) {
	url := "/homegym/api/events/"

	testEventMeta := workoutlog.EventMeta{
		Mood:       1,
		Energy:     2,
		Motivation: 3,
		Overall:    4,
		Notes:      "test note",
	}

	testExType := testExerciseTypeNonComposite()

	testInstance := workoutlog.ExerciseInstance{
		TypeID:   testExType.ID,
		Segments: []workoutlog.ExerciseSegment{},
	}

	testInstance.Segments = append(testInstance.Segments, workoutlog.ExerciseSegment{
		Intensity: 5.0,
		Volume:    [][]float32{{1, 1, 1, 1}},
	})

	testExercises := map[uint64]workoutlog.ExerciseInstance{1: testInstance}

	testEvents := []workoutlog.Event{
		{
			ID:         "test-event-1",
			ActivityID: "test-activity-id",
			Date:       time.Now().Unix(),
			EventMeta:  testEventMeta,
		},
		{
			ID:         "test-event-2",
			ActivityID: "test-activity-id",
			Date:       time.Now().Unix(),
			EventMeta:  testEventMeta,
		},
	}

	Convey("Given an event manager and an exercise manager", t, func() {

		mockEventManager := newMockEventAdmin()
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

			// doest recognize the eventStruct arg as being equal
			//So(mockEventManager.AssertCalled(t, "UpdateEvent", GymCtxValue(testContext(), usernameKey), 1234, *eventStruct), ShouldBeTrue)

			So(mockEventManager.Calls[0].Method, ShouldEqual, "UpdateEvent")
			updateEventCalledArgs := mockEventManager.Calls[0].Arguments
			So(updateEventCalledArgs.Get(0), ShouldEqual, GymContextValue(testContext(), usernameKey))
			So(updateEventCalledArgs.Get(1), ShouldEqual, 1234)
			So(updateEventCalledArgs.Get(2), ShouldResemble, *eventStruct)
		})

		Convey("When we add an exercise to an event", func() {
			mockEventManager.On("AddExercisesToEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockEventManager.On("GetCachedExerciseType", mock.Anything).Return(&testExType, nil)

			exJson, err := json.Marshal(map[uint64]workoutlog.ExerciseInstance{0: testInstance})
			if err != nil {
				t.Fail()
			}

			reqUrl := fmt.Sprintf("%s13456/test-event-id/exercises", url)

			req := httptest.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(exJson))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			EventsApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
		})

		Convey("When we get event exercises", func() {
			mockEventManager.On("GetEventExercises", mock.Anything, mock.Anything, mock.Anything).Return(testExercises, nil)

			reqUrl := fmt.Sprintf("%s13456/test-event-id/exercises", url)

			req := httptest.NewRequest(http.MethodGet, reqUrl, nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			EventsApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)

			returnedExercises := new(map[uint64]workoutlog.ExerciseInstance)
			mockEventManager.On("GetCachedExerciseType", mock.Anything).Return(&testExType, nil)
			So(json.NewDecoder(w.Result().Body).Decode(returnedExercises), ShouldBeNil)
			So(*returnedExercises, ShouldResemble, testExercises)
		})

		Convey("When we get a page of events", func() {
			mockEventManager.On("GetPageOfEvents", mock.Anything, mock.Anything, mock.Anything).Return(testEvents, nil)
			mockExerciseManager.On("GetExerciseType", mock.Anything, mock.Anything).Return(&testExType, nil)
			mockEventManager.On("GetCachedExerciseType", mock.Anything).Return(&testExType, nil).Times(len(testEvents))

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

			So(returnedEvents, ShouldResemble, testEvents)
		})
	})
}
