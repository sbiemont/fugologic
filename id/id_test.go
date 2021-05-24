package id

import (
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewID(t *testing.T) {
	Convey("new id", t, func() {
		identifier := NewID()
		So(identifier, ShouldHaveLength, 36)

		_, err := uuid.Parse(string(identifier))
		So(err, ShouldBeNil)
	})
}
