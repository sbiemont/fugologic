package fuzzy

import (
	"testing"

	"fugologic.git/crisp"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewIDVal(t *testing.T) {
	Convey("id val", t, func() {
		fvA := NewIDVal(crisp.Set{})
		fsA1 := NewIDSet(nil, &fvA)

		So(fsA1.parent, ShouldEqual, &fvA)
	})
}
