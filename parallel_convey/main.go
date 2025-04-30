package parallel_convey

import (
	"sync"

	"github.com/smartystreets/goconvey/convey"
)

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
