package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/scottbrodersen/homegym/dal"
	"github.com/scottbrodersen/homegym/programs"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"testing"
)

const (
	testProgramID                = "test-program-id"
	testProgramTitle             = "program title"
	testBlockTitle               = "block title"
	testMicroCycleTitle          = "microcycle title"
	testSpan                 int = 7
	testWorkoutTitle             = "workout title"
	testIntensityStr             = "test intensity"
	testVolumeStr                = "test volume"
	testProgramInstanceID        = "test-instance-id"
	testProgramInstanceTitle     = "program instance title"
	testEventID                  = "test-event-id"
)

var testStartDate int64 = time.Now().Local().Unix()

func testSegment() programs.WorkoutSegment {
	return programs.WorkoutSegment{
		ExerciseTypeID: testExerciseID,
		Prescription:   testVolumeStr,
	}
}

func testWorkout(index int) programs.Workout {
	return programs.Workout{
		Title: fmt.Sprintf("%s %d", testWorkoutTitle, index),
		Segments: []programs.WorkoutSegment{
			testSegment(),
			testSegment(),
		},
		Description: testIntensityStr,
	}
}

func testMicroCycle(index int) programs.MicroCycle {
	return programs.MicroCycle{
		Title:       fmt.Sprintf("%s %d", testMicroCycleTitle, index),
		Span:        testSpan,
		Description: testIntensityStr,
		Workouts: []programs.Workout{
			testWorkout(0),
			testWorkout(1),
			testWorkout(2),
			testWorkout(3),
			testWorkout(4),
			testWorkout(5),
			testWorkout(6),
		},
	}
}

func testBlock(index int) programs.Block {
	return programs.Block{
		Title: fmt.Sprintf("%s %d", testBlockTitle, index),
		MicroCycles: []programs.MicroCycle{
			testMicroCycle(0),
			testMicroCycle(1),
			testMicroCycle(2),
			testMicroCycle(3),
		},
		Description: samesiteString(),
	}
}

func testProgram() programs.Program {
	return programs.Program{
		ID:         testProgramID,
		Title:      testProgramTitle,
		ActivityID: testActivityID,
		Blocks:     []programs.Block{testBlock(0), testBlock(1)},
	}
}

func testProgramInstance() programs.ProgramInstance {
	inst := programs.ProgramInstance{
		Program:   testProgram(),
		ID:        testProgramInstanceID,
		ProgramID: testProgramID,
		StartTime: testStartDate,
	}

	inst.Program.ID = ""
	instTime := time.UnixMicro(inst.StartTime * 1000)
	inst.Program.Title = fmt.Sprintf("%s - %s", inst.Program.Title, instTime.Format("Jan _2, 2006"))

	return inst
}

func TestHandlePrograms(t *testing.T) {
	url := fmt.Sprintf("/homegym/api/activities/%s/programs", testActivityID)
	Convey("Given a mock dal and program manager", t, func() {
		db := dal.NewMockDal()
		dal.DB = db

		mpm := newMockProgramManager()
		programs.ProgramManager = mpm

		Convey("When we receive a request to add a program", func() {
			mpm.On("AddProgram", mock.Anything, mock.Anything).Return(&testActivityID, nil)
			testProgram := testProgram()
			testProgram.ID = ""
			jsonStr, err := json.Marshal(testProgram)

			if err != nil {
				t.Errorf("failed to marshal test program: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			body := struct{ ID string }{}

			if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
				t.Fail()
			}

			So(body.ID, ShouldNotBeEmpty)
		})

		Convey("When we receive a request to update a program", func() {
			testProgram := testProgram()
			jsonByte, err := json.Marshal(testProgram)

			if err != nil {
				t.Errorf("failed to marshal test program: %s", err.Error())
			}

			mpm.On("UpdateProgram", mock.Anything, mock.Anything).Return(nil)
			db.On("GetProgramPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([][]byte{jsonByte}, nil)
			db.On("AddProgram", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			db.On("ReadActivity", mock.Anything, mock.Anything).Return(&testActivityName, []string{}, nil)

			if err != nil {
				t.Errorf("failed to marshal test program: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", url, testProgramID), bytes.NewBuffer(jsonByte))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
		})

		Convey("When we receive a request for a page of programs", func() {
			pagesize := 5
			page := []programs.Program{}
			for i := 0; i < pagesize; i++ {
				page = append(page, testProgram())
			}

			mpm.On("GetProgramsPageForActivity", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(page, nil)

			requestURL := fmt.Sprintf("%s?previous=%s&size=%d", url, testProgramID, pagesize)
			req := httptest.NewRequest(http.MethodGet, requestURL, nil).WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)

			body := []programs.Program{}

			if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
				t.Fail()
			}

			So(body, ShouldHaveLength, pagesize)
			So(body[0], ShouldResemble, testProgram())
		})

		Convey("When we receive a request to get a program", func() {
			program := []programs.Program{testProgram()}

			mpm.On("GetProgramsPageForActivity", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(program, nil)

			requestURL := fmt.Sprintf("%s/%s", url, testProgramID)
			req := httptest.NewRequest(http.MethodGet, requestURL, nil).WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)

			body := programs.Program{}

			if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
				t.Fail()
			}

			So(body, ShouldResemble, testProgram())
		})

		Convey("When we receive a request to add a program instance", func() {
			piURL := fmt.Sprintf("%s/%s/instances", url, testProgramID)
			id := testProgramInstanceID
			mpm.On("AddProgramInstance", mock.Anything, mock.Anything).Return(nil)
			testInstance := testProgramInstance()
			testInstance.ID = ""
			jsonStr, err := json.Marshal(testInstance)

			if err != nil {
				t.Errorf("failed to marshal test program instance: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, piURL, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			body := new(programs.ProgramInstance)

			if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
				t.Fail()
			}

			So(body.ID, ShouldNotBeEmpty)
			So(body.ID, ShouldEqual, id)
		})

		Convey("When we receive a request to update a program instance", func() {
			piURL := fmt.Sprintf("%s/%s/instances/%s", url, testProgramID, testProgramInstanceID)
			testPI := testProgramInstance()
			mpm.On("UpdateProgramInstance", mock.Anything, mock.Anything).Return(&testPI, nil)
			jsonStr, err := json.Marshal(testPI)

			if err != nil {
				t.Errorf("failed to marshal test program instance: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, piURL, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			body := programs.ProgramInstance{}

			if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
				t.Fail()
			}

			So(body, ShouldResemble, testProgramInstance())
		})

		Convey("When we receive a request to get a program instance", func() {
			piURL := fmt.Sprintf("%s/%s/instances/%s", url, testProgramID, testProgramInstanceID)
			testInstance := testProgramInstance()
			mpm.On("GetProgramInstancesPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]programs.ProgramInstance{testInstance}, nil)
			jsonStr, err := json.Marshal(testInstance)

			if err != nil {
				t.Errorf("failed to marshal test program instance: %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodGet, piURL, bytes.NewBuffer(jsonStr))
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			body := programs.ProgramInstance{}

			if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
				t.Fail()
			}

			So(body, ShouldResemble, testProgramInstance())
		})

		Convey("When we receive a request for a page of program instances", func() {
			pagesize := 5
			page := []programs.ProgramInstance{}
			for i := 0; i < pagesize; i++ {
				page = append(page, testProgramInstance())
			}

			mpm.On("GetProgramInstancesPage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(page, nil)

			piURL := fmt.Sprintf("%s/%s/instances", url, testProgramID)

			requestURL := fmt.Sprintf("%s?previous=%s&size=%d", piURL, testProgramID, pagesize)
			req := httptest.NewRequest(http.MethodGet, requestURL, nil).WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)

			body := []programs.ProgramInstance{}

			if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
				t.Fail()
			}

			So(body, ShouldHaveLength, pagesize)
			So(body[0], ShouldResemble, testProgramInstance())
		})

		Convey("When we receive a request to set the active program instance", func() {
			piURL := fmt.Sprintf("%s/instances/active?programid=%s&instanceid=%s", url, testProgramID, testProgramInstanceID)
			mpm.On("SetActiveProgramInstance", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			req := httptest.NewRequest(http.MethodPost, piURL, nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")
		})

		Convey("When we receive a request to get the active program instance", func() {
			piURL := fmt.Sprintf("%s/instances/active", url)
			instance := testProgramInstance()
			mpm.On("GetActiveProgramInstance", mock.Anything, mock.Anything, mock.Anything).Return(&instance, nil)

			req := httptest.NewRequest(http.MethodGet, piURL, nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			So(w.Result().Header.Get("content-type"), ShouldEqual, "application/json")

			body := programs.ProgramInstance{}

			if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
				t.Fail()
			}

			So(body, ShouldResemble, testProgramInstance())
		})

		Convey("When we receive a request to deactivate the active program instance", func() {
			piURL := fmt.Sprintf("%s/instances/active", url)
			mpm.On("DeactivateProgramInstance", mock.Anything, mock.Anything).Return(nil)

			req := httptest.NewRequest(http.MethodDelete, piURL, nil)
			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			ActivitiesApi(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
		})
	})
}
