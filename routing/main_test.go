package routing

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVersionablePrefix(t *testing.T) {
	t.Parallel()

	Convey("should prefix the given module correctly", t, func() {
		prefix := VersionablePrefix("users")
		So(prefix(1), ShouldEqual, "/v1/users")
		So(prefix(2), ShouldEqual, "/v2/users")
	})
}
