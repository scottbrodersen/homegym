package workoutlog

import (
	"errors"
	"fmt"
	"testing"

	"github.com/scottbrodersen/homegym/dal"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

var (
	testUserID       = "testUserID"
	testActivityID   = "testActivityID"
	testActivityName = "testActivityName"
)

func TestActivities(t *testing.T) {
	Convey("Given a dal client", t, func() {
		db := dal.NewMockDal()
		dal.DB = db
		Convey("When we add an activity", func() {
			db.On("GetActivityNames", mock.Anything).Return(map[string]string{}, nil)
			db.On("AddActivity", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			activity, err := ActivityManager.NewActivity(testUserID, testActivityName)
			So(err, ShouldBeNil)
			So(activity.ID, ShouldNotBeEmpty)
			So(activity.Name, ShouldEqual, testActivityName)
		})

		Convey("When we add an activity with a name that is already used", func() {
			//a := Activity{ID: testActivityID, Name: testActivityName}
			db.On("GetActivityNames", mock.Anything).Return(map[string]string{testActivityID: testActivityName}, nil)
			activity, err := ActivityManager.NewActivity(testUserID, testActivityName)
			So(err, ShouldNotBeNil)
			So(errors.Is(err, ErrActivityNameTaken), ShouldBeTrue)
			So(activity, ShouldBeNil)
		})

		Convey("When we attempt to rename an activity with an already-used name", func() {
			db.On("GetActivityNames", mock.Anything).Return(map[string]string{"some-id": "any name"}, nil)
			u := Activity{ID: "some-id", Name: "any name"}
			err := ActivityManager.RenameActivity(testUserID, u)
			So(err, ShouldNotBeNil)
		})

		Convey("When we attempt to rename an activity that has not been added", func() {
			db.On("GetActivityNames", mock.Anything).Return(map[string]string{}, nil)
			db.On("UpdateActivity", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf(""))
			u := Activity{ID: testActivityID, Name: "any name"}
			err := ActivityManager.RenameActivity(testUserID, u)
			So(err, ShouldNotBeNil)
		})

		Convey("When we rename an activity", func() {
			db.On("GetActivityNames", mock.Anything).Return(map[string]string{}, nil)
			db.On("UpdateActivity", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			activity := Activity{ID: testActivityID, Name: "newname"}
			err := ActivityManager.RenameActivity(testUserID, activity)
			So(err, ShouldBeNil)
		})

		Convey("When we add an exercise to an activity", func() {
			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testActivityName, []string{}, nil)
			db.On("AddExerciseToActivity", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			db.On("GetExercise", mock.Anything, mock.Anything).Return([]byte("not nil"), nil)
			activity := Activity{ID: testActivityID}
			err := activity.AddExerciseToActivity(testUserID, testExerciseID)
			So(err, ShouldBeNil)
			So(activity.ExerciseIDs, ShouldNotBeEmpty)
			So(len(activity.ExerciseIDs), ShouldEqual, 1)
			So(activity.ExerciseIDs[0], ShouldEqual, testExerciseID)
		})

		Convey("When we get activities for a user", func() {
			db.On("GetActivityNames", mock.Anything).Return(map[string]string{"id1": "name1", "id2": "name2"}, nil)
			activites, err := ActivityManager.GetActivityNames(testUserName)
			So(err, ShouldBeNil)
			So(activites, ShouldNotBeNil)
			So(len(activites), ShouldEqual, 2)
			for _, a := range activites {
				if a.ID == "id1" {
					So(a.Name, ShouldEqual, "name1")
				} else if a.ID == "id2" {
					So(a.Name, ShouldEqual, "name2")
				}
			}
		})

		Convey("When we get exerciseIDs for an activity", func() {
			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testActivityName, []string{"testid1", "testid2"}, nil)
			testActivity := Activity{ID: testActivityID, Name: testActivityName}
			err := testActivity.GetActivityExercises(testUserID)
			So(err, ShouldBeNil)
			So(testActivity.ID, ShouldEqual, testActivityID)
			So(testActivity.Name, ShouldEqual, testActivityName)
			So(testActivity.ExerciseIDs, ShouldResemble, []string{"testid1", "testid2"})
		})

		Convey("When we remove an activity's exercise", func() {
			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testActivityName, []string{"testid1", "testid2"}, nil)
			db.On("UpdateActivityExercises", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			testActivity := Activity{ID: testActivityID, ExerciseIDs: []string{"testid2"}}

			err := ActivityManager.UpdateActivityExercises(testUserID, testActivity)

			So(err, ShouldBeNil)
			db.AssertCalled(t, "UpdateActivityExercises", testUserID, testActivityID, []string{}, []string{"testid1"})
		})

		Convey("When we add and remove exercises from an activity", func() {

			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testActivityName, []string{"testid1", "testid2"}, nil)
			db.On("UpdateActivityExercises", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			testActivity := Activity{ID: testActivityID, ExerciseIDs: []string{"testid3", "testid1"}}

			err := ActivityManager.UpdateActivityExercises(testUserID, testActivity)

			So(err, ShouldBeNil)
			db.AssertCalled(t, "ReadActivity", testUserID, testActivityID)
			db.AssertCalled(t, "UpdateActivityExercises", testUserID, testActivityID, []string{"testid3"}, []string{"testid2"})
		})
	})
}
