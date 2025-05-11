package pc

import (
	"os"
	"sync"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// Creates a new WaitGroup and returns a Wrapper around `Convey` which
// automatically adds to the WaitGroup and calls `Done` when the test is
// complete. This allows you to run multiple tests in parallel without
// blocking the main thread. The returned function can be used to call
// `Convey` with the same arguments as `Convey`, but it will run in a
// separate goroutine and will not block the main thread. The WaitGroup
// can be used to wait for all tests to complete before exiting the program.
// However, if you do not provide a *testing.T to `ParallelConvey` at the time of usage,
// then you should provide it once at the time of instantiation.
//
// !!WARNING: Using this will jumble up the test reporting on the console due to parallel execution (The ones with the ticks and crosses).
//
// Example usage:
//
//	ParallelConvey, Wait := pc.New(t)
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
//
// If you want to disable parallel execution, for example for debugging purposes,
// you can set the environment variable PARALLEL_CONVEY to `false` before running the tests:
// PARALLEL_CONVEY=false go test -v
// This will make it behave like the original `Convey` function and run the tests sequentially.
func New(t ...*testing.T) (func(...any), func()) {
	wg := &sync.WaitGroup{}
	return func(items ...any) {
		if os.Getenv("PARALLEL_CONVEY") == "false" {
			convey.Convey(items...)
			return
		}
		wg.Add(1)
		if len(items) > 1 && len(t) > 0 { // If *testing.T is provided, inject it into convey if not already present
			if _, ok := items[1].(func()); ok {
				items = append(items[:1], append([]any{t[0]}, items[1:]...)...)
			}
		}
		go func() {
			defer wg.Done()
			convey.Convey(items...)
		}()
	}, wg.Wait
}
