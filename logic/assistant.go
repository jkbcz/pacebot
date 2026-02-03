package logic

import (
	"math"
	"time"

	"github.com/samber/lo"
)

var assistantMilestones = []Milestone{
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-02-03T00:00:00Z")), Milestone: 0},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-02-28T00:00:00Z")), Milestone: 22},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-06-30T00:00:00Z")), Milestone: 66},
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
