package assert

import (
	"bytes"
	"testing"
)

func True(t *testing.T, given bool) {
	if !given {
		t.Fatal("Should be true but was false")
	}
}

func False(t *testing.T, given bool) {
	if given {
		t.Fatal("Should be false but was true")
	}
}

func Nil(t *testing.T, given interface{}) {
	if given != nil {
		t.Fatal("Should be nil.")
	}
}

func NotNil(t *testing.T, given interface{}) {
	if given == nil {
		t.Fatal("Should not be nil.")
	}
}
func EqualString(t *testing.T, expected, given string) {
	if expected != given {
		t.Fatal("Strings not equal. Wanted:", expected, "\nGot:", given)
	}
}
func Equal(t *testing.T, a, b interface{}) {
	if a != b {
		t.Errorf("Objects not equal, expected `%s` and got `%s`", a, b)
	}
}

func EqualStringsSlices(t *testing.T, a, b []string) {
	if len(a) != len(b) {
		t.Errorf("String slices not equal, expected %v and got %v", a, b)
		return
	}

	for i := range a {
		if a[i] != b[i] {
			t.Errorf("String slices not equal, expected %v and got %v", a, b)
			return
		}
	}
}

func EqualByteSlices(t *testing.T, a, b []byte) {
	if !bytes.Equal(a, b) {
		t.Errorf("Byte slices not equal, expected %v and got %v", a, b)
		return
	}
}
