package fuzzy

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOperator(t *testing.T) {
	Convey("operator zadeh", t, func() {
		So(OperatorZadeh{}.And(42, 43), ShouldEqual, 42)
		So(OperatorZadeh{}.Or(42, 43), ShouldEqual, 43)
		So(OperatorZadeh{}.XOr(42, 43), ShouldEqual, 1) // 42+43-2*min(42,43)
	})

	Convey("operator hyperbolic", t, func() {
		So(OperatorHyperbolic{}.And(42, 43), ShouldEqual, 1806)  // 42*43
		So(OperatorHyperbolic{}.Or(42, 43), ShouldEqual, -1721)  // 42+43-42*43
		So(OperatorHyperbolic{}.XOr(42, 43), ShouldEqual, -3527) // 42+43-2*42*43
	})
}
