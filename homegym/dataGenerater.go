package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/scottbrodersen/homegym/auth"
	"github.com/scottbrodersen/homegym/programs"
	"github.com/scottbrodersen/homegym/workoutlog"
)

const (
	username = "test"
	pwd      = "testtesttest"
	email    = "test@example.com"
	activity = "weightlifting"
)

var squatID string = ""
var snatchID string = ""

func AddData() error {

	testDate := time.Now().Unix()
	event := workoutlog.Event{Date: testDate}
	page, err := workoutlog.EventManager.GetPageOfEvents("test", event, 5)
	if err != nil {
		fmt.Println("error getting page: %w", err)
	}
	for _, v := range page {
		fmt.Printf("%v\n", v)
	}

	log.Print("adding data")
	err = createUser()
	if err != nil {
		return err
	}

	err = createExerciseTypes()
	if err != nil {
		return err
	}

	activity, err := createActivity()
	if err != nil {
		return err
	}

	now := time.Now().Unix() - int64(60*60*24*30)

	for i := 0; i < 25; i++ {
		date := now + int64(i*3600*25)
		_, err := newEvent(date, *activity)
		if err != nil {
			return err
		}
	}

	addProgram(*activity, "program 1")
	p2 := addProgram(*activity, "program 2")

	programInstanceMaker(p2, testDate)

	return nil
}

func createUser() error {
	_, err := workoutlog.FrontDesk.NewUser(username, auth.User, email, pwd)
	if err != nil {
		return err
	}
	log.Print("user created")
	return nil
}

func createExerciseTypes() error {
	id, err := workoutlog.ExerciseManager.NewExerciseType(username, "squat", "weight", "count", 1, nil, "")
	if err != nil {
		return err
	}
	log.Print("squat exercise type created")
	squatID = *id

	id, err = workoutlog.ExerciseManager.NewExerciseType(username, "snatch", "weight", "count", 2, nil, "")
	if err != nil {
		return err
	}
	log.Print("snatch exercise type created")
	snatchID = *id

	return nil
}

func createActivity() (*workoutlog.Activity, error) {
	activity, err := workoutlog.ActivityManager.NewActivity(username, activity)
	if err != nil {
		return nil, err
	}
	log.Print("activity created")

	if err := activity.AddExerciseToActivity(username, squatID); err != nil {
		return nil, err
	}

	if err := activity.AddExerciseToActivity(username, snatchID); err != nil {
		return nil, err
	}

	log.Print("exercise types added to activity")

	return activity, nil
}

func newEvent(date int64, activity workoutlog.Activity) (*workoutlog.Event, error) {
	event := workoutlog.Event{
		Date:       date,
		ActivityID: activity.ID,
		EventMeta:  metaMaker(),
	}

	exInstances, err := exerciseInstancesMaker(activity)
	if err != nil {
		return nil, err
	}

	event.Exercises = exInstances

	eventID, err := workoutlog.EventManager.NewEvent(username, event)
	if err != nil {
		return nil, err
	}

	event.ID = *eventID
	log.Print("event created")
	return &event, nil
}

func metaMaker() workoutlog.EventMeta {
	return workoutlog.EventMeta{
		Mood:       int(rand.Int63n(5)),
		Motivation: int(rand.Int63n(5)),
		Energy:     int(rand.Int63n(5)),
		Overall:    int(rand.Int63n(5)),
		Notes:      "random note",
	}

}

func exerciseInstancesMaker(activity workoutlog.Activity) (map[int]workoutlog.ExerciseInstance, error) {
	numInstances := rand.Intn(4)

	exerciseTypes := []workoutlog.ExerciseType{}
	for _, etID := range activity.ExerciseIDs {
		e, err := workoutlog.ExerciseManager.GetExerciseType(username, etID)
		if err != nil {
			return nil, err
		}
		exerciseTypes = append(exerciseTypes, *e)
	}

	exerciseInstances := map[int]workoutlog.ExerciseInstance{}

	for i := 0; i < numInstances; i++ {
		etID := rand.Intn(len(exerciseTypes))
		inst := exerciseTypes[etID].CreateInstance()
		inst.Index = int(i)
		numParts := rand.Intn(3)
		for p := 0; p <= numParts; p++ {
			inst.Segments = append(inst.Segments, exerciseSegmentMaker())
		}
		exerciseInstances[int(i)] = inst
	}

	return exerciseInstances, nil
}

func exerciseSegmentMaker() workoutlog.ExerciseSegment {
	max := float32(150.0)
	min := float32(20.0)
	intensity := float32(min + rand.Float32()*(max-min))

	vol := [][]float32{}

	numSets := rand.Intn(4) + 1

	for i := 0; i < numSets; i++ {
		set := []float32{}
		numReps := rand.Intn(6) + 1
		for j := 0; j < numReps; j++ {
			set = append(set, float32(rand.Intn(2)))
		}
		vol = append(vol, set)
	}
	return workoutlog.ExerciseSegment{Intensity: intensity, Volume: vol}

}

func addProgram(activity workoutlog.Activity, title string) programs.Program {

	microcycle1 := programs.MicroCycle{
		Title:       title + " week 1",
		Span:        7,
		Description: "Intensity of the first week",
		Workouts: []programs.Workout{
			workoutMaker("Day1", "heavy single X 5", "heavy double X 3"),
			workoutMaker("Day2", "80% single X 5", "heavy double X 3"),
			workoutMaker("Day3", "", ""),
			workoutMaker("Day4", "70% doubles X 5", "heavy double X 3"),
			workoutMaker("Day5", "max out", "heavy double X 3"),
			workoutMaker("Day6", "", ""),
			workoutMaker("Day7", "", ""),
		},
	}

	microcycle2 := programs.MicroCycle{
		Title:       title + " week 2",
		Span:        7,
		Description: "Intensity of the second week",
		Workouts: []programs.Workout{
			workoutMaker("Monday", "80% doubles x 3", "heavy double X 3"),
			workoutMaker("Tuesday", "80% single X 5", "heavy double X 3"),
			workoutMaker("Wednesday", "", ""),
			workoutMaker("Thursday", "70% doubles X 5", "heavy double X 3"),
			workoutMaker("Friday", "max out", "heavy double X 3"),
			workoutMaker("Saturday", "", ""),
			workoutMaker("Sunday", "", ""),
		},
	}

	block1 := programs.Block{
		Title:       title + " Max out",
		Description: "High intensity",
		MicroCycles: []programs.MicroCycle{
			microcycle1,
			microcycle2,
		},
	}
	block2 := programs.Block{
		Title:       title + " Taper",
		Description: "Test",
		MicroCycles: []programs.MicroCycle{
			microcycle1,
			microcycle2,
		},
	}
	program := programs.Program{
		Title:      title,
		ActivityID: activity.ID,
		Blocks: []programs.Block{
			block1,
			block2,
		},
	}

	programID, err := programs.ProgramManager.AddProgram(username, program)
	if err != nil {
		fmt.Println("error creating program")
		fmt.Println(err.Error())
	}

	program.ID = *programID

	fmt.Println("program created")

	return program
}

func workoutMaker(title, pSn, pSq string) programs.Workout {
	if pSn == "" && pSq == "" {
		return programs.Workout{
			Title:       title,
			RestDay:     true,
			Description: "Rest Day",
		}
	}
	ws1 := programs.WorkoutSegment{
		ExerciseTypeID: snatchID,
		Prescription:   pSn,
	}

	ws2 := programs.WorkoutSegment{
		ExerciseTypeID: squatID,
		Prescription:   pSq,
	}

	segments := []programs.WorkoutSegment{
		ws1,
		ws2,
	}

	return programs.Workout{
		Title:       title,
		Segments:    segments,
		Description: "Leg day",
	}

}

func programInstanceMaker(program programs.Program, testDate int64) {
	startTime := testDate - 3.5*24*60*60
	program.Title = "Test program instance"

	pi := programs.ProgramInstance{
		Program:   program,
		StartTime: startTime,
		ProgramID: program.ID,
	}

	pi.Program.ID = ""
	programs.ProgramManager.AddProgramInstance(username, &pi)

}
