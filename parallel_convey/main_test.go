package pc

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestParallelConveyWithRoot(t *testing.T) {
	t.Parallel()

	Convey("root", t, func() {
		ParallelConvey, Wait := New(t)

		start := time.Now()

		ParallelConvey("1", func() {
			time.Sleep(3 * time.Second)
			So(1, ShouldEqual, 1)
		})

		ParallelConvey("2", func() {
			time.Sleep(3 * time.Second)
			So(2, ShouldEqual, 2)
		})

		Wait()

		elapsed := time.Since(start)

		So(elapsed.Seconds(), ShouldBeLessThan, 3.5)
	})
}

func TestParallelConveyTopLevel(t *testing.T) {
	t.Parallel()

	ParallelConvey, Wait := New()

	start := time.Now()

	ParallelConvey("3", t, func() {
		time.Sleep(3 * time.Second)
		So(1, ShouldEqual, 1)
	})

	ParallelConvey("4", t, func() {
		time.Sleep(3 * time.Second)
		So(2, ShouldEqual, 2)
	})

	Wait()

	elapsed := time.Since(start)

	if elapsed.Seconds() > 3.5 {
		t.Fatalf("TestParallelConveyTopLevel took too long: %v", elapsed)
	}
}

func TestParallelConveyWithInheritedT(t *testing.T) {
	t.Parallel()

	ParallelConvey, Wait := New(t)

	start := time.Now()

	ParallelConvey("5", func() {
		time.Sleep(3 * time.Second)
		So(1, ShouldEqual, 1)
	})

	ParallelConvey("6", func() {
		time.Sleep(3 * time.Second)
		So(2, ShouldEqual, 2)
	})

	Wait()

	elapsed := time.Since(start)

	if elapsed.Seconds() > 3.5 {
		t.Fatalf("TestParallelConveyWithNoT took too long: %v", elapsed)
	}
}
