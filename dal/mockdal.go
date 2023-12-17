package dal

import (
	"github.com/stretchr/testify/mock"
)

type MockDal struct {
	mock.Mock
}

func NewMockDal() *MockDal {
	return new(MockDal)
}

func (d *MockDal) NewUser(id, email, pwdHash, pwdHashVersion, role string) error {
	args := d.Called(id, email, pwdHash, pwdHashVersion, role)
	return args.Error(0)
}

func (d *MockDal) AddActivity(userID, activityID, activityName string) error {
	args := d.Called(userID, activityName)
	if args.Error(0) != nil {
		return args.Error(1)
	}
	return nil
}

func (d *MockDal) ReadActivity(userID, activityID string) (*string, []string, error) {
	args := d.Called(userID, activityID)
	if args.Error(2) != nil {
		return nil, nil, args.Error(2)
	}

	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil, nil
	}

	return args.Get(0).(*string), args.Get(1).([]string), nil
}

func (d *MockDal) GetActivityNames(userID string) (map[string]string, error) {
	args := d.Called(userID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), nil
}

func (d *MockDal) UpdateActivity(userID, activityID, activityName string) error {
	args := d.Called(userID, activityID, activityName)
	return args.Error(0)
}

func (d *MockDal) UpdateActivityExercises(userID, activityID string, exIDsToAdd, exIDsToDelete []string) error {
	args := d.Called(userID, activityID, exIDsToAdd, exIDsToDelete)
	return args.Error(0)
}

func (d *MockDal) ReadUser(id string) (*string, *string, *string, *string, error) {
	args := d.Called(id)
	if args.Error(4) != nil {
		return nil, nil, nil, nil, args.Error(4)
	}

	return args.Get(0).(*string), args.Get(1).(*string), args.Get(2).(*string), args.Get(3).(*string), nil
}

func (d *MockDal) UpdateUserProfile(id, email string) error {
	args := d.Called(id, email)

	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}
func (d *MockDal) UpdateUserPassword(id, pwdHash, hashVersion string) error {
	args := d.Called(id, pwdHash, hashVersion)

	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}
func (d *MockDal) ChangeUserRole(id, role string) error {
	args := d.Called(id, role)

	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}
func (d *MockDal) UpdatePwdVersion(userID, version string) error {
	args := d.Called(userID, version)

	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}
func (d *MockDal) AddExercise(userID, exerciseID string, exercise []byte) error {
	args := d.Called(userID, exerciseID, exercise)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (d *MockDal) UpdateExercise(userID, exerciseID string, exercise []byte) error {
	args := d.Called(userID, exerciseID, exercise)

	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (d *MockDal) GetExercise(userID, exerciseID string) ([]byte, error) {
	args := d.Called(userID, exerciseID)
	if args.Error(1) != nil {
		return nil, args.Error(4)
	}

	return args.Get(0).([]byte), nil
}

func (d *MockDal) GetExercises(userID string) ([][]byte, error) {
	args := d.Called(userID)

	return args.Get(0).([][]byte), args.Error(1)
}

func (d *MockDal) AddExerciseToActivity(userID, activityID, exerciseID string) error {
	args := d.Called(userID, activityID, exerciseID)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}
func (d *MockDal) AddEvent(userID, eventID, activityID string, date int64, event []byte) error {
	args := d.Called(userID, eventID, activityID, date, event)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (d *MockDal) ShiftEvent(userID, eventID, activityID string, currentDate, newDate int64, event []byte) error {
	args := d.Called(userID, eventID, activityID, currentDate, newDate, event)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (d *MockDal) GetEvent(userID, eventID string, eventDate int64) ([]byte, error) {
	args := d.Called(userID, eventID, eventDate)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), nil
}

func (d *MockDal) GetEventActivity(userID, eventID string, eventDate int64) (*string, *string, []string, error) {
	args := d.Called(userID, eventID, eventDate)

	if args.Error(3) != nil {
		return nil, nil, nil, args.Error(3)
	}
	return args.Get(0).(*string), args.Get(1).(*string), args.Get(2).([]string), nil
}

func (d *MockDal) GetEventPage(userID, previousEventID string, previousDate int64, pageSize int) (
	[][]byte, error) {
	args := d.Called(userID, previousEventID, previousDate, pageSize)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([][]byte), nil

}

func (d *MockDal) AddExercisesToEvent(userID, eventID string, exerciseIDs map[int]string, exerciseInstances map[int][]byte) error {
	args := d.Called(userID, eventID, exerciseIDs, exerciseInstances)

	if args.Error(0) != nil {
		return args.Error(0)
	}

	return nil

}

func (d *MockDal) GetEventExercises(userID, eventID string) ([][]byte, error) {
	args := d.Called(userID, eventID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([][]byte), nil
}

func (d *MockDal) GetKeys(usage string) (map[string][]byte, map[string][]byte, error) {
	args := d.Called(usage)

	if args.Error(2) != nil {
		return map[string][]byte{}, map[string][]byte{}, args.Error(2)
	}

	if args.Get(0) == nil {
		return map[string][]byte{}, map[string][]byte{}, nil
	}

	return args.Get(0).(map[string][]byte), args.Get(1).(map[string][]byte), nil
}

func (d *MockDal) RotateKeys(key []byte, keyID, usage string) error {
	args := d.Called(key, keyID, usage)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (d *MockDal) DeleteKey(keyID, usage string) error {
	args := d.Called(keyID, usage)
	return args.Error(0)
}

func (d *MockDal) AddSession(username, sessionID string, ttl int) error {
	args := d.Called(username, sessionID)
	return args.Error(0)
}

func (d *MockDal) GetSession(sessionID string) (*string, *int64, error) {
	args := d.Called(sessionID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*string), args.Get(1).(*int64), args.Error(2)
}

func (d *MockDal) DeleteSession(sessionID string) error {
	args := d.Called(sessionID)
	return args.Error(0)
}

func (d *MockDal) GetSessionExpiries() (map[string]int64, error) {
	args := d.Called()
	if args.Error(1) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), nil
}

func (d *MockDal) Destroy() {}

func (d *MockDal) AddProgram(userID, activityID, programID string, program []byte) error {
	args := d.Called(userID, activityID, programID, program)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (d *MockDal) GetProgramPage(userID, activityID, previousProgramID string, pageSize int) ([][]byte, error) {
	args := d.Called(userID, activityID, previousProgramID, pageSize)
	if args.Error(1) != nil {
		return nil, args.Error(0)
	}
	return args.Get(0).([][]byte), nil
}

func (d *MockDal) AddProgramInstance(userID, activityID, programID, instanceID string, instance []byte) error {
	args := d.Called(userID, activityID, programID, instanceID, instance)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}
func (d *MockDal) GetProgramInstancePage(userID, activityID, programID, instanceID string, pageSize int) ([][]byte, error) {
	args := d.Called(userID, activityID, programID, instanceID)
	if args.Error(1) != nil {
		return nil, args.Error(0)
	}
	return args.Get(0).([][]byte), nil
}

func (d *MockDal) SetActiveProgramInstance(userID, activityID, programID, instanceID string) error {
	args := d.Called(userID, activityID, programID, instanceID)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (d *MockDal) GetActiveProgramInstance(userID, activityID, programID string) ([]byte, error) {
	args := d.Called(userID, activityID)
	if args.Error(1) != nil {
		return nil, args.Error(0)
	}
	return args.Get(0).([]byte), nil
}
