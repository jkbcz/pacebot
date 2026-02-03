package logic

import (
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestAssistantPercentage(t *testing.T) {
	fmt.Println(time.Now())
	type testCase struct {
		time       time.Time
		percentage int
	}

	tcs := []testCase{
		{lo.Must(time.Parse(time.RFC3339, "2026-02-01T10:10:00+01:00")), 0},
		{lo.Must(time.Parse(time.RFC3339, "2026-02-03T10:10:00+01:00")), 0},
		{lo.Must(time.Parse(time.RFC3339, "2026-02-04T10:10:00+01:00")), 1},
		{lo.Must(time.Parse(time.RFC3339, "2026-02-05T10:10:00+01:00")), 2},
		{lo.Must(time.Parse(time.RFC3339, "2026-02-06T10:10:00+01:00")), 3},
		{lo.Must(time.Parse(time.RFC3339, "2026-02-07T10:10:00+01:00")), 4},
		{lo.Must(time.Parse(time.RFC3339, "2026-02-08T10:10:00+01:00")), 4},
		{lo.Must(time.Parse(time.RFC3339, "2026-02-27T10:10:00+01:00")), 21},
		{lo.Must(time.Parse(time.RFC3339, "2026-02-28T10:10:00+01:00")), 22},
		{lo.Must(time.Parse(time.RFC3339, "2026-03-01T01:10:00+01:00")), 22},
		{lo.Must(time.Parse(time.RFC3339, "2026-06-30T01:10:00+01:00")), 66},
	}

	for _, tc := range tcs {
		percentage := getAssistantPercentage(tc.time)
		assert.Equal(t, tc.percentage, percentage)
	}
}
