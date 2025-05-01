package pc

import (
	"sync"

	"github.com/smartystreets/goconvey/convey"
)

// Creates a new WaitGroup and returns a Wrapper around `Convey` which
// automatically adds to the WaitGroup and calls `Done` when the test is
// complete. This allows you to run multiple tests in parallel without
// blocking the main thread. The returned function can be used to call
// `Convey` with the same arguments as `Convey`, but it will run in a
// separate goroutine and will not block the main thread. The WaitGroup
// can be used to wait for all tests to complete before exiting the program.
//
// Example usage:
//
//	ParallelConvey, Wait := pc.New()
//
//	ParallelConvey("Test 1", func() {
//		// Test code here
//	})
//
//	ParallelConvey("Test 2", func() {
//		// Test code here
//	})
//
//	Wait()
func New() (func(), func(...any)) {
	wg := &sync.WaitGroup{}
	return wg.Wait, func(items ...any) {
		wg.Add(1)
		go func() {
			convey.Convey(items...)
			wg.Done()
		}()
	}
}
