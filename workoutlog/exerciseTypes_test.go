package workoutlog

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func testExerciseType() ExerciseType {
	return ExerciseType{
		Name:             testExerciseName,
		ID:               testExerciseID,
		IntensityType:    testIntensity,
		VolumeType:       testVolume,
		VolumeConstraint: testVolConstraint,
	}
}

var testSubmittedVolumeString string = `[ [1, 1, 1, 1], [1, 0, 1, 1] ]`
var testSubmittedVolumeStringBadRep string = `[ [1, 4, 1, 1], [1, 0, 1, 1] ]`

var testSubmittedExInstanceJSON string = `{
	"typeID": "testExerciseID",
	"index": 3,
	"parts": [
		 {
			"intensity": 5.0,
			"volume": ` + testSubmittedVolumeString + `
		}
	]
}`

var testIncomingWithBadIndex string = `{
	"typeID": "testExerciseID",
	"index": -1,
	"parts": [
		 {
			"intensity": 5.0,
			"volume": ` + testSubmittedVolumeString + `
		}
	]
}`

var testIncomingWithBadRep string = `{
	"typeID": "testExerciseID",
	"index": 3,
	"parts": [
		{
			"intensity": 5.0,
			"volume": ` + testSubmittedVolumeStringBadRep + `
		}
	]
}`

var testExerciseInstance ExerciseInstance = ExerciseInstance{
	TypeID: "testExerciseID",
	Index:  3,
	Segments: []ExerciseSegment{
		{
			Intensity: 5.0,
			Volume:    [][]float32{{float32(1), float32(0)}},
		},
	},
}

func TestExerciseTypes(t *testing.T) {
	Convey("Given a table of ExerciseType values", t, func() {
		type testEtype = struct {
			Etype ExerciseType
			Valid bool
		}
		// initialize with invalid types
		tests := []testEtype{
			{
				Etype: func() ExerciseType {
					e := testExerciseType()
					e.IntensityType = "unsupported"
					return e
				}(),
				Valid: false,
			},
			{
				Etype: func() ExerciseType {
					e := testExerciseType()
					e.VolumeType = "unsupported"
					return e
				}(),
				Valid: false,
			},
			{
				Etype: func() ExerciseType {
					e := testExerciseType()
					e.Name = ""
					return e
				}(),
				Valid: false,
			},
			{
				Etype: func() ExerciseType {
					e := testExerciseType()
					e.ID = ""
					return e
				}(),
				Valid: false,
			},
			{
				Etype: func() ExerciseType {
					e := testExerciseType()
					e.IntensityType = ""
					return e
				}(),
				Valid: false,
			},
			{
				Etype: func() ExerciseType {
					e := testExerciseType()
					e.VolumeType = ""
					return e
				}(),
				Valid: false,
			},
			{
				Etype: func() ExerciseType {
					e := testExerciseType()
					e.VolumeConstraint = 0
					return e
				}(),
				Valid: false,
			},
			{
				Etype: func() ExerciseType {
					e := testExerciseType()
					e.VolumeType = "distance"
					e.Composition = testComposition
					return e
				}(),
				Valid: false,
			},
			{
				Etype: func() ExerciseType {
					e := testExerciseType()
					e.Basis = "anyID"
					e.Composition = testComposition
					return e
				}(),
				Valid: false,
			},
		}
		// add a valid type for each intensity type
		for _, iType := range intensityTypes {
			e := testExerciseType()
			e.IntensityType = iType

			tests = append(tests, testEtype{Etype: e, Valid: true})
		}
		// add a valid type for each volume type
		for _, vType := range volumeTypes {
			e := testExerciseType()
			e.VolumeType = vType
			if vType != "count" {
				e.VolumeConstraint = 0
			} else {
				e.VolumeConstraint = 1
			}
			tests = append(tests, testEtype{Etype: e, Valid: true})
		}
		testComposite := testExerciseType()
		testComposite.Composition = testComposition
		tests = append(tests, testEtype{Etype: testComposite, Valid: true})

		testVariation := testExerciseType()
		testVariation.Basis = "anyID"
		tests = append(tests, testEtype{Etype: testVariation, Valid: true})

		Convey("Then the validation is as expected", func() {
			for _, v := range tests {
				So(v.Etype.validate() == nil, ShouldEqual, v.Valid)
			}
		})
	})

	Convey("Given an exercise type", t, func() {
		exType := testExerciseType()
		exerciseTypeCache.Store(exType.ID, exType)

		Convey("When we validate a correct incoming exercise instance", func() {
			exIncoming := ExerciseInstance{}

			err := json.Unmarshal([]byte(testSubmittedExInstanceJSON), &exIncoming)
			if err != nil {
				t.Fatal(err)
			}

			err = exType.validateInstance(&exIncoming)

			So(err, ShouldBeNil)
		})

		Convey("When we ingest an exerise instance with a bad rep value in the volume part", func() {
			exIncoming := ExerciseInstance{}

			if err := json.Unmarshal([]byte(testIncomingWithBadRep), &exIncoming); err != nil {
				t.Fatal(err)
			}

			err := exType.validateInstance(&exIncoming)

			So(err, ShouldNotBeNil)

		})
	})

}
