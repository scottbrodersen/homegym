package dal

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestBioDal(t *testing.T) {
	testDate := int64(1720181149)
	testBioStat := []byte("test bio stats")

	defer cleanup()
	Convey("Given a dal client", t, func() {
		client, err := InitClient(testPath)
		if err != nil {
			t.Fatal()
		}
		defer client.Destroy()

		Convey("when we add a bio stat", func() {
			err := client.AddBioStats(testUserID, testDate, testBioStat)
			Convey("Then nil error is returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When we get the bio stat", func() {
			bioStats, err := client.GetBioStatsPage(testUserID, testDate, 1)

			So(err, ShouldBeNil)

			So(bioStats, ShouldNotBeNil)
			So(len(bioStats), ShouldEqual, 1)
			So(bioStats[0], ShouldResemble, testBioStat)

		})

		Convey("When we get a page of bio stats", func() {
			err := client.AddBioStats(testUserID, testDate+int64(1), testBioStat)
			if err != nil {
				t.Fail()
			}
			bioStats, err := client.GetBioStatsPage(testUserID, 0, 10)

			So(err, ShouldBeNil)

			So(bioStats, ShouldNotBeNil)
			So(len(bioStats), ShouldEqual, 2)
			So(bioStats[0], ShouldResemble, testBioStat)
			So(bioStats[1], ShouldResemble, testBioStat)

		})
	})
}
