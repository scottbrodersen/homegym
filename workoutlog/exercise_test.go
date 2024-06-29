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

var testComposition = map[string]int{"id1": 2, "id2": 3}
var exerciseType1 = ExerciseType{ID: "id1"}
var exerciseType2 = ExerciseType{ID: "id2"}
var exerciseType3 = ExerciseType{ID: "id3"}

var testEvents []Event

func testEventMaker(date int64) Event {
	return Event{
		Date:      date,
		Exercises: map[int]ExerciseInstance{1: {TypeID: exerciseType1.ID}, 2: {TypeID: exerciseType2.ID}, 3: {TypeID: exerciseType3.ID}},
	}
}

// func testExerciseInstanceMaker(typeID string) ExerciseInstance {
// 	return ExerciseInstance{TypeID: typeID}
// }

func TestExercises(t *testing.T) {
	Convey("Given a dal client and an exercise manager", t, func() {
		db := dal.NewMockDal()
		dal.DB = db

		ExerciseManager = new(exerciseManager)

		Convey("When we create a non-composite ExerciseType", func() {
			db.On("GetExercises", mock.Anything).Return([][]byte{}, nil)
			db.On("AddExercise", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			id, err := ExerciseManager.NewExerciseType(testUserID, testExerciseName, testIntensity, testVolume, testVolConstraint, nil, "")

			Convey("Then the exercise id is returned", func() {
				So(err, ShouldBeNil)
				So(id, ShouldNotBeEmpty)
			})
		})

		Convey("When we create a composite ExerciseType", func() {

			ex1Json, _ := json.Marshal(exerciseType1)
			ex2Json, _ := json.Marshal(exerciseType2)
			exercises := [][]byte{ex1Json, ex2Json}
			db.On("GetExercises", mock.Anything).Return(exercises, nil)
			db.On("AddExercise", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			id, err := ExerciseManager.NewExerciseType(testUserID, testExerciseName, testIntensity, testVolume, testVolConstraint, testComposition, "")

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
			ed, err := ExerciseManager.NewExerciseType(testUserID, testExerciseName, testIntensity, testVolume, testVolConstraint, nil, "")
			Convey("Then no exercise is created", func() {
				So(err, ShouldNotBeNil)
				So(errors.Is(err, ErrNameNotUnique), ShouldBeTrue)
				So(ed, ShouldBeNil)
			})
		})

		Convey("When we attempt to create an exercise composed of non-existent types", func() {
			db.On("GetExercises", mock.Anything).Return([][]byte{}, nil)
			db.On("AddExercise", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			id, err := ExerciseManager.NewExerciseType(testUserID, testExerciseName, testIntensity, testVolume, testVolConstraint, testComposition, "")

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(id, ShouldBeNil)
			})

		})

	})
}
