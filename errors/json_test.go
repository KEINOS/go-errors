package errors

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrameMarshalText(t *testing.T) {
	var tests = []struct {
		Frame
		wantExp string
	}{{
		Frame:   initpc,
		wantExp: `^github.com/KEINOS/go-errors/errors\.init(\.ializers)? .+/go-errors/errors/stack_test.go:\d+$`,
	}, {
		Frame:   0,
		wantExp: `^unknown$`,
	}}

	for index, tt := range tests {
		got, err := tt.Frame.MarshalText()
		require.NoError(t, err, "test #%d: MarshalText failed during test setup", index+1)

		assert.Regexp(t, regexp.MustCompile(tt.wantExp), string(got),
			"test #%d: MarshalText failed:\n got %q\n want %q", index+1, got, tt.wantExp)
	}
}

func TestFrameMarshalJSON(t *testing.T) {
	var tests = []struct {
		Frame
		wantExp string
	}{{
		Frame:   initpc,
		wantExp: `^"github\.com/KEINOS/go-errors/errors\.init(\.ializers)? .+/go-errors/errors/stack_test.go:\d+"$`,
	}, {
		Frame:   0,
		wantExp: `^"unknown"$`,
	}}

	for index, tt := range tests {
		got, err := json.Marshal(tt.Frame)
		require.NoError(t, err, "test #%d: MarshalJSON failed during test setup", index+1)

		assert.Regexp(t, regexp.MustCompile(tt.wantExp), string(got),
			"test #%d: MarshalText failed:\n got %q\n want %q", index+1, got, tt.wantExp)
	}
}
