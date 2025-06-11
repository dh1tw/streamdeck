package streamdeck

import (
	"testing"

	"go.viam.com/test"
)

func TestState(t *testing.T) {
	s := State{}

	// button 0
	test.That(t, s.Update([]byte{1, 0, 8, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), test.ShouldBeNil)
	test.That(t, s.Keys[0], test.ShouldBeTrue)
	test.That(t, s.Keys[1], test.ShouldBeFalse)
	test.That(t, s.Update([]byte{1, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), test.ShouldBeNil)
	test.That(t, s.Keys[0], test.ShouldBeFalse)
	test.That(t, s.Keys[1], test.ShouldBeFalse)

	// dial turns
	test.That(t, s.Update([]byte{1, 3, 5, 0, 1, 2}), test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 52)

	test.That(t, s.Update([]byte{1, 3, 5, 0, 1, 1}), test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 53)

	test.That(t, s.Update([]byte{1, 3, 5, 0, 1, 50}), test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 100)

	test.That(t, s.Update([]byte{1, 3, 5, 0, 1, 255}), test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 99)

	test.That(t, s.Update([]byte{1, 3, 5, 0, 1, 254}), test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 97)

	s.DialPos[0] = 0
	test.That(t, s.Update([]byte{1, 3, 5, 0, 1, 255}), test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 0)

	// dial pushes
	test.That(t, s.Update([]byte{1, 3, 5, 0, 0, 1}), test.ShouldBeNil)
	test.That(t, s.DialPush[0], test.ShouldBeTrue)
	test.That(t, s.Update([]byte{1, 3, 5, 0, 0, 0}), test.ShouldBeNil)
	test.That(t, s.DialPush[0], test.ShouldBeFalse)

}
