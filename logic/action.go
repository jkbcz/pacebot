package logic

import (
	"time"

	"github.com/samber/lo"
)

type Milestone struct {
	Date      time.Time
	Milestone float64
}

var actionMilestones = []Milestone{
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-02-01T00:00:00Z")), Milestone: 11},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-03-01T00:00:00Z")), Milestone: 22},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-04-01T00:00:00Z")), Milestone: 33},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-05-01T00:00:00Z")), Milestone: 44},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-06-01T00:00:00Z")), Milestone: 55},
	{Date: lo.Must(time.Parse(time.RFC3339, "2025-07-01T00:00:00Z")), Milestone: 66},
}

func getStatusForNextMilestone(now time.Time) float64 {
	for _, m := range actionMilestones {
		if m.Date.After(now) {
			return m.Milestone
		}
	}
	return 100
}
