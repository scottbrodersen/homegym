package programs

import (
	"encoding/json"
	"errors"
	"fmt"
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

func testProgramInstance() ProgramInstance {
	return ProgramInstance{
		Program: Program{
			Title:      testProgramInstanceTitle,
			ActivityID: testProgram.ActivityID,
		},
		ID:        testProgramInstanceID,
		ProgramID: testProgramID,
	}
}

func goodEvent() map[int]string {
	return map[int]string{
		0: "",
		1: "",
		2: "",
		3: "",
		4: "",
	}
}

func TestPrograms(t *testing.T) {

	Convey("When we sanitize well-formed program events", t, func() {
		testEvents := goodEvent()
		err := sanitizeEvents(testEvents)

		So(err, ShouldBeNil)
		So(testEvents, ShouldResemble, goodEvent())
	})
	Convey("When we sanitize events with missing days", t, func() {
		testEvents := map[int]string{
			0: "",
			1: "",
			4: "",
		}
		err := sanitizeEvents(testEvents)

		So(err, ShouldBeNil)
		So(testEvents, ShouldResemble, goodEvent())
	})
	Convey("When we sanitize events with an early missing day", t, func() {
		testEvents := map[int]string{
			0: "",
			2: "",
			3: "",
			4: "",
		}
		err := sanitizeEvents(testEvents)

		So(err, ShouldNotBeNil)
	})

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
			db.On("ActivateProgramInstance", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
			testInstanceByte, err := json.Marshal(testProgramInstance())
			if err != nil {
				t.FailNow()
			}
			db.On("GetProgramInstancePage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{testInstanceByte}, nil)

			instances, err := ProgramManager.GetProgramInstancesPage(testUserID, testProgramID, testProgramInstanceID, 1)

			So(err, ShouldBeNil)
			So(instances, ShouldHaveLength, 1)
			So(instances[0], ShouldResemble, testProgramInstance())
		})

		Convey("When we update a program instance", func() {
			db.On("GetProgramInstancePage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{[]byte("any")}, nil)
			db.On("AddProgramInstance", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			pi, err := ProgramManager.UpdateProgramInstance(testUserID, testProgramInstance())

			So(err, ShouldBeNil)
			So(*pi, ShouldResemble, (testProgramInstance()))
		})

		Convey("When we update a program instance with non-contiguous events", func() {
			db.On("GetProgramInstancePage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{[]byte("any")}, nil)
			db.On("AddProgramInstance", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			testPI := testProgramInstance()
			testPI.Events = map[int]string{0: "", 1: "", 3: ""}
			expectedProgramInstance := testProgramInstance()
			expectedProgramInstance.Events = map[int]string{0: "", 1: "", 2: "", 3: ""}

			pi, err := ProgramManager.UpdateProgramInstance(testUserID, testPI)

			So(err, ShouldBeNil)
			So(*pi, ShouldResemble, expectedProgramInstance)
		})

		Convey("When we update a program instance with duplicate events", func() {
			db.On("GetProgramInstancePage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{[]byte("any")}, nil)
			db.On("AddProgramInstance", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			testPI := testProgramInstance()
			testPI.Events = map[int]string{0: "", 1: "", -1: "", 3: ""}

			pi, err := ProgramManager.UpdateProgramInstance(testUserID, testPI)

			So(err, ShouldNotBeNil)
			So(pi, ShouldBeNil)
		})

		Convey("When we activate the program instance", func() {
			db.On("ActivateProgramInstance", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			testPGInstance := testProgramInstance()
			err := ProgramManager.ActivateProgramInstance(testUserID, testPGInstance.ActivityID, testProgramID, testPGInstance.ID)

			So(err, ShouldBeNil)
		})

		Convey("When we get the active program instance", func() {
			instanceIDBytes := [][]byte{([]byte)(fmt.Sprintf("%s:%s", testProgramID, testProgramInstanceID))}
			// testPGInstance := testProgramInstance()

			// // convert the test pg instance to [][]byte
			// jsonPGI, err := json.Marshal(testPGInstance)
			// if err != nil {
			// 	fmt.Printf("Error marshalling test program instance: %v\n", err)
			// 	return
			// }
			db.On("GetActiveProgramInstancePage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(instanceIDBytes, nil)
			//db.On("GetProgramInstancePage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{jsonPGI}, nil)

			inst, err := ProgramManager.GetActiveProgramInstancesPage(testUserID, testActivityID, testProgramInstanceID, 1)

			So(err, ShouldBeNil)
			So(len(inst), ShouldEqual, 1)
			So(inst[0], ShouldEqual, string(fmt.Sprintf("%s:%s", testProgramID, testProgramInstanceID)))
		})

		Convey("When we deactivate the active program instance", func() {

			db.On("DeactivateProgramInstance", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			err := ProgramManager.DeactivateProgramInstance(testUserID, testActivityID, testProgramInstanceID)

			So(err, ShouldBeNil)
		})
	})
}
