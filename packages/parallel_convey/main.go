package parallel_convey

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
func New() (*sync.WaitGroup, func(...any)) {
	wg := &sync.WaitGroup{}
	return wg, func(items ...any) {
		wg.Add(1)
		go func() {
			convey.Convey(items...)
			wg.Done()
		}()
	}
}
