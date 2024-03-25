package workoutlog

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/scottbrodersen/homegym/dal"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

var ErrNameNotUnique = fmt.Errorf("name is not unique")

var ExerciseManager ExerciseAdmin = new(exerciseManager)

type ExerciseAdmin interface {
	NewExerciseType(userID, name, intensity, volume string, volConstraint int, composition map[string]int, basis string) (*string, error)
	UpdateExerciseType(userID, exerciseID, name, intensity, volume string, volConstraint int, composition map[string]int, basis string) error
	GetExerciseTypes(userID string) ([]ExerciseType, error)
	GetExerciseType(userID, exerciseID string) (*ExerciseType, error)
}

type exerciseManager struct{}

func (ea *exerciseManager) NewExerciseType(userID, name, intensity, volume string, volumeConstraint int, composition map[string]int, basis string) (*string, error) {
	id := uuid.New().String()

	newType := ExerciseType{
		ID:               id,
		Name:             name,
		IntensityType:    intensity,
		VolumeType:       volume,
		VolumeConstraint: volumeConstraint,
		Composition:      composition,
		Basis:            basis,
	}

	if err := newType.validate(); err != nil {
		return nil, fmt.Errorf("cannot add exercise type: %w", err)
	}

	// Check that the name isn't already used
	available, err := isTypeNameAvailable(*ea, userID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check name availability: %w", err)
	}

	if !available {
		return nil, ErrNameNotUnique
	}

	if err := checkDependencies(*ea, userID, newType); err != nil {
		return nil, fmt.Errorf("problem with exercise type: %w", err)
	}

	exerciseJson, err := json.Marshal(newType)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the exercise type: %w", err)
	}

	if err := dal.DB.AddExercise(userID, id, exerciseJson); err != nil {
		return nil, fmt.Errorf("failed to add exercise: %w", err)
	}

	return &id, nil
}

func (ea *exerciseManager) UpdateExerciseType(userID, exerciseID, name, intensity, volume string, volConstraint int, composition map[string]int, basis string) error {
	updated := ExerciseType{
		ID:               exerciseID,
		Name:             name,
		IntensityType:    intensity,
		VolumeType:       volume,
		VolumeConstraint: volConstraint,
		Composition:      composition,
		Basis:            basis,
	}

	if err := updated.validate(); err != nil {
		return fmt.Errorf("invalid exercise type: %w", err)
	}

	// remove from cache so it doesn't become stale
	exerciseTypeCache.Delete(exerciseID)

	exerciseByte, err := dal.DB.GetExercise(userID, exerciseID)
	if err != nil {
		return fmt.Errorf("failed to update exercise: %w", err)
	}

	if exerciseByte == nil {

		return fmt.Errorf("exercise type not found")
	}

	eType := new(ExerciseType)

	err = json.Unmarshal(exerciseByte, eType)
	if err != nil {
		log.Warn("the stored exercise type that we're updating was found to be invalid")
	}

	// Check that the name isn't already used
	if name != eType.Name {
		available, err := isTypeNameAvailable(*ea, userID, name)
		if err != nil {
			return fmt.Errorf("failed to check name availability: %w", err)
		}

		if !available {
			return ErrNameNotUnique
		}
	}

	if err := checkDependencies(*ea, userID, *eType); err != nil {
		return fmt.Errorf("problem with exercise type: %w", err)
	}

	eTypeJSON, err := json.Marshal(updated)
	if err != nil {
		return fmt.Errorf("failed to marshal exercise type: %w", err)
	}

	if err := dal.DB.AddExercise(userID, exerciseID, eTypeJSON); err != nil {
		return fmt.Errorf("failed to update exercise: %w", err)
	}

	return nil
}

func (ea *exerciseManager) GetExerciseTypes(userID string) ([]ExerciseType, error) {
	exercises, err := dal.DB.GetExercises(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}

	exerciseTypes := []ExerciseType{}
	for _, n := range exercises {
		e := ExerciseType{}
		if err := json.Unmarshal(n, &e); err != nil {
			return nil, fmt.Errorf("failed to unmarshal exercise type: %w", err)
		}

		exerciseTypes = append(exerciseTypes, e)
		exerciseTypeCache.Store(e.ID, e)
	}

	return exerciseTypes, nil
}

func (ea *exerciseManager) GetExerciseType(userID, exerciseID string) (*ExerciseType, error) {
	exerciseType := EventManager.GetCachedExerciseType(exerciseID)

	if exerciseType == nil {
		exerciseTypeByte, err := dal.DB.GetExercise(userID, exerciseID)
		if err != nil {
			return nil, fmt.Errorf("failed to get exercise type: %w", err)
		}

		exerciseType = new(ExerciseType)

		if err := json.Unmarshal(exerciseTypeByte, exerciseType); err != nil {
			return nil, fmt.Errorf("failed to unmarshal exercise type: %w", err)
		}
	}

	return exerciseType, nil
}

func isTypeNameAvailable(ea exerciseManager, userID, name string) (bool, error) {
	types, err := ea.GetExerciseTypes(userID)
	if err != nil {
		return false, err
	}

	for _, t := range types {
		if t.Name == name {
			return false, nil
		}
	}

	return true, nil
}

func checkDependencies(em exerciseManager, userID string, eType ExerciseType) error {
	exerciseTypes, err := em.GetExerciseTypes(userID)
	if err != nil {
		return fmt.Errorf("failed to get exercises: %w", err)
	}

	// Check that the composition references existing types
	for id := range eType.Composition {
		found := false
		for _, e := range exerciseTypes {
			if id == e.ID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("composition includes unfound exercise type: %s", id)
		}
	}

	// Check that the basis exists
	if eType.Basis != "" {
		foundBasis := false
		for _, e := range exerciseTypes {
			if e.ID == eType.Basis {
				foundBasis = true
				break
			}
		}
		if !foundBasis {
			return fmt.Errorf("basis references unfound exercise type: %s", eType.Basis)
		}
	}

	return nil
}

func NewMockExerciseManager() *mockExerciseManager {
	return new(mockExerciseManager)
}

type mockExerciseManager struct {
	mock.Mock
}

func (m *mockExerciseManager) NewExerciseType(userID, name, intensity, volume string, volConstraint int, composition map[string]int, basis string) (*string, error) {
	args := m.Called(userID, name, intensity, volume, volConstraint, composition, basis)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*string), nil
}

func (m *mockExerciseManager) UpdateExerciseType(userID, exerciseID, name, intensity, volume string, volConstraint int, composition map[string]int, basis string) error {
	args := m.Called(userID, name, intensity, volume, volConstraint, composition, basis)

	return args.Error(0)
}

func (m *mockExerciseManager) GetExerciseTypes(userID string) ([]ExerciseType, error) {
	args := m.Called(userID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]ExerciseType), nil
}

func (m *mockExerciseManager) GetExerciseType(userID, exerciseID string) (*ExerciseType, error) {
	args := m.Called(userID, exerciseID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ExerciseType), nil
}
