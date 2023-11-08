package workoutlog

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/scottbrodersen/homegym/dal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

const (
	testExerciseID    = "testExerciseID"
	testExerciseName  = "testExerciseName"
	testVolume        = "count"
	testIntensity     = "weight"
	testVolConstraint = 2
)

func TestExercises(t *testing.T) {
	Convey("Given a dal client", t, func() {
		db := dal.NewMockDal()
		dal.DB = db
		Convey("When we create an ExerciseType", func() {
			db.On("GetExercises", mock.Anything).Return([][]byte{}, nil)
			db.On("AddExercise", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			id, err := ExerciseManager.NewExerciseType(testUserID, testExerciseName, testIntensity, testVolume, testVolConstraint)

			Convey("Then the exercise id is returned", func() {
				So(err, ShouldBeNil)
				So(id, ShouldNotBeEmpty)
			})
		})

		Convey("When we attempt to create an exercise with a used name", func() {
			ted := testExerciseType()
			exerciseJson, err := json.Marshal(ted)
			if err != nil {
				t.Fail()
			}

			db.On("GetExercises", mock.Anything).Return([][]byte{exerciseJson}, nil)
			ed, err := ExerciseManager.NewExerciseType(testUserID, testExerciseName, testIntensity, testVolume, testVolConstraint)
			Convey("Then no exercise is created", func() {
				So(err, ShouldNotBeNil)
				So(errors.Is(err, ErrNameNotUnique), ShouldBeTrue)
				So(ed, ShouldBeNil)
			})
		})
	})
}
