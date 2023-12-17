package dal

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testfolder       = "./temp/"
	testdb           = "testdb"
	testUserID       = "testuser"
	testActivtyName  = "testActivity"
	testEmail        = "test@example.com"
	testUserName     = "test user"
	testHash         = "testhash"
	testPwdVersion   = "v1"
	testExerciseName = "testExercise"
	tokenUsage       = "token"
	testSessionID    = "test-session-id"
	testRole         = "testRole"
	testSessionTTL   = 14
)

var testEvent []byte = []byte("test event")
var testExercise = []byte("test exercise type")
var testpath string = fmt.Sprintf("%s%s", testfolder, testdb)

func init() {
	cleanup()
}
func TestDal(t *testing.T) {
	var client *DBClient
	var err error
	defer cleanup()
	Convey("Given an even number of strings", t, func() {
		items := []string{"one", "two", "three", "four"}
		Convey("When we generate a key", func() {
			key := key(items)
			Convey("the key should be as expected", func() {
				So(bytes.Equal(key, []byte("one:two#three:four")), ShouldBeTrue)
			})
		})
		Convey("When we add another item and generate a key", func() {
			odditems := append(items[:], "five")
			key := key(odditems)
			Convey("the key should be as expected", func() {
				So(bytes.Equal(key, []byte("one:two#three:four#five")), ShouldBeTrue)
			})
		})
	})
	Convey("When we create a dal client", t, func() {
		client, err = InitClient(testpath)
		defer client.Destroy()
		Convey("Then the database is created", func() {
			So(err, ShouldBeNil)
		})
	})
}

func TestUsersDal(t *testing.T) {
	defer cleanup()
	Convey("Given a dal client", t, func() {
		client, err := InitClient(testpath)
		if err != nil {
			log.Fatal("failed to create dal client")
		}
		defer client.Destroy()
		Convey("When we add a user", func() {
			err := client.NewUser(testUser(testUserID))
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("When we add the user again", func() {
			err := client.NewUser(testUser(testUserID))
			Convey("Then the expected error is returned", func() {
				So(err, ShouldEqual, ErrNotUnique)
			})
		})
		Convey("When we add a user with id that starts with the entire id of the other user", func() {
			err := client.NewUser(testUser(testUserID + "2"))
			Convey("Then noerror is returned", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("When we read the user", func() {
			e, ph, phv, r, err := client.ReadUser(testUserID)

			Convey("The returned user is as expected", func() {
				So(err, ShouldBeNil)
				So(*e, ShouldEqual, testEmail)
				So(*ph, ShouldEqual, testHash)
				So(*phv, ShouldEqual, testPwdVersion)
				So(*r, ShouldEqual, testRole)
			})
		})

		Convey("When we update the user details", func() {
			err := client.UpdateUserProfile(testUserID, "newemail@example.com")
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)

				er, _, _, _, _ := client.ReadUser(testUserID)
				So(*er, ShouldEqual, "newemail@example.com")
			})
		})

		Convey("When we update the user password", func() {
			newHash := "newhash"
			pwdVersion := "newversion"
			err := client.UpdateUserPassword(testUserID, newHash, pwdVersion)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)

				_, ph, phv, _, _ := client.ReadUser(testUserID)
				So(*ph, ShouldEqual, newHash)
				So(*phv, ShouldEqual, pwdVersion)
			})
		})

		Convey("When we update the user role", func() {
			newRole := "newRole"
			err := client.ChangeUserRole(testUserID, newRole)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)

				_, _, _, rr, _ := client.ReadUser(testUserID)
				So(*rr, ShouldEqual, newRole)
			})
		})

		Convey("When we update the user with a new hashing version", func() {
			newPwdHashVersion := "blah"
			err := client.UpdatePwdVersion(testUserID, newPwdHashVersion)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				_, _, phv, _, _ := client.ReadUser(testUserID)
				So(*phv, ShouldEqual, newPwdHashVersion)
			})
		})
	})
}

var (
	testActivityID       string = "testActivityID"
	testExerciseID              = "testExerciseID"
	testEventTime               = time.Now().Unix()
	testEventID                 = "testEventID"
	testExerciseInstance        = []byte("testExercisInstance")
)

func TestLogItemsDal(t *testing.T) {

	testExerciseIDs := map[int]string{0: "id1", 1: "id2"}
	testExerciseInstances := map[int][]byte{0: testExerciseInstance, 1: testExerciseInstance}
	testExerciseIDs2 := map[int]string{0: "id1", 1: "id4", 2: "id3"}
	testExerciseInstances2 := map[int][]byte{0: testExerciseInstance, 1: testExerciseInstance, 2: testExerciseInstance}

	defer cleanup()
	Convey("Given a dal client", t, func() {
		client, err := InitClient(testpath)
		if err != nil {
			log.Fatal()
		}
		defer client.Destroy()
		Convey("when we add an activity", func() {
			err = client.AddActivity(testUserID, testActivityID, testActivtyName)
			Convey("Then nil error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("When we add the activity again", func() {
			err := client.AddActivity(testUserID, testActivityID, testActivtyName)
			Convey("Then nils are returned", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("When we change the activity name to a different name", func() {
			err := client.UpdateActivity(testUserID, testActivityID, "newname")
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("When we add an exercise to the activity", func() {
			err = client.AddExerciseToActivity(testUserID, testActivityID, testExerciseID)
			Convey("A nil error is returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When we read the activity", func() {
			name, exerciseIDs, err := client.ReadActivity(testUserID, testActivityID)
			Convey("Then the returned activity is as expected", func() {
				So(err, ShouldBeNil)
				So(name, ShouldNotBeNil)
				So(name, ShouldNotBeEmpty)
				So(exerciseIDs, ShouldNotBeNil)
				So(len(exerciseIDs), ShouldEqual, 1)
				So(exerciseIDs[0], ShouldEqual, testExerciseID)
			})
		})

		Convey("When we attempt to read an activity that has not been added", func() {
			name, exerciseIDs, err := client.ReadActivity(testUserID, "badID")
			Convey("Then an error is returned", func() {
				So(err, ShouldBeNil)
				So(name, ShouldBeNil)
				So(exerciseIDs, ShouldBeNil)
			})
		})

		Convey("When we add an event to the activity", func() {
			err = client.AddExerciseToActivity(testUserID, testActivityID, testExerciseID)
			So(err, ShouldBeNil)

		})

		Convey("When we add and delete exercises to an activity", func() {
			adds := []string{"a", "b", "c", "d", "e"}
			deletes := []string{testExerciseID}
			err := client.UpdateActivityExercises(testUserID, testActivityID, adds, deletes)
			So(err, ShouldBeNil)

			Convey("And when we read the exercises", func() {
				_, exercises, err := client.ReadActivity(testUserID, testActivityID)
				So(err, ShouldBeNil)
				So(exercises, ShouldResemble, adds)
			})
		})

		Convey("When we add an exercise", func() {
			err := client.AddExercise(testUserID, testExerciseID, testExercise)

			So(err, ShouldBeNil)
		})

		Convey("When we read the exercises", func() {
			exercises, err := client.GetExercises(testUserID)

			So(err, ShouldBeNil)
			So(exercises, ShouldNotBeEmpty)
			So(len(exercises), ShouldEqual, 1)
			So(string(exercises[0]), ShouldEqual, string(testExercise))
		})

		Convey("When we update an exercise", func() {
			err := client.UpdateExercise(testUserID, testExerciseID, []byte("New"))

			So(err, ShouldBeNil)
		})

		Convey("When we read the exercise again", func() {
			exercise, err := client.GetExercise(testUserID, testExerciseID)

			So(err, ShouldBeNil)
			So(exercise, ShouldNotBeEmpty)
			So(string(exercise), ShouldEqual, "New")
		})

		Convey("When we add an event", func() {
			err := client.AddEvent(testUserID, testEventID, testActivityID, testEventTime, testEvent)

			So(err, ShouldBeNil)
		})

		Convey("When we add an exercise to the event", func() {
			err := client.AddExercisesToEvent(testUserID, testEventID, testExerciseIDs, testExerciseInstances)

			So(err, ShouldBeNil)
		})

		Convey("When we get the event's exercises", func() {
			eventExercises, err := client.GetEventExercises(testUserID, testEventID)

			So(err, ShouldBeNil)
			So(len(eventExercises), ShouldEqual, len(testExerciseInstances))
			So(eventExercises, ShouldResemble, [][]byte{testExerciseInstances[0], testExerciseInstances[1]})
			So(eventExercises[1], ShouldResemble, testExerciseInstances[1])
		})

		Convey("When we update the event exercises", func() {
			err := client.AddExercisesToEvent(testUserID, testEventID, testExerciseIDs2, testExerciseInstances2)

			So(err, ShouldBeNil)
		})

		Convey("When we get the event's exercises after the update", func() {
			eventExercises, err := client.GetEventExercises(testUserID, testEventID)

			So(err, ShouldBeNil)
			So(len(eventExercises), ShouldEqual, len(testExerciseInstances2))
			So(eventExercises, ShouldResemble, [][]byte{testExerciseInstances2[0], testExerciseInstances2[1], testExerciseInstances2[2]})
			So(eventExercises[1], ShouldResemble, testExerciseInstances[1])
		})

		Convey("When we get the event", func() {
			eventByte, err := client.GetEvent(testUserID, testEventID, testEventTime)

			So(err, ShouldBeNil)
			So(eventByte, ShouldNotBeNil)
			So(eventByte, ShouldResemble, testEvent)
		})

		newTime := testEventTime + 1
		Convey("When we shift the event's time", func() {
			err := client.ShiftEvent(testUserID, testEventID, testActivityID, testEventTime, newTime, testEvent)

			So(err, ShouldBeNil)
		})

		Convey("When we get the shifted event", func() {
			eventByte, err := client.GetEvent(testUserID, testEventID, newTime)

			So(err, ShouldBeNil)
			So(eventByte, ShouldNotBeNil)
		})

		Convey("When we attempt to get the event with using the previous event time", func() {
			event, err := client.GetEvent(testUserID, testEventID, testEventTime)

			So(err, ShouldBeNil)
			So(event, ShouldBeNil)
		})

		Convey("When we add many events", func() {
			timeIncrement := 100
			testEvents := [][]byte{testEvent}
			testEventTimes := []int64{testEventTime}
			for i := 1; i < 25; i++ {
				eventID := fmt.Sprintf("%s%d", testEventID, i)
				eventTime := testEventTime + int64(timeIncrement*i)
				event := []byte(fmt.Sprintf("%s_%d", eventID, eventTime))
				err := client.AddEvent(testUserID, eventID, testActivityID, eventTime, event)
				So(err, ShouldBeNil)
				testEvents = append(testEvents, event)
				testEventTimes = append(testEventTimes, eventTime)
			}
			pagesize := 10
			middleEvent := ""
			lastEvent := ""
			Convey("And then we get a page of events", func() {
				// startime must be greater than the most recent test event
				startTime := testEventTimes[len(testEventTimes)-1] + 1
				eventsByte, err := client.GetEventPage(testUserID, "", int64(startTime), int(pagesize))
				So(err, ShouldBeNil)
				So(len(eventsByte), ShouldEqual, pagesize)
				So(eventsByte[0], ShouldResemble, testEvents[len(testEvents)-1])
				So(eventsByte[pagesize-1], ShouldResemble, testEvents[len(testEvents)-pagesize])

				middleEvent = string(eventsByte[pagesize/2-1])
				lastEvent = string(eventsByte[pagesize-1])

				Convey("And then we get a half-size page of events starting at the middle event of the previous full page", func() {
					startEventID := strings.Split(middleEvent, "_")[0]
					startTime, _ := strconv.Atoi(strings.Split(middleEvent, "_")[1])
					events, err := client.GetEventPage(testUserID, startEventID, int64(startTime), int(pagesize/2))
					So(err, ShouldBeNil)
					So(len(events), ShouldEqual, pagesize/2)
					lastInPage := string(events[pagesize/2-1])
					So(lastEvent, ShouldEqual, lastInPage)
				})
			})
		})
	})
}

func TestKeysDal(t *testing.T) {
	id1 := uuid.NewString()
	id2 := uuid.NewString()
	id3 := uuid.NewString()
	testKey1 := []byte("testKey1")
	testKey2 := []byte("testKey2")
	testKey3 := []byte("testKey3")
	defer cleanup()

	Convey("Given a dal client", t, func() {
		db, err := InitClient(testpath)
		if err != nil {
			log.Fatal()
		}
		defer db.Destroy()
		Convey("When we add the first key", func() {
			err := db.RotateKeys(testKey1, id1, tokenUsage)
			So(err, ShouldBeNil)
			Convey("And when we get the key", func() {
				active, retired, err := db.GetKeys(tokenUsage)
				So(err, ShouldBeNil)
				So(len(active), ShouldEqual, 1)
				So(len(retired), ShouldEqual, 0)
				for k, v := range active {
					So(k, ShouldEqual, id1)
					So(v, ShouldResemble, testKey1)
				}
			})
		})

		Convey("When we rotate the keys and get all keys", func() {
			err = db.RotateKeys(testKey2, id2, tokenUsage)
			So(err, ShouldBeNil)
			active, retired, err := db.GetKeys(tokenUsage)

			So(err, ShouldBeNil)

			So(len(active), ShouldEqual, 1)
			activeKey, ok := active[id2]
			So(ok, ShouldBeTrue)
			So(activeKey, ShouldResemble, testKey2)

			So(len(retired), ShouldEqual, 1)
			retiredKey, ok := retired[id1]
			So(ok, ShouldBeTrue)
			So(retiredKey, ShouldResemble, testKey1)
		})

		Convey("When we rotate the keys again and get all keys", func() {
			err := db.RotateKeys(testKey3, id3, tokenUsage)
			So(err, ShouldBeNil)
			active, retired, err := db.GetKeys(tokenUsage)

			So(err, ShouldBeNil)
			activeKey, ok := active[id3]
			So(ok, ShouldBeTrue)
			So(activeKey, ShouldResemble, testKey3)

			So(len(retired), ShouldEqual, 2)
		})

		Convey("When we delete the retired keys", func() {
			_, retired, err := db.GetKeys(tokenUsage)
			if err != nil {
				log.Fatal()
			}
			errCount := 0
			for id := range retired {
				err := db.DeleteKey(id, tokenUsage)
				if err != nil {
					errCount++
				}
			}

			So(errCount, ShouldEqual, 0)
			Convey("There should be one active key and no retired keys", func() {
				active, retired, err := db.GetKeys(tokenUsage)
				if err != nil {
					log.Fatal()
				}
				So(len(active), ShouldEqual, 1)
				So(len(retired), ShouldEqual, 0)
			})
		})
	})
}

func TestSessionssDal(t *testing.T) {
	defer cleanup()
	Convey("Given a dal client", t, func() {
		client, err := InitClient(testpath)
		if err != nil {
			log.Fatal()
		}
		defer client.Destroy()
		Convey("When we track an authenticated user", func() {
			err := client.AddSession(testUserID, testSessionID, testSessionTTL)
			So(err, ShouldBeNil)
			username, expiry, err := client.GetSession(testSessionID)
			So(err, ShouldBeNil)
			So(username, ShouldNotBeNil)
			So(*username, ShouldEqual, testUserID)
			So(expiry, ShouldNotBeNil)
		})
		Convey("When we add the user again", func() {
			err = client.AddSession(testUserID, testSessionID, testSessionTTL)
			So(err, ShouldBeNil)
		})
		Convey("When we delete the session", func() {
			err = client.DeleteSession(testSessionID)
			So(err, ShouldBeNil)
			username, exprity, err := client.GetSession(testSessionID)
			So(err, ShouldBeNil)
			So(username, ShouldBeNil)
			So(exprity, ShouldBeNil)
		})
		Convey("When we get session expiry times", func() {
			expiries, err := client.GetSessionExpiries()
			So(err, ShouldBeNil)
			So(expiries, ShouldBeEmpty)
			err = client.AddSession(testUserID, testSessionID, testSessionTTL)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = client.AddSession("testuserid2", "testsessionid2", testSessionTTL)
			if err != nil {
				log.Fatal(err.Error())
			}
			expiries, err = client.GetSessionExpiries()
			So(err, ShouldBeNil)
			So(len(expiries), ShouldEqual, 2)
		})
	})
}

func cleanup() {
	err := os.RemoveAll(testfolder)
	if err != nil {
		log.WithError(err).Error("failed to delete test db")
	}
}

func testUser(id string) (string, string, string, string, string) {
	return id, testEmail, testHash, testPwdVersion, testRole

}
