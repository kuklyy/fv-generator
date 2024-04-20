package fvtime

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newDate(s string) time.Time {
	t, _ := time.Parse(time.DateOnly, s)
	return t
}

const exampleTime = "2006-01-02"
const exampleDefault = "0001-01-01"

func TestNew(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		nowEnv   string
		expected time.Time
	}{
		{desc: "valid env", nowEnv: exampleTime, expected: newDate(exampleTime)},
		{desc: "invalid format", nowEnv: "11111111", expected: newDate(exampleDefault)},
		{desc: "wrong order", nowEnv: "02-01-2006", expected: newDate(exampleDefault)},
		{desc: "empty", nowEnv: "", expected: newDate(exampleDefault)},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			// Test can't be run in parallel
			os.Setenv(FV_NOW_ENV_NAME, tc.nowEnv)
			defer os.Clearenv()

			result := New(newDate(exampleDefault))

			assert.Equal(t, tc.expected, result)
		})
	}

}
