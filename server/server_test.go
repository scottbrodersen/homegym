package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/scottbrodersen/homegym/auth"
	"github.com/scottbrodersen/homegym/programs"
	"github.com/scottbrodersen/homegym/workoutlog"
	"github.com/stretchr/testify/mock"
)

// const (
// 	testpath string = "./temp"
// )

// func TestServer(t *testing.T) {

// 	defer cleanup()
// 	db, err := dal.InitClient(testpath)
// 	if err != nil {
// 		log.Fatal("cannot create database client")
// 	}
// 	dal.DB = db

// 	types, err := workoutlog.ExerciseManager.GetExerciseTypes("test")

// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println(types)

// }

// 	if err = auth.InitiateKeyRotation(auth.KeyTypes.Token); err != nil {
// 		log.Fatal(err)
// 	}
// 	_, err = workoutlog.FrontDesk.NewUser(testUserName, auth.Admin, testEmail, testPassword)
// 	if err != nil && err.Error() != "value not unique" {
// 		log.Fatal(err)
// 	}
// 	StartUnsafe(cleanupShutDown)
// }

func samesiteString() string {
	switch samesite {
	case http.SameSiteLaxMode:
		return "Lax"
	case http.SameSiteStrictMode:
		return "Strict"
	case http.SameSiteNoneMode:
		return "None"
	default:
		return ""
	}
}

// func cleanup() {
// 	_ = recover()
// 	dberr := os.RemoveAll(testpath)
// 	if dberr != nil {
// 		log.Printf("did not delete test db: %s", dberr.Error())
// 	}
// }
// func cleanupShutDown(err error) {
// 	cleanup()
// 	DefaultShutdown(err)
// }

func testContext() context.Context {
	c := context.Background()
	return context.WithValue(context.WithValue(c, usernameKey, testUserName), roleKey, string(auth.User))
}

func newMockEventAdmin() *mockEventAdmin {
	return new(mockEventAdmin)
}

type mockEventAdmin struct {
	mock.Mock
}

func (e *mockEventAdmin) GetCachedExerciseType(exerciseTypeID string) *workoutlog.ExerciseType {
	args := e.Called(exerciseTypeID)

	return args.Get(0).(*workoutlog.ExerciseType)

}

func (e *mockEventAdmin) NewEvent(userID string, event workoutlog.Event) (*string, error) {
	args := e.Called(userID, event)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*string), nil
}

func (e *mockEventAdmin) UpdateEvent(userID string, currentDate int64, event workoutlog.Event) error {
	args := e.Called(userID, currentDate, event)

	return args.Error(0)
}

func (e *mockEventAdmin) AddExercisesToEvent(userID, eventID string, eventDate int64, instances []workoutlog.ExerciseInstance) error {
	args := e.Called(userID, eventID, eventDate, instances)

	return args.Error(0)
}

func (e *mockEventAdmin) GetEventExercises(userID, eventID string) (map[int]workoutlog.ExerciseInstance, error) {
	args := e.Called(userID, eventID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[int]workoutlog.ExerciseInstance), nil
}

func (e *mockEventAdmin) GetPageOfEvents(userID string, previousEvent workoutlog.Event, pageSize int) ([]workoutlog.Event, error) {
	args := e.Called(userID, previousEvent, pageSize)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]workoutlog.Event), nil
}

func newMockActivityAdmin() *mockActivityAdmin {
	return new(mockActivityAdmin)
}

type mockActivityAdmin struct {
	mock.Mock
}

func (mma *mockActivityAdmin) RenameActivity(userID string, activity workoutlog.Activity) error {
	args := mma.Called(userID, activity)

	return args.Error(0)
}

func (maa *mockActivityAdmin) NewActivity(userID, name string) (*workoutlog.Activity, error) {
	args := maa.Called(userID, name)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*workoutlog.Activity), nil
}

func (maa *mockActivityAdmin) GetActivityNames(userID string) ([]*workoutlog.Activity, error) {
	args := maa.Called(userID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*workoutlog.Activity), nil
}

func (maa *mockActivityAdmin) UpdateActivityExercises(userID string, updated workoutlog.Activity) error {
	args := maa.Called(userID, updated)

	return args.Error(0)
}

func newMockUserAdmin() *mockUserAdmin {
	return new(mockUserAdmin)
}

type mockUserAdmin struct {
	mock.Mock
}

func (m *mockUserAdmin) NewUser(username string, role auth.Role, email, password string) (*workoutlog.User, error) {
	args := m.Called(username, role, email, password)

	if args.Error(1) != nil {
		return nil, fmt.Errorf("test error")
	}

	return args.Get(0).(*workoutlog.User), nil

}

func (m *mockUserAdmin) GetUser(username string) (*workoutlog.User, error) {

	args := m.Called(username)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*workoutlog.User), nil

}

type MockAuthorizer struct {
	tokenTTL   int
	sessionTTL int
	mock.Mock
}

func NewMockAuthorizer() *MockAuthorizer {
	ma := MockAuthorizer{tokenTTL: 300, sessionTTL: 1440}

	return &ma
}

func (a *MockAuthorizer) TokenTTL() int {
	return a.tokenTTL
}

func (a *MockAuthorizer) SessionTTL() int {
	return a.sessionTTL
}

func (a *MockAuthorizer) IssueToken(username string, pwd string) (*string, *string, error) {
	args := a.Called(username, pwd)

	if args.Error(2) != nil {
		return nil, nil, args.Error(2)
	}

	return args.Get(0).(*string), args.Get(1).(*string), nil

}

func (a *MockAuthorizer) ValidateToken(tokenString, sessionID string) (*string, error) {
	args := a.Called(tokenString, sessionID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*string), nil
}

func (a *MockAuthorizer) TokenClaims(tokenString string) (auth.Claims, error) {
	args := a.Called(tokenString)

	if args.Error(1) != nil {
		return auth.Claims{}, args.Error(1)
	}

	return args.Get(0).(auth.Claims), nil
}

func newMockProgramManager() *MockProgramManager {
	return new(MockProgramManager)
}

type MockProgramManager struct {
	mock.Mock
}

func (mpm *MockProgramManager) AddProgram(userID string, program programs.Program) (*string, error) {
	args := mpm.Called(userID, program)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*string), nil
}

func (mpm *MockProgramManager) UpdateProgram(userID string, program programs.Program) error {
	args := mpm.Called(userID, program)

	if args.Error(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (mpm *MockProgramManager) GetProgramsPageForActivity(userID, activityID, previousProgramID string, pageSize int) ([]programs.Program, error) {
	args := mpm.Called(userID, activityID, previousProgramID, pageSize)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]programs.Program), nil
}

func (mpm *MockProgramManager) AddProgramInstance(userID string, instance *programs.ProgramInstance) error {
	args := mpm.Called(userID, instance)

	if args.Error(0) != nil {
		return args.Error(0)
	}

	instance.ID = testProgramInstanceID

	return nil
}

func (mpm *MockProgramManager) UpdateProgramInstance(userID string, instance programs.ProgramInstance) (*programs.ProgramInstance, error) {
	args := mpm.Called(userID, instance)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*programs.ProgramInstance), nil
}

func (mpm *MockProgramManager) GetProgramInstancesPage(userID, programID, previousProgramInstanceID string, pageSize int) ([]programs.ProgramInstance, error) {
	args := mpm.Called(userID, previousProgramInstanceID, pageSize)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]programs.ProgramInstance), nil
}

func (mpm *MockProgramManager) SetActiveProgramInstance(userID, activityID, programID, instanceID string) error {
	args := mpm.Called(userID, activityID, programID, instanceID)

	if args.Error(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (mpm *MockProgramManager) GetActiveProgramInstance(userID, activityID string) (*programs.ProgramInstance, error) {
	args := mpm.Called(userID, activityID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*programs.ProgramInstance), nil
}
