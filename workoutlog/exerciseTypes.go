package workoutlog

import (
	"fmt"
	"math"
	"regexp"
	"slices"
)

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
}

// Indicates the type of values that can be expressed for volumes
// 0 is no restriction (any float32)
// 1 is restricted to the value 1 (for counts)
// 2 is restricted to either 1 or 0 (success/failure)
var volumeConstraints []int = []int{0, 1, 2}
var volumeTypes []string = []string{"count", "distance", "time"}
var intensityTypes []string = []string{"weight", "distance", "hrZone", "rpe", "percentOfMax", "bodyweight"}

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
func (et ExerciseType) validateInstance(ei *ExerciseInstance) error {
	for i, segment := range ei.Segments {
		// Validate intenstiy values
		if segment.Intensity <= 0 {
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
				return fmt.Errorf("hrZone must be between 1 and 5")
			}
			fallthrough
		case "rpe":
			if segment.Intensity > 10 && segment.Intensity < 1 {
				return fmt.Errorf("RPI must be between 1 and 10")
			}
			fallthrough
		case "percentOfMax":
			ei.Segments[i].Intensity = float32(math.Floor(float64(segment.Intensity)))
		}

		// Validate volume values
		for j, set := range segment.Volume {

			if len(set) == 0 {
				return fmt.Errorf("volume is a required value")
			}
			for k, rep := range set {
				if rep < 0 {
					return fmt.Errorf("volume must be a positive number")
				}
				switch et.VolumeConstraint {
				case 0:
					ei.Segments[i].Volume[j][k] = float32(math.Floor(float64(rep*10)) / 10)

				case 1:
					fallthrough
				case 2:
					if rep != 0 && rep != 1 {
						return fmt.Errorf("invalid rep value: must be 0 or 1")
					}
				}
			}
		}
	}
	return nil
}

// ExerciseInstance stores data about the performance of an exercise type.
// Index stores the location in the order of performed exercies in an event
type ExerciseInstance struct {
	TypeID   string            `json:"typeID"`
	Index    uint64            `json:"index"`
	Segments []ExerciseSegment `json:"parts"`
}

// ExerciseSegment stores the performance metrics for an exercise performed at a specific intensity.
// The untyped Volume field enables flexibility in expressing volume data.
// The volume type of the associated exercise type should perform type assertion on the Volume fields.
type ExerciseSegment struct {
	Intensity float32     `json:"intensity"`
	Volume    [][]float32 `json:"volume"`
}

func (e ExerciseType) validate() error {
	if e.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	namerxp := regexp.MustCompile(`^[a-zA-Z0-9_\-&$*@. \+]+$`)
	if ok := namerxp.MatchString(e.Name); !ok {
		return fmt.Errorf("name can include letters, digits, _, -, &, $, *, @, ., and spaces")
	}

	if e.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if slices.Index(intensityTypes, e.IntensityType) < 0 {
		return fmt.Errorf("invalid intensity type: %s", e.IntensityType)
	}

	if slices.Index(volumeTypes, e.VolumeType) < 0 {
		return fmt.Errorf("invalid volume type: %s", e.VolumeType)
	}

	if slices.Index(volumeConstraints, e.VolumeConstraint) < 0 {
		return fmt.Errorf("invalid volume constraint: %d", e.VolumeConstraint)
	}

	if e.VolumeType == "count" {
		if e.VolumeConstraint != 1 && e.VolumeConstraint != 2 {
			return fmt.Errorf("invalid volume constraint: %d", e.VolumeConstraint)
		}

	} else {
		if e.VolumeConstraint != 0 {
			return fmt.Errorf("invalid volume constraint: %d", e.VolumeConstraint)
		}
	}

	if e.Composition != nil {
		if e.VolumeType != "count" && e.VolumeConstraint != 1 {
			return fmt.Errorf("composites must use the count volume type with volume constraint of 1")
		}
	}

	return nil
}
