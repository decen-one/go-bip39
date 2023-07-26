package assert

import (
	"testing"
)

func TestTrue(t *testing.T) {
	True(t, true)
	True(t, "" == "")
}

func TestFalse(t *testing.T) {
	False(t, false)
	False(t, 1 == 2)
}
func TestNil(t *testing.T) {
	Nil(t, nil)
}

func TestNotNil(t *testing.T) {
	NotNil(t, "hello")
	NotNil(t, "")
}

func TestEqualString(t *testing.T) {
	EqualString(t, "hello", "hello")
	EqualString(t, "", "")
	EqualString(t, "hello world", "hello"+" world")
}

func TestEqual(t *testing.T) {
	Equal(t, 12, 12)
	Equal(t, 11+2, 13)
}

func TestEqualStringsSlices(t *testing.T) {
	EqualStringsSlices(t, []string{"hello", "world"}, []string{"hello", "world"})
	EqualStringsSlices(t, []string{"", "1"}, []string{"", "1"})
}

func TestEqualByteSlices(t *testing.T) {
	EqualByteSlices(t, []byte{0, 1}, []byte{0, 1})
	EqualByteSlices(t, []byte{12, 11}, []byte{12, 11})

}
