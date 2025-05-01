package pc

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestParallelConvey(t *testing.T) {

	Convey("root", t, func() {
		ParallelConvey, Wait := New()

		start := time.Now()

		ParallelConvey("1", t, func() {
			time.Sleep(3 * time.Second)
			So(1, ShouldEqual, 1)
		})

		ParallelConvey("2", t, func() {
			time.Sleep(3 * time.Second)
			So(2, ShouldEqual, 2)
		})

		Wait()

		elapsed := time.Since(start)

		So(elapsed.Seconds(), ShouldBeLessThan, 3.5)
	})
}
