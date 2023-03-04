package errors

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		err    string
		expect error
	}{
		{"", fmt.Errorf("")},
		{"foo", fmt.Errorf("foo")},
		{"foo", New("foo")},
		{"string with format specifiers: %v", errors.New("string with format specifiers: %v")},
	}

	for _, tt := range tests {
		expect := tt.expect.Error()
		actual := New(tt.err).Error()

		require.Equal(t, expect, actual)
	}
}

func TestWrapNil(t *testing.T) {
	got := Wrap(nil, "no error")

	require.Nil(t, got, "if err is nil, Wrap should return nil")
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err     error
		message string
		expect  string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrap(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		expect := tt.expect
		actual := Wrap(tt.err, tt.message).Error()

		require.Equal(t, expect, actual)
	}
}

// dummy error type to test against
type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	x := New("error")

	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		// caused error returns cause
		err:  Wrap(io.EOF, "ignored"),
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}, {
		WithMessage(nil, "whoops"),
		nil,
	}, {
		WithMessage(io.EOF, "whoops"),
		io.EOF,
	}, {
		WithStack(nil),
		nil,
	}, {
		WithStack(io.EOF),
		io.EOF,
	}}

	for i, tt := range tests {
		heystack := Cause(tt.err) // Get underlying error
		needle := tt.want

		require.Equal(t, heystack, needle, "test #%d failed: no match in error chain", i+1)
	}
}

func TestWrapfNil(t *testing.T) {
	got := Wrapf(nil, "no error")

	require.Nil(t, got, "if err is nil, Wrapf should return nil")
}

func TestWrapf(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{
			io.EOF,
			"read error",
			"read error: EOF",
		},
		{
			Wrapf(io.EOF, "read error without format specifiers"),
			"client error",
			"client error: read error without format specifiers: EOF",
		},
		{
			Wrapf(io.EOF, "read error with %d format specifier", 1),
			"client error",
			"client error: read error with 1 format specifier: EOF",
		},
	}

	for index, tt := range tests {
		expect := tt.want
		actual := Wrapf(tt.err, tt.message).Error()

		require.Equal(t, expect, actual,
			"test #%d failed: Wrapf(%q, %q)", index, tt.err, tt.message)
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{
			Errorf("read error without format specifiers"),
			"read error without format specifiers",
		},
		{
			Errorf("read error with %d format specifier", 1),
			"read error with 1 format specifier",
		},
	}

	for index, tt := range tests {
		expect := tt.want
		actual := tt.err.Error()

		require.Equal(t, expect, actual, "test #%d failed", index)
	}
}

func TestWithStackNil(t *testing.T) {
	got := WithStack(nil)

	require.Nil(t, got, "if err is nil, WithStack should return nil")
}

func TestWithStack(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{io.EOF, "EOF"},
		{WithStack(io.EOF), "EOF"}, // double stack
	}

	for _, tt := range tests {
		expect := tt.want
		actual := WithStack(tt.err).Error()

		require.Equal(t, expect, actual)
	}
}

func TestWithMessageNil(t *testing.T) {
	got := WithMessage(nil, "no error")

	require.Nil(t, got, "if err is nil, WithMessage should return nil")
}

func TestWithMessage(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessage(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for index, tt := range tests {
		expect := tt.want
		actual := WithMessage(tt.err, tt.message).Error()

		require.Equal(t, expect, actual,
			"test #%d failed: WithMessage(%q, %q).Error() did not return as expected",
			index, tt.err, tt.message)
	}
}

func TestWithMessagefNil(t *testing.T) {
	got := WithMessagef(nil, "no error")

	require.Nil(t, got, "if err is nil, WithMessagef should return nil")
}

func TestWithMessagef(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{
			io.EOF,
			"read error",
			"read error: EOF",
		},
		{WithMessagef(
			io.EOF,
			"read error without format specifier"),
			"client error", "client error: read error without format specifier: EOF",
		},
		{
			WithMessagef(io.EOF, "read error with %d format specifier", 1),
			"client error",
			"client error: read error with 1 format specifier: EOF",
		},
	}

	for _, tt := range tests {
		expect := tt.want
		actual := WithMessagef(tt.err, tt.message).Error()

		require.Equal(t, expect, actual,
			"test #%d failed: WithMessagef(%q, %q).Error() did not return as expected",
			tt.err, tt.message)
	}
}

// errors.New, etc values are not expected to be compared by value
// but the change in errors#27 made them incomparable. Assert that
// various kinds of errors have a functional equality operator, even
// if the result of that equality is always false.
func TestErrorEquality(t *testing.T) {
	vals := []error{
		nil,
		io.EOF,
		errors.New("EOF"),
		New("EOF"),
		Errorf("EOF"),
		Wrap(io.EOF, "EOF"),
		Wrapf(io.EOF, "EOF%d", 2),
		WithMessage(nil, "whoops"),
		WithMessage(io.EOF, "whoops"),
		WithStack(io.EOF),
		WithStack(nil),
	}

	for i := range vals {
		for j := range vals {
			assert.NotPanics(t, func() {
				_ = vals[i] == vals[j]
			}, "comparing %T and %T", vals[i], vals[j])
		}
	}
}

func Test_panicOnWriteErr(t *testing.T) {
	dummyIOWriteString := func() (int, error) {
		return 0, errors.New("forced error")
	}

	require.Panics(t, func() {
		panicOnWriteErr(dummyIOWriteString())
	})

	require.PanicsWithError(t, "forced error (file: errors_test.go, line: 299)", func() {
		panicOnWriteErr(dummyIOWriteString())
	})
}
