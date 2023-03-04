package errors

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// ============================================================================
//  Preparations and helpers
// ============================================================================

var initpc = caller()

type X struct{}

// val returns a Frame pointing to itself.
func (x X) val() Frame {
	return caller()
}

// ptr returns a Frame pointing to itself.
func (x *X) ptr() Frame {
	return caller()
}

// a test version of runtime.Caller that returns a Frame, not a uintptr.
func caller() Frame {
	var pcs [3]uintptr

	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()

	return Frame(frame.PC)
}

// a test version of runtime.Callers that returns a StackTrace, not a []uintptr.
func stackTrace() StackTrace {
	const depth = 8

	var pcs [depth]uintptr

	n := runtime.Callers(1, pcs[:])

	var st stack = pcs[0:n]

	return st.StackTrace()
}

// ============================================================================
//  Tests
// ============================================================================
//  Note: On refactoring, please keep in mind that the error line numbers are
//        hard-coded in the tests. So if you change the code, you may need to
//        update the line numbers in the tests as well.

// ----------------------------------------------------------------------------
//  Frame
// ----------------------------------------------------------------------------

// nolint: funlen // Allow longer than 60 lines as this is a test
func TestFrame_format(t *testing.T) {
	var tests = []struct {
		Frame
		format  string
		wantExp string
	}{
		{
			Frame:   initpc,
			format:  "%s",
			wantExp: "stack_test.go",
		},
		{
			Frame:  initpc,
			format: "%+s",
			wantExp: "github.com/KEINOS/go-errors/errors.init\n" +
				"\t.+/go-errors/errors/stack_test.go",
		},
		{
			Frame:   0,
			format:  "%s",
			wantExp: "unknown",
		},
		{
			Frame:   0,
			format:  "%+s",
			wantExp: "unknown",
		},
		{
			Frame:   initpc,
			format:  "%d",
			wantExp: "15",
		},
		{
			Frame:   0,
			format:  "%d",
			wantExp: "0",
		},
		{
			Frame:   initpc,
			format:  "%n",
			wantExp: "init",
		},
		{
			Frame: func() Frame {
				var x X
				return x.ptr()
			}(),
			format:  "%n",
			wantExp: `\(\*X\).ptr`,
		},
		{
			Frame: func() Frame {
				var x X
				return x.val()
			}(),
			format:  "%n",
			wantExp: "X.val",
		},
		{
			Frame:   0,
			format:  "%n",
			wantExp: "",
		},
		{
			Frame:   initpc,
			format:  "%v",
			wantExp: "stack_test.go:15",
		},
		{
			Frame:  initpc,
			format: "%+v",
			wantExp: "github.com/KEINOS/go-errors/errors.init\n" +
				"\t.+/go-errors/errors/stack_test.go:15",
		},
		{
			Frame:   0,
			format:  "%v",
			wantExp: "unknown:0",
		},
	}

	for index, tt := range tests {
		testFormatRegexp(t, index, tt.Frame, tt.format, tt.wantExp)
	}
}

// ----------------------------------------------------------------------------
//  funcname
// ----------------------------------------------------------------------------

func TestFuncname(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"github.com/KEINOS/go-errors/errors.funcname", "funcname"},
		{"funcname", "funcname"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}

	for _, tt := range tests {
		expect := tt.want
		actual := funcname(tt.name)

		require.Equal(t, expect, actual)
	}
}

// ----------------------------------------------------------------------------
//  StackTrace
// ----------------------------------------------------------------------------

// nolint: funlen // Allow longer than 60 lines as this is a test
func TestStackTrace(t *testing.T) {
	tests := []struct {
		err      error
		wantExps []string
	}{
		{
			err: New("ooh"),
			wantExps: []string{
				"github.com/KEINOS/go-errors/errors.TestStackTrace\n" +
					"\t.+/go-errors/errors/stack_test.go:186",
			},
		},
		{
			err: Wrap(New("ooh"), "ahh"),
			wantExps: []string{
				"github.com/KEINOS/go-errors/errors.TestStackTrace\n" +
					"\t.+/go-errors/errors/stack_test.go:193", // this is the stack of Wrap, not New
			},
		},
		{
			err: Cause(Wrap(New("ooh"), "ahh")),
			wantExps: []string{
				"github.com/KEINOS/go-errors/errors.TestStackTrace\n" +
					"\t.+/go-errors/errors/stack_test.go:200", // this is the stack of New
			},
		},
		{
			err: func() error { return New("ooh") }(),
			wantExps: []string{
				`github.com/KEINOS/go-errors/errors.TestStackTrace.func1` +
					"\n\t.+/go-errors/errors/stack_test.go:207", // this is the stack of New
				"github.com/KEINOS/go-errors/errors.TestStackTrace\n" +
					"\t.+/go-errors/errors/stack_test.go:207", // this is the stack of New's caller
			},
		},
		{
			err: Cause(func() error {
				return func() error {
					return Errorf("hello %s", fmt.Sprintf("world: %s", "ooh"))
				}()
			}()),
			wantExps: []string{
				`github.com/KEINOS/go-errors/errors.TestStackTrace.func2.1` +
					"\n\t.+go-errors/errors/stack_test.go:218", // this is the stack of Errorf
				`github.com/KEINOS/go-errors/errors.TestStackTrace.func2` +
					"\n\t.+go-errors/errors/stack_test.go:219", // this is the stack of Errorf's caller
				"github.com/KEINOS/go-errors/errors.TestStackTrace\n" +
					"\t.+go-errors/errors/stack_test.go:220", // this is the stack of Errorf's caller's caller
			},
		},
	}

	for i, tt := range tests {
		x, ok := tt.err.(interface {
			StackTrace() StackTrace
		})
		require.True(t, ok, "expected %#v to implement StackTrace() StackTrace", tt.err)

		st := x.StackTrace()
		for j, wantExp := range tt.wantExps {
			testFormatRegexp(t, i, st[j], "%+v", wantExp)
		}
	}
}

// nolint: funlen // Allow longer than 60 lines as this is a test
func TestStackTrace_format(t *testing.T) {
	tests := []struct {
		StackTrace
		format  string
		wantExp string
	}{
		{
			StackTrace: nil,
			format:     "%s",
			wantExp:    `\[\]`,
		},
		{
			StackTrace: nil,
			format:     "%v",
			wantExp:    `\[\]`,
		},
		{
			StackTrace: nil,
			format:     "%+v",
			wantExp:    "",
		},
		{
			StackTrace: nil,
			format:     "%#v",
			wantExp:    `\[\]errors.Frame\(nil\)`,
		},
		{
			StackTrace: make(StackTrace, 0),
			format:     "%s",
			wantExp:    `\[\]`,
		},
		{
			StackTrace: make(StackTrace, 0),
			format:     "%v",
			wantExp:    `\[\]`,
		},
		{
			StackTrace: make(StackTrace, 0),
			format:     "%+v",
			wantExp:    "",
		},
		{
			StackTrace: make(StackTrace, 0),
			format:     "%#v",
			wantExp:    `\[\]errors.Frame{}`,
		},
		{
			StackTrace: stackTrace()[:2],
			format:     "%s",
			wantExp:    `\[stack_test.go stack_test.go\]`,
		},
		{
			StackTrace: stackTrace()[:2],
			format:     "%v",
			wantExp:    `\[stack_test.go:46 stack_test.go:298\]`,
		},
		{
			StackTrace: stackTrace()[:2],
			format:     "%+v",
			wantExp: "\n" +
				"github.com/KEINOS/go-errors/errors.stackTrace\n" +
				"\t.+go-errors/errors/stack_test.go:46\n" +
				"github.com/KEINOS/go-errors/errors.TestStackTrace_format\n" +
				"\t.+go-errors/errors/stack_test.go:303",
		},
		{
			StackTrace: stackTrace()[:2],
			format:     "%#v",
			wantExp:    `\[\]errors.Frame{stack_test.go:46, stack_test.go:312}`,
		},
	}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.StackTrace, tt.format, tt.wantExp)
	}
}
