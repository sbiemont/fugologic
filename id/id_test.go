package id

import (
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestID(t *testing.T) {
	Convey("new id", t, func() {
		identifier := NewID()
		So(identifier, ShouldHaveLength, 36)

		_, err := uuid.Parse(string(identifier))
		So(err, ShouldBeNil)
	})

	Convey("empty", t, func() {
		identifier := NewID()
		So(identifier.Empty(), ShouldBeFalse)
		So(ID("").Empty(), ShouldBeTrue)
	})
}
