package workoutlog

import (
	"fmt"
	"math"
	"regexp"
	"slices"
)

type ErrInvalidExercise struct {
	Message string
}

func (ec ErrInvalidExercise) Error() string {
	return fmt.Sprintf("invalid exercise: %s", ec.Message)
}

// Defines an exercise
// Factory for ExerciseInstance
// Translates instances between the user interface and the db
// Composition indicates that an exercise is composed of other exercises. Limited to Count types.
type ExerciseType struct {
	Name             string         `json:"name"`
	ID               string         `json:"id"`
	IntensityType    string         `json:"intensityType"`    // for UI
	VolumeType       string         `json:"volumeType"`       // for data type interpretation and validation
	VolumeConstraint int            `json:"volumeConstraint"` // for ui, not generally useful for aerobic activities
	Composition      map[string]int `json:"composition"`      // key is exercise ID, value is number of reps
	Basis            string         `json:"basis"`            // id of exercise of which this is a variation
}

// ExerciseInstance stores data about the performance of an exercise type.
// Index stores the location in the order of performed exercise in an event
type ExerciseInstance struct {
	TypeID   string            `json:"typeID"`
	Index    int               `json:"index"`
	Segments []ExerciseSegment `json:"parts"`
}

// ExerciseSegment stores the performance data for an exercise at a specific intensity.
// Volume is an array of the amount of activity done at the intensity.
// The volume type and constraint of the exercise type dictates the shape of the values in the inner array.
// Time and distance types store a single float32 in the inner array.
// Count types store one or more values of 1 or 0 in the array, depending on the volume constraint
type ExerciseSegment struct {
	Intensity float32     `json:"intensity"`
	Volume    [][]float32 `json:"volume"`
}

// volumeConstraints indicates the type of values that can be expressed for volumes.
// 0 is no restriction (any float32)
// 1 is restricted to the value 1 (for counts)
// 2 is restricted to either 1 or 0 (success/failure)
var volumeConstraints []int = []int{0, 1, 2}
var volumeTypes []string = []string{"count", "distance", "time"}
var intensityTypes []string = []string{"weight", "hrZone", "rpe", "percentOfMax", "bodyweight", "pace"}

func (e ExerciseType) CreateInstance() ExerciseInstance {
	return ExerciseInstance{
		TypeID:   e.ID,
		Segments: []ExerciseSegment{},
	}
}

// validateInstance ensures intensity and volume values are valid for the exercise type
// It also massages data for some types:
//   - scale-based intensity values are stripped of decimals
//   - distance and weight intensity values are truncated to single decimals
//   - bodyweight intensity values are set to 1
//   - non-rep-based volume values are truncated to single decimals
//   - time-based intensities are stripped of decimals
func (et ExerciseType) validateInstance(ei *ExerciseInstance) error {
	if ei.Index < 0 {
		return fmt.Errorf("index must be > 0")
	}

	for i, segment := range ei.Segments {
		// Validate intensity values
		if segment.Intensity <= 0 && et.IntensityType != "bodyweight" {
			return fmt.Errorf("intensity must be greater than zero")
		}

		switch et.IntensityType {

		case "bodyweight":
			ei.Segments[i].Intensity = 1
		case "weight":
			fallthrough
		case "distance":
			ei.Segments[i].Intensity = float32(math.Floor(float64(segment.Intensity*10)) / 10)
		case "hrZone":
			if segment.Intensity > 5 && segment.Intensity < 1 {
				return ErrInvalidExercise{Message: "hrZone must be between 1 and 5"}
			}
			fallthrough
		case "rpe":
			if segment.Intensity > 10 && segment.Intensity < 1 {
				return ErrInvalidExercise{Message: "RPE must be between 1 and 10"}
			}
			fallthrough
		case "pace":
			fallthrough
		case "percentOfMax":
			ei.Segments[i].Intensity = float32(math.Floor(float64(segment.Intensity)))
		}

		// Validate volume values
		for j, set := range segment.Volume {

			if len(set) == 0 {
				return ErrInvalidExercise{Message: "volume is a required value"}
			}
			for k, rep := range set {
				if rep < 0 {
					return ErrInvalidExercise{Message: "volume must be a positive number"}
				}
				switch et.VolumeConstraint {
				case 0:
					ei.Segments[i].Volume[j][k] = float32(math.Floor(float64(rep*10)) / 10)

				case 1:
					// for simple counts, allow 0 so that the value can be interpreted as either tracking failures or not.
					fallthrough
				case 2:
					if rep != 0 && rep != 1 {
						return ErrInvalidExercise{Message: "invalid rep value: must be 0 or 1"}
					}
				}
			}
		}
	}
	return nil
}

// validate ensures the exercise type is correctly defined.
func (e ExerciseType) validate() error {
	if e.Name == "" {
		return ErrInvalidExercise{Message: "name cannot be empty"}
	}

	nameRxp := regexp.MustCompile(`^[a-zA-Z0-9_\-&$*@. \+]+$`)
	if ok := nameRxp.MatchString(e.Name); !ok {
		return ErrInvalidExercise{Message: "name can include letters, digits, _, -, &, $, *, @, ., and spaces"}
	}

	if e.ID == "" {
		return ErrInvalidExercise{Message: "id cannot be empty"}
	}

	if slices.Index(intensityTypes, e.IntensityType) < 0 {
		return ErrInvalidExercise{Message: fmt.Sprintf("invalid intensity type: %s", e.IntensityType)}
	}

	if slices.Index(volumeTypes, e.VolumeType) < 0 {
		return ErrInvalidExercise{Message: fmt.Sprintf("invalid volume type: %s", e.VolumeType)}
	}

	if slices.Index(volumeConstraints, e.VolumeConstraint) < 0 {
		return ErrInvalidExercise{Message: fmt.Sprintf("invalid volume constraint: %d", e.VolumeConstraint)}
	}

	if e.VolumeType == "count" {
		if e.VolumeConstraint != 1 && e.VolumeConstraint != 2 {
			return ErrInvalidExercise{Message: fmt.Sprintf("invalid volume constraint: %d", e.VolumeConstraint)}
		}

	} else {
		if e.VolumeConstraint != 0 {
			return ErrInvalidExercise{Message: fmt.Sprintf("invalid volume constraint: %d", e.VolumeConstraint)}
		}
	}

	if e.Composition != nil || len(e.Composition) > 0 {
		if e.VolumeType != "count" && e.VolumeConstraint != 1 {
			return ErrInvalidExercise{Message: "composites must use the count volume type with volume constraint of 1"}
		}
	}

	if (e.Composition != nil || len(e.Composition) > 0) && e.Basis != "" {
		return ErrInvalidExercise{Message: "cannot be both a composite and a variation"}
	}

	// Restrict to sensible combinations of intensity and volume types
	if e.IntensityType == "weight" || e.IntensityType == "bodyweight" || e.IntensityType == "percentOfMax" {
		if e.VolumeType != "count" {
			return ErrInvalidExercise{Message: "weight-based intensities muse use count as volume type"}
		}
	} else if e.IntensityType == "hrZone" {
		if e.VolumeType != "time" {
			return ErrInvalidExercise{Message: "HR Zone intensities must use time as volume type"}
		}
	} else if e.IntensityType == "pace" {
		if e.VolumeType != "time" && e.VolumeType != "distance" {
			return ErrInvalidExercise{Message: "Pace intensities must use time as volume type"}
		}
	}
	// RPE is valid for any volume type

	return nil
}

func (et ExerciseType) CalculateMetrics(ei *ExerciseInstance) (load, volume float32) {
	load = 0
	volume = 0

	volumeFactor := float32(1)

	// time volume is converted from seconds to hours
	if et.VolumeType == "time" {
		volumeFactor = float32(1) / 60 / 60
	}

	// distance is converted from m to km
	if et.VolumeType == "distance" {
		volumeFactor = float32(1) / 1000
	}

	// for weight and count, volume is the total reps and load is total weight * reps
	// for distance, volume is number of km
	// for time, volume is number of hours
	for _, segment := range ei.Segments {
		for _, set := range segment.Volume {
			for _, reps := range set {
				volume = volume + reps*volumeFactor
				load = load + segment.Intensity*reps*volumeFactor
			}

		}
	}
	return
}
