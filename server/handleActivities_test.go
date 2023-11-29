package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scottbrodersen/homegym/dal"
	"github.com/scottbrodersen/homegym/workoutlog"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

var (
	testActivityName = "test activity name"
	testActivityID   = "test-activity-id"
)

var testExerciseIDs []string = []string{"eid1", "eid2"}

func testActivity() *workoutlog.Activity {
	return &workoutlog.Activity{
		ID:          testActivityID,
		Name:        testActivityName,
		ExerciseIDs: testExerciseIDs,
	}
}

func TestHandleActivities(t *testing.T) {
	url := "/homegym/api/activities/"

	Convey("Given an activity manager, exercise manager, and db client", t, func() {
		mockActivityManager := newMockActivityAdmin()
		workoutlog.ActivityManager = mockActivityManager

		db := dal.NewMockDal()
		dal.DB = db

		Convey("When we receive a request to create a new activity", func() {
			mockActivityManager.On("NewActivity", mock.Anything, mock.Anything).Return(testActivity(), nil)

			jsonStr, err := json.Marshal(testActivity())
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			returnedActivity := workoutlog.Activity{}

			if err := json.NewDecoder(w.Result().Body).Decode(&returnedActivity); err != nil {
				t.Fail()
			}

			So(returnedActivity, ShouldResemble, *testActivity())
		})

		Convey("When we receive a request to update an activity", func() {
			mockActivityManager.On("RenameActivity", mock.Anything, mock.Anything).Return(nil)

			jsonStr, err := json.Marshal(testActivity())
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", url, testActivityID), bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
		})

		Convey("When we receive a request to list activities", func() {
			testActivities := []*workoutlog.Activity{testActivity(), testActivity(), testActivity()}
			mockActivityManager.On("GetActivityNames", mock.Anything).Return(testActivities, nil)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			returnedAcivities := []workoutlog.Activity{}

			if err := json.NewDecoder(w.Result().Body).Decode(&returnedAcivities); err != nil {
				t.Fail()
			}

			for i, a := range testActivities {
				So(returnedAcivities[i], ShouldResemble, *a)
			}
		})

		Convey("When we receive a request for the exercise types of an activity", func() {

			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testActivityName, testExerciseIDs, nil)

			endpointUrl := fmt.Sprintf("%s%s/exercises", url, testActivityID)
			req := httptest.NewRequest(http.MethodGet, endpointUrl, nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()
			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			returned := []string{}

			if err := json.NewDecoder(w.Result().Body).Decode(&returned); err != nil {
				t.Fail()
			}

			So(returned, ShouldResemble, testExerciseIDs)
		})

		Convey("When we receive a request to update an activity's exercise types", func() {
			mockActivityManager.On("UpdateActivityExercises", mock.Anything, mock.Anything).Return(nil)
			updated := workoutlog.Activity{
				ID:          testActivityID,
				Name:        testActivityName,
				ExerciseIDs: []string{"eid1", "eid2"},
			}
			jsonStr, err := json.Marshal(updated)
			if err != nil {
				t.Errorf("failed to marshal test exercise ids: %s", err.Error())
			}

			endpointUrl := fmt.Sprintf("%s%s/exercises/", url, updated.ID)
			req := httptest.NewRequest(http.MethodPost, endpointUrl, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()
			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
		})
	})
}
