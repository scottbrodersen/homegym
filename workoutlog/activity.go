package workoutlog

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/scottbrodersen/homegym/dal"
)

var ErrActivityNameTaken error = errors.New("activity name not unique")
var ErrNotFound error = errors.New("not found")
var ActivityManager ActivityAdmin = &ActivityMaker{}

type ActivityAdmin interface {
	NewActivity(userID, name string) (*Activity, error)
	GetActivityNames(userID string) ([]*Activity, error)
	RenameActivity(userID string, a Activity) error
	UpdateActivityExercises(userID string, updated Activity) error
}

type ActivityMaker struct{}

type Activity struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	ExerciseIDs []string `json:"exercises"`
}

func (am *ActivityMaker) NewActivity(userID, name string) (*Activity, error) {
	activityNames, err := am.GetActivityNames(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to add activity:%w", err)
	}

	for _, v := range activityNames {
		if v.Name == name {
			return nil, ErrActivityNameTaken
		}
	}

	activityID := uuid.NewString()
	err = dal.DB.AddActivity(userID, activityID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to add activity: %w", err)
	}

	activity := Activity{ID: activityID, Name: name}
	return &activity, nil
}

func (am *ActivityMaker) GetActivityNames(userID string) ([]*Activity, error) {
	names, err := dal.DB.GetActivityNames(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read activities: %w", err)
	}

	activities := []*Activity{}

	for k, v := range names {
		a := Activity{
			ID:   k,
			Name: v,
		}
		activities = append(activities, &a)
	}

	return activities, nil
}

func (a *Activity) GetActivityExercises(userID string) error {
	activity, err := getActivity(userID, a.ID)
	if err != nil {
		return fmt.Errorf("failed to read activity: %w", err)
	}

	a.ExerciseIDs = activity.ExerciseIDs

	return nil
}

func getActivity(userID, activityID string) (*Activity, error) {
	aName, exerciseIDs, err := dal.DB.ReadActivity(userID, activityID)

	if err != nil {
		return nil, fmt.Errorf("failed to read activity: %w", err)
	}

	if aName == nil && exerciseIDs == nil {
		return nil, ErrNotFound
	}

	if len(exerciseIDs) == 0 {
		exerciseIDs = []string{}
	}

	a := Activity{
		ID:          activityID,
		Name:        *aName,
		ExerciseIDs: exerciseIDs,
	}

	return &a, nil
}

func (am *ActivityMaker) RenameActivity(userID string, a Activity) error {
	// check that the new name is not already used
	activityNames, err := am.GetActivityNames(userID)
	if err != nil {
		return fmt.Errorf("failed to rename activity:%w", err)
	}

	for _, v := range activityNames {
		if v.Name == a.Name {
			return ErrActivityNameTaken
		}
	}

	if err := dal.DB.UpdateActivity(userID, a.ID, a.Name); err != nil {
		return (fmt.Errorf("failed to rename activity: %w", err))
	}

	return nil
}

func (a *Activity) AddExerciseToActivity(userID, exerciseID string) error {
	if userID == "" {
		return fmt.Errorf("userID must not be empty")
	}

	if exerciseID == "" {
		return fmt.Errorf("exerciseID must not be empty")
	}

	activity, err := getActivity(userID, a.ID)
	if err != nil {
		return fmt.Errorf("failed to rename activity: %w", err)
	}

	for _, v := range activity.ExerciseIDs {
		if v == exerciseID {
			return nil
		}
	}

	// check that the exercise type exists
	eType, err := dal.DB.GetExercise(userID, exerciseID)

	if err != nil {
		return fmt.Errorf("failure: %w", err)
	}

	if eType == nil {
		return fmt.Errorf("no exercise with that ID was found")
	}

	if err := dal.DB.AddExerciseToActivity(userID, a.ID, exerciseID); err != nil {
		return fmt.Errorf("exercise not added: %w", err)
	}

	a.ExerciseIDs = append(a.ExerciseIDs, exerciseID)

	return nil
}

func (am *ActivityMaker) UpdateActivityExercises(userID string, updated Activity) error {
	activity, err := getActivity(userID, updated.ID)
	if err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
	}

	deletes := []string{}
	adds := []string{}

	for _, v := range activity.ExerciseIDs {
		found := false
		for _, v2 := range updated.ExerciseIDs {
			if v == v2 {
				found = true
				break
			}
		}
		if !found {
			deletes = append(deletes, v)
		}
	}

	for _, v := range updated.ExerciseIDs {
		found := false
		for _, v2 := range activity.ExerciseIDs {
			if v == v2 {
				found = true
				break
			}
		}
		if !found {
			adds = append(adds, v)
		}
	}

	if len(deletes) == 0 && len(adds) == 0 {
		return nil
	}

	if err := dal.DB.UpdateActivityExercises(userID, activity.ID, adds, deletes); err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
	}

	return nil
}
