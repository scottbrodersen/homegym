package dailystats

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/scottbrodersen/homegym/dal"
)

const (
	testUserID     = "testUserID"
	testDate       = int64(1720181150)
	testBG         = float32(5)
	testSleep      = float32(8)
	testProtein    = 20
	testFat        = 20
	testCarbs      = 2
	testFiber      = 10
	testBodyWeight = 135
	testMood       = 3
	testStress     = 3
	testEnergy     = 3
)

var testBP []int = []int{165, 75}

func testStats() DailyStats {
	return DailyStats{
		Date:          testDate,
		BloodGlucose:  testBG,
		BloodPressure: testBP,
		Sleep:         testSleep,
		Food: Food{
			Protein: testProtein,
			Carbs:   testCarbs,
			Fat:     testFat,
			Fiber:   testFiber,
		},
		BodyWeight: testBodyWeight,
		Mood:       testMood,
		Stress:     testStress,
		Energy:     testEnergy,
	}
}

func TestDailyStats(t *testing.T) {
	Convey("Given a dal client", t, func() {
		db := dal.NewMockDal()
		dal.DB = db

		Convey("When we add a daily stat", func() {
			db.On("AddBioStats", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			stat := testStats()

			statJSON, err := json.Marshal(stat)
			if err != nil {
				t.Fail()
			}

			err = DailyStatsManager.AddStats(testUserID, testDate, statJSON)

			So(err, ShouldBeNil)

		})

		Convey("When we get a page of daily stats", func() {
			count := 0
			statsJSON := [][]byte{}
			for count < 2 {
				stat := testStats()
				stat.Date = stat.Date + int64(count)
				statJSON, err := json.Marshal(stat)
				if err != nil {
					t.Fail()
				}

				statsJSON = append(statsJSON, statJSON)
				count++
			}
			db.On("GetBioStatsPage", mock.Anything, mock.Anything, mock.Anything).Return(statsJSON, nil)
			dailyStats, err := DailyStatsManager.GetBioStatsPage(testUserID, testDate+int64(1), testDate-int64(1), 10)

			So(err, ShouldBeNil)
			So(dailyStats, ShouldNotBeNil)

			stats := []DailyStats{}

			err = json.Unmarshal(dailyStats, &stats)
			So(err, ShouldBeNil)
			So(len(stats), ShouldEqual, 2)
		})
	})
}
