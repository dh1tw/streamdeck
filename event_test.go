package streamdeck

import (
	"testing"

	"go.viam.com/test"
)

func TestState(t *testing.T) {
	s := State{}

	// button 0
	myEvent, err := s.Update(nil, []byte{1, 0, 8, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.Keys[0], test.ShouldBeTrue)
	test.That(t, s.Keys[1], test.ShouldBeFalse)

	myEvent, err = s.Update(nil, []byte{1, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.Keys[0], test.ShouldBeFalse)
	test.That(t, s.Keys[1], test.ShouldBeFalse)
	test.That(t, myEvent.String(), test.ShouldEqual, "key-unpush:0")

	myEvent, err = s.Update(nil, []byte{1, 0, 8, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.Keys[0], test.ShouldBeFalse)
	test.That(t, s.Keys[1], test.ShouldBeTrue)
	test.That(t, myEvent.String(), test.ShouldEqual, "key-push:1")

	// dial turns
	myEvent, err = s.Update(nil, []byte{1, 3, 5, 0, 1, 2})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 52)
	test.That(t, myEvent.String(), test.ShouldEqual, "dial-turn:0")

	myEvent, err = s.Update(nil, []byte{1, 3, 5, 0, 1, 1})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 53)

	myEvent, err = s.Update(nil, []byte{1, 3, 5, 0, 1, 50})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 100)

	myEvent, err = s.Update(nil, []byte{1, 3, 5, 0, 1, 255})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 99)

	myEvent, err = s.Update(nil, []byte{1, 3, 5, 0, 1, 254})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 97)
	test.That(t, myEvent.String(), test.ShouldEqual, "dial-turn:0")

	s.DialPos[0] = 0
	myEvent, err = s.Update(nil, []byte{1, 3, 5, 0, 1, 255})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.DialPos[0], test.ShouldEqual, 0)
	test.That(t, myEvent.String(), test.ShouldEqual, "dial-turn:0")

	// dial pushes
	myEvent, err = s.Update(nil, []byte{1, 3, 5, 0, 0, 1})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.DialPush[0], test.ShouldBeTrue)
	test.That(t, myEvent.String(), test.ShouldEqual, "dial-push:0")

	myEvent, err = s.Update(nil, []byte{1, 3, 5, 0, 0, 0})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.DialPush[0], test.ShouldBeFalse)
	test.That(t, myEvent.String(), test.ShouldEqual, "dial-unpush:0")

}

func TestEventString(t *testing.T) {
	test.That(t, Event{EventDialPressed, 5}.String(), test.ShouldEqual, "dial-push:5")
}

func TestStateOriginal(t *testing.T) {
	s := State{}

	// button 0
	myEvent, err := s.Update(&Original, []byte{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	t.Logf("myEvent 1: %v", myEvent)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.Keys[0], test.ShouldBeTrue)
	test.That(t, s.Keys[1], test.ShouldBeFalse)
	test.That(t, myEvent.String(), test.ShouldEqual, "key-push:0")

	myEvent, err = s.Update(&Original, []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	t.Logf("myEvent 2: %v", myEvent)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, s.Keys[0], test.ShouldBeFalse)
	test.That(t, s.Keys[1], test.ShouldBeFalse)
	test.That(t, myEvent.String(), test.ShouldEqual, "key-unpush:0")
}
