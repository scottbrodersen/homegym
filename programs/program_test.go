package programs

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/scottbrodersen/homegym/dal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

const (
	testProgramID            = "test-program-id"
	testProgramTitle         = "test program title"
	testActivityID           = "test-activity-id"
	testUserID               = "test-user-id"
	testProgramInstanceID    = "test-program-instance-id"
	testProgramInstanceTitle = "test program instance title"
)

var testProgram Program = Program{
	ID:         testProgramID,
	Title:      testProgramTitle,
	ActivityID: testActivityID,
}

var testProgramInstance ProgramInstance = ProgramInstance{
	Program: Program{
		Title:      testProgramInstanceTitle,
		ActivityID: testProgram.ActivityID,
	},
	ID:        testProgramInstanceID,
	ProgramID: testProgramID,
}

func TestPrograms(t *testing.T) {

	Convey("Given a dal client", t, func() {
		db := dal.NewMockDal()
		dal.DB = db

		Convey("When we add a program", func() {
			db.On("AddProgram", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			testActivityID := "test ID"
			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testActivityID, []string{"a name"}, nil)
			program := Program{
				Title:      testProgramTitle,
				ActivityID: testActivityID,
			}

			id, err := ProgramManager.AddProgram(testUserID, program)

			So(err, ShouldBeNil)
			So(id, ShouldNotBeEmpty)
		})

		Convey("When we get the program", func() {
			testProgramByte, err := json.Marshal(testProgram)
			if err != nil {
				t.FailNow()
			}
			db.On("GetProgramPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{testProgramByte}, nil)

			programs, err := ProgramManager.GetProgramsPageForActivity(testUserID, testActivityID, testProgramID, 1)

			So(err, ShouldBeNil)
			So(programs, ShouldHaveLength, 1)
			So(programs[0], ShouldResemble, testProgram)
		})

		Convey("When we get the first page of programs", func() {
			testProgramByte, err := json.Marshal(testProgram)
			if err != nil {
				t.FailNow()
			}
			db.On("GetProgramPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{testProgramByte}, nil)

			programs, err := ProgramManager.GetProgramsPageForActivity(testUserID, testActivityID, "", 1)

			So(err, ShouldBeNil)
			So(programs, ShouldHaveLength, 1)
			So(programs[0], ShouldResemble, testProgram)
		})

		Convey("When we update a program", func() {
			db.On("GetProgramPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{[]byte(testActivityID)}, nil)
			db.On("AddProgram", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testProgram.ActivityID, []string{}, nil)

			program := Program{
				ID:         testProgramID,
				Title:      testProgramTitle,
				ActivityID: testActivityID,
			}

			err := ProgramManager.UpdateProgram(testUserID, program)

			So(err, ShouldBeNil)
		})

		Convey("When we attempt to add an invalid program", func() {
			program := Program{
				ActivityID: testActivityID,
			}
			id, err := ProgramManager.AddProgram(testUserID, program)

			So(err, ShouldNotBeNil)
			So(errors.As(err, new(ErrInvalidProgram)), ShouldBeTrue)
			So(id, ShouldBeNil)
		})

		Convey("When we add a program instance", func() {
			db.On("AddProgramInstance", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testProgram.ActivityID, []string{}, nil)
			db.On("GetProgramPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{[]byte(testActivityID)}, nil)

			programInstance := ProgramInstance{
				Program: Program{
					Title:      testProgramInstanceTitle,
					ActivityID: testActivityID,
				},
				ID:        "",
				ProgramID: testProgramID,
				StartTime: time.Now().Local().Unix(),
			}

			err := ProgramManager.AddProgramInstance(testUserID, &programInstance)

			So(err, ShouldBeNil)
			So(programInstance.ID, ShouldNotBeEmpty)
		})

		Convey("When we get the program instance", func() {
			testInstanceByte, err := json.Marshal(testProgramInstance)
			if err != nil {
				t.FailNow()
			}
			db.On("GetProgramInstancePage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{testInstanceByte}, nil)

			instances, err := ProgramManager.GetProgramInstancesPage(testUserID, testProgramID, testProgramInstanceID, 1)

			So(err, ShouldBeNil)
			So(instances, ShouldHaveLength, 1)
			So(instances[0], ShouldResemble, testProgramInstance)
		})

		Convey("When we update a program instance", func() {
			db.On("GetProgramInstancePage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{[]byte("any")}, nil)
			db.On("AddProgramInstance", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			err := ProgramManager.UpdateProgramInstance(testUserID, testProgramInstance)

			So(err, ShouldBeNil)
		})

		Convey("When we set the active program instance", func() {
			db.On("SetActiveProgramInstance", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			err := ProgramManager.SetActiveProgramInstance(testUserID, testProgramInstance.ActivityID, testProgramInstance.ActivityID, testProgramInstance.ID)

			So(err, ShouldBeNil)
		})

		Convey("When we get the active program instance", func() {
			instanceBytes, err := json.Marshal(testProgramInstance)
			if err != nil {
				t.FailNow()
			}
			db.On("GetActiveProgramInstance", mock.Anything, mock.Anything, mock.Anything).Return(instanceBytes, nil)

			inst, err := ProgramManager.GetActiveProgramInstance(testUserID, testActivityID)

			So(err, ShouldBeNil)
			So(*inst, ShouldResemble, testProgramInstance)
		})
	})
}
