package dal

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testProgramID             = "test-program-id"
	testProgramInstanceID     = "test-program-instance-id"
	numberOfTestProgramsToAdd = 10
)

var testProgram []byte = []byte("test program")
var testProgramInstance []byte = []byte("test program instance")

func init() {
	cleanup()
}
func TestProgramDal(t *testing.T) {
	defer cleanup()

	Convey("Given a dal client", t, func() {
		db, err := InitClient(testPath)
		if err != nil {
			log.Fatal("failed to create dal client")
		}
		defer db.Destroy()

		Convey("When we add a program", func() {
			err := db.AddProgram(testUserID, testActivityID, testProgramID, testProgram)

			So(err, ShouldBeNil)
		})

		Convey("When we get a non-full page of programs", func() {
			page, err := db.GetProgramPage(testUserID, testActivityID, "", 10)

			So(err, ShouldBeNil)
			So(page, ShouldHaveLength, 1)
			So(page[0], ShouldResemble, testProgram)
		})

		Convey("When we add lots of programs", func() {
			for i := 0; i < numberOfTestProgramsToAdd-1; i++ {
				err := db.AddProgram(testUserID, testActivityID, fmt.Sprintf("%s-%d", testProgramID, i), []byte(fmt.Sprintf("test program instance %d", i)))
				So(err, ShouldBeNil)
			}
		})

		Convey("When we get a page of programs", func() {
			page, err := db.GetProgramPage(testUserID, testActivityID, "", 4)

			So(err, ShouldBeNil)
			So(page, ShouldHaveLength, 4)
			So(page[0], ShouldResemble, testProgram)
			So(page[3], ShouldResemble, []byte("test program instance 2"))
		})

		Convey("When we get a second page of programs", func() {
			page, err := db.GetProgramPage(testUserID, testActivityID, fmt.Sprintf("%s-2", testProgramID), 4)

			So(err, ShouldBeNil)
			So(page, ShouldHaveLength, 4)
			So(page[0], ShouldResemble, []byte("test program instance 3"))
			So(page[3], ShouldResemble, []byte("test program instance 6"))
		})

		Convey("When we add a program instance", func() {
			err := db.AddProgramInstance(testUserID, testProgramID, testProgramInstanceID, testActivityID, testProgramInstance)

			So(err, ShouldBeNil)
		})

		Convey("When we get all programs", func() {
			page, err := db.GetProgramPage(testUserID, testActivityID, "", numberOfTestProgramsToAdd+1)

			So(err, ShouldBeNil)
			So(page, ShouldHaveLength, numberOfTestProgramsToAdd)
		})

		Convey("When we get the program instance", func() {
			instance, err := db.GetProgramInstancePage(testUserID, testProgramID, testProgramInstanceID, 1)

			So(err, ShouldBeNil)
			So(instance, ShouldResemble, [][]byte{testProgramInstance})
		})

		Convey("When we get the active program", func() {
			activeProgram, err := db.GetActiveProgramInstance(testUserID, testActivityID)

			So(err, ShouldBeNil)
			So(activeProgram, ShouldResemble, testProgramInstance)

		})

		Convey("When we set the active program", func() {
			err := db.SetActiveProgramInstance(testUserID, testActivityID, testProgramID, testProgramInstanceID)

			So(err, ShouldBeNil)
		})

		Convey("When we get a program instance using the wrong ID", func() {
			expectNil, err := db.GetProgramInstancePage(testUserID, testProgramID, "bad id", 1)

			So(err, ShouldBeNil)
			So(expectNil, ShouldBeNil)
		})
	})
}
