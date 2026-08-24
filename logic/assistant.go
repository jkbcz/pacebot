package logic

import (
	"math"
	"time"

	"github.com/samber/lo"
)

var assistantMilestones = []Milestone{
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-02-03T00:00:00Z")), Milestone: 0},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-02-28T00:00:00Z")), Milestone: 22},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-03-17T00:00:00Z")), Milestone: 28},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-03-21T00:00:00Z")), Milestone: 33},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-04-25T00:00:00Z")), Milestone: 44},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-05-23T00:00:00Z")), Milestone: 55},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-06-20T00:00:00Z")), Milestone: 66},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-07-05T00:00:00Z")), Milestone: 72},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-08-01T00:00:00Z")), Milestone: 77},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-09-01T00:00:00Z")), Milestone: 88},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-10-01T00:00:00Z")), Milestone: 100},
}

func getAssistantPercentage(time time.Time) int {

	lastMilestone := assistantMilestones[0]
	for _, m := range assistantMilestones {
		if m.Date.Before(time) {
			lastMilestone = m
		}
	}

	nextMilestone := lastMilestone
	for _, m := range assistantMilestones {
		if m.Date.After(time) {
			nextMilestone = m
			break
		}
	}

	daysBetweenMilestones := math.Floor(nextMilestone.Date.Sub(lastMilestone.Date).Hours() / 24)
	if daysBetweenMilestones == 0 {
		return int(math.Round(lastMilestone.Milestone))
	}

	daysSinceLastMilestone := math.Floor(time.Sub(lastMilestone.Date).Hours() / 24)
	percentageDiff := nextMilestone.Milestone - lastMilestone.Milestone

	return int(math.Round(lastMilestone.Milestone + percentageDiff*daysSinceLastMilestone/daysBetweenMilestones))

}
