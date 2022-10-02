package scalar

import (
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type textUnmarshaler struct {
	val int
}

func (f *textUnmarshaler) UnmarshalText(b []byte) error {
	f.val = len(b)
	return nil
}

func assertParse(t *testing.T, expected interface{}, str string) {
	v := reflect.New(reflect.TypeOf(expected)).Elem()
	err := ParseValue(v, str)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, v.Interface())
	}

	ptr := reflect.New(reflect.PtrTo(reflect.TypeOf(expected))).Elem()
	err = ParseValue(ptr, str)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, ptr.Elem().Interface())
	}

	assert.True(t, CanParse(v.Type()))
	assert.True(t, CanParse(ptr.Type()))
}

func TestParseValue(t *testing.T) {
	// strings
	assertParse(t, "abc", "abc")

	// booleans
	assertParse(t, true, "true")
	assertParse(t, false, "false")

	// integers
	assertParse(t, int(123), "123")
	assertParse(t, int(123), "1_2_3")
	assertParse(t, int8(123), "123")
	assertParse(t, int16(123), "123")
	assertParse(t, int32(123), "123")
	assertParse(t, int64(123), "123")

	// unsigned integers
	assertParse(t, uint(123), "123")
	assertParse(t, uint(123), "1_2_3")
	assertParse(t, byte(123), "123")
	assertParse(t, uint8(123), "123")
	assertParse(t, uint16(123), "123")
	assertParse(t, uint32(123), "123")
	assertParse(t, uint64(123), "123")
	assertParse(t, uintptr(123), "123")
	assertParse(t, rune(123), "123")

	// floats
	assertParse(t, float32(123), "123")
	assertParse(t, float64(123), "123")

	// durations
	assertParse(t, 3*time.Hour+15*time.Minute, "3h15m")

	// IP addresses
	assertParse(t, net.IPv4(1, 2, 3, 4), "1.2.3.4")

	// email addresses
	assertParse(t, mail.Address{Address: "joe@example.com"}, "joe@example.com")

	// MAC addresses
	assertParse(t, net.HardwareAddr("\x01\x23\x45\x67\x89\xab"), "01:23:45:67:89:ab")

	// MAC addresses
	assertParse(t, net.HardwareAddr("\x01\x23\x45\x67\x89\xab"), "01:23:45:67:89:ab")

	// URL
	assertParse(t, url.URL{Scheme: "https", Host: "example.com", Path: "/a/b/c"}, "https://example.com/a/b/c")

	// custom text unmarshaler
	assertParse(t, textUnmarshaler{3}, "abc")
}

func TestParseErrors(t *testing.T) {
	var err error

	// this should fail because the pointer is nil and will not be settable
	var p *int
	err = ParseValue(reflect.ValueOf(p), "123")
	assert.Equal(t, errPtrNotSettable, err)

	// this should fail because the value will not be settable
	var v int
	err = ParseValue(reflect.ValueOf(v), "123")
	assert.Equal(t, errNotSettable, err)

	// this should fail due to a malformed boolean
	var b bool
	err = ParseValue(reflect.ValueOf(&b), "malformed")
	assert.Error(t, err)

	// this should fail due to a malformed boolean
	var i int
	err = ParseValue(reflect.ValueOf(&i), "malformed")
	assert.Error(t, err)

	// this should fail due to a malformed boolean
	var u uint
	err = ParseValue(reflect.ValueOf(&u), "malformed")
	assert.Error(t, err)

	// this should fail due to a malformed boolean
	var f float64
	err = ParseValue(reflect.ValueOf(&f), "malformed")
	assert.Error(t, err)

	// this should fail due to a malformed time duration
	var d time.Duration
	err = ParseValue(reflect.ValueOf(&d), "malfomed")
	assert.Error(t, err)

	// this should fail due to a malformed email address
	var email mail.Address
	err = ParseValue(reflect.ValueOf(&email), "malfomed")
	assert.Error(t, err)

	// this should fail due to a malformed time duration
	var mac net.HardwareAddr
	err = ParseValue(reflect.ValueOf(&mac), "malfomed")
	assert.Error(t, err)

	// this should fail due to a malformed time duration
	var url url.URL
	err = ParseValue(reflect.ValueOf(&url), "$:")
	assert.Error(t, err)

	// this should fail due to an unsupported type
	var x struct{}
	err = ParseValue(reflect.ValueOf(&x), "$")
	assert.Error(t, err)
}

func TestParse(t *testing.T) {
	var v int
	err := Parse(&v, "123")
	require.NoError(t, err)
	assert.Equal(t, 123, v)
}

func TestCanParseReturnsFalse(t *testing.T) {
	var x struct{}
	assert.Equal(t, false, CanParse(reflect.TypeOf(x)))
}
