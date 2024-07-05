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

var testExerciseID = "test-exercise-id"

const (
	testExerciseName     = "testExerciseName"
	testVolume           = "count"
	testIntensity        = "weight"
	testVolumeConstraint = 2
)

var testComposition = map[string]int{"id1": 2, "id2": 1}

func testExerciseTypeNonComposite() workoutlog.ExerciseType {
	return workoutlog.ExerciseType{
		Name:             testExerciseName,
		ID:               testExerciseID,
		IntensityType:    testIntensity,
		VolumeType:       testVolume,
		VolumeConstraint: testVolumeConstraint,
	}
}

func testExerciseType() workoutlog.ExerciseType {
	return workoutlog.ExerciseType{
		Name:             testExerciseName,
		ID:               testExerciseID,
		IntensityType:    testIntensity,
		VolumeType:       testVolume,
		VolumeConstraint: testVolumeConstraint,
		Composition:      testComposition,
	}
}

func TestHandleExercises(t *testing.T) {
	baseURL := "/homegym/api/exercises/"

	Convey("Given an exercise manager", t, func() {

		mockEmgr := workoutlog.NewMockExerciseManager()
		workoutlog.ExerciseManager = mockEmgr

		Convey("When we receive a request to create a new exercise type", func() {
			mockEmgr.On("NewExerciseType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&testExerciseID, nil)

			jsonStr, err := json.Marshal(testExerciseType())
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ExerciseTypesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			id := new(returnedID)

			if err := json.NewDecoder(w.Result().Body).Decode(id); err != nil {
				t.Fail()
			}
			So(id, ShouldResemble, id)
		})

		Convey("When we receive a request to create a new non-composite exercise type", func() {
			mockEmgr.On("NewExerciseType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&testExerciseID, nil)

			jsonStr, err := json.Marshal(testExerciseTypeNonComposite())
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ExerciseTypesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			id := new(returnedID)

			if err := json.NewDecoder(w.Result().Body).Decode(id); err != nil {
				t.Fail()
			}
			So(id, ShouldResemble, id)
		})

		Convey("When we receive a request to create a new exercise type of a non-unique name", func() {
			mockEmgr.On("NewExerciseType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, workoutlog.ErrNameNotUnique)

			jsonStr, err := json.Marshal(testExerciseTypeNonComposite())
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ExerciseTypesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("When we receive a request for all exercise types", func() {
			mockEmgr.On("GetExerciseTypes", mock.Anything).Return([]workoutlog.ExerciseType{testExerciseTypeNonComposite(), testExerciseTypeNonComposite(), testExerciseTypeNonComposite()}, nil)

			req := httptest.NewRequest(http.MethodGet, baseURL, nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ExerciseTypesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			types := []workoutlog.ExerciseType{}

			if err := json.NewDecoder(w.Result().Body).Decode(&types); err != nil {
				t.Fail()
			}
			So(len(types), ShouldEqual, 3)
			for _, v := range types {
				So(v, ShouldResemble, testExerciseTypeNonComposite())
			}
		})

		Convey("When we receive a request to update an exercise type", func() {
			mockEmgr.On("UpdateExerciseType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			jsonStr, err := json.Marshal(testExerciseType())
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			url := fmt.Sprintf("%s%s", baseURL, testExerciseID)
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ExerciseTypesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusNoContent)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")
		})

		Convey("When we receive a request to update a non-composite exercise type", func() {
			mockEmgr.On("UpdateExerciseType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			jsonStr, err := json.Marshal(testExerciseTypeNonComposite())
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			url := fmt.Sprintf("%s%s", baseURL, testExerciseID)
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ExerciseTypesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusNoContent)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")
		})

		Convey("When we receive a request to update an exercise type where the ID does not match the url parameter", func() {
			jsonStr, err := json.Marshal(testExerciseTypeNonComposite())
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			url := fmt.Sprintf("%s%s", baseURL, "badID")
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ExerciseTypesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("When we receive a request to update an exercise with an invalid definition", func() {
			badType := testExerciseTypeNonComposite()
			badType.IntensityType = "badType"
			jsonStr, err := json.Marshal(badType)
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			url := fmt.Sprintf("%s%s", baseURL, "badID")
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ExerciseTypesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("When updating an exercise type returns an error", func() {
			mockEmgr.On("UpdateExerciseType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("an error"))

			jsonStr, err := json.Marshal(testExerciseTypeNonComposite())
			if err != nil {
				t.Errorf("failed to marshal test activity: %s", err.Error())
			}

			url := fmt.Sprintf("%s%s", baseURL, testExerciseID)
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ExerciseTypesApi(w, req)

			So(w.Result().StatusCode, ShouldNotEqual, http.StatusNoContent)
			So(w.Result().StatusCode, ShouldNotEqual, http.StatusOK)
		})
	})
}
