package streamdeck

import (
	"testing"

	"go.viam.com/test"
)

func TestFixKey(t *testing.T) {
	test.That(t, Plus.fixKey(5), test.ShouldEqual, 5)
	test.That(t, Original2.fixKey(5), test.ShouldEqual, 5)

	test.That(t, Original.fixKey(4), test.ShouldEqual, 1)
	test.That(t, Original.fixKey(5), test.ShouldEqual, 10)
}
