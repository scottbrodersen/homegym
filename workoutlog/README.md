### Exercise Types

`ExerciseType` defines exercises in terms of intensity and volume. It is a factory for `ExerciseInstance` structs that hold information about the performance of an exercise. `ExerciseType` also validates the data.

A limited number of intensity types and volume types are used to define how to express the performance of an exercise. For example, squats use weight as intensity and sets and reps to indicate volume. Tempo runs can use heart rate zones as intensity and time as volume.

Exercise performance data is stored in the same format for all exercises. The `VolumeConstraint` property indicates how volume data should be interpreted by the client.

```mermaid
    classDiagram
direction LR
    class ExerciseType{
        <<type>>
        Name string
        ID string
        IntensityType string
        VolumeType string
        VolumeConstraint int
        validate() error
        validateInstance() error
        CreateInstance() ExerciseInstance
    }

    class volumeTypes {
      << collection >>
      id string
    }

    class intensityTypes {
      << collection >>
      id string
    }

    class volumeConstraints {
       << collection >>
     enum int
    }

    class ExerciseInstance{
      <<type>>
      TypeID string
      Index int
      Segments []ExerciseSegment
    }

    class ExerciseSegment{
      Intensity float32
      Volume [][]float32
    }

    ExerciseType --> ExerciseInstance :creates
    ExerciseInstance o-- ExerciseSegment
    ExerciseType ..> volumeTypes
    ExerciseType ..> intensityTypes
    ExerciseType ..> volumeConstraints
```
