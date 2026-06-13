package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"time"

	_ "embed"

	"github.com/samber/lo"
)

type Milestone struct {
	Date      time.Time
	Milestone float64
}

type Target struct {
	PersonID     int     `json:"personID"`
	Amount       float64 `json:"amount"`
	ArchivedDate *string `json:"archivedDate"`
}

type Contribution struct {
	PersonID         int
	Amount           float64
	ContributionDate time.Time
}

type Person struct {
	PersonID    int
	DisplayName string
}

var assistantMilestones = []Milestone{
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-02-03T00:00:00Z")), Milestone: 0},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-02-28T00:00:00Z")), Milestone: 22},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-03-17T00:00:00Z")), Milestone: 28},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-03-21T00:00:00Z")), Milestone: 33},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-04-25T00:00:00Z")), Milestone: 44},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-05-23T00:00:00Z")), Milestone: 55},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-06-20T00:00:00Z")), Milestone: 66},
	{Date: lo.Must(time.Parse(time.RFC3339, "2026-07-05T00:00:00Z")), Milestone: 72},
}

//go:embed targets.json
var targetsFile []byte

//go:embed contributions.json
var contributionsFile []byte

//go:embed persons.json
var personsFile []byte

func main() {
	var allTargets []Target
	var contributions []Contribution
	var persons []Person

	err := json.Unmarshal(targetsFile, &allTargets)
	check(err)
	err = json.Unmarshal(contributionsFile, &contributions)
	check(err)
	err = json.Unmarshal(personsFile, &persons)
	check(err)

	var targets []Target
	for _, t := range allTargets {
		if t.ArchivedDate == nil {
			targets = append(targets, t)
		}
	}

	progress := getPersonProgressPercentages(targets, contributions)

	writeProgressCSV("progress.csv", progress, persons)
}

func getPersonProgressPercentages(targets []Target, contributions []Contribution) map[int][]float64 {
	startDate := assistantMilestones[0].Date.UTC().Truncate(24 * time.Hour)
	endDate := assistantMilestones[len(assistantMilestones)-1].Date.UTC().Truncate(24 * time.Hour)
	numDays := int(endDate.Sub(startDate).Hours()/24) + 1

	targetMap := map[int]float64{}
	for _, t := range targets {
		if _, exists := targetMap[t.PersonID]; !exists {
			targetMap[t.PersonID] = t.Amount
		}
	}

	personDayAmounts := map[int]map[int]float64{}
	for _, c := range contributions {
		if _, ok := targetMap[c.PersonID]; !ok {
			continue
		}
		day := c.ContributionDate.UTC().Truncate(24 * time.Hour)
		if day.After(endDate) {
			continue
		}
		dayIdx := 0
		if !day.Before(startDate) {
			dayIdx = int(day.Sub(startDate).Hours() / 24)
		}
		if personDayAmounts[c.PersonID] == nil {
			personDayAmounts[c.PersonID] = map[int]float64{}
		}
		personDayAmounts[c.PersonID][dayIdx] += c.Amount
	}

	result := map[int][]float64{}
	for personID, targetAmount := range targetMap {
		percentages := make([]float64, numDays)
		cumulative := 0.0
		dayAmounts := personDayAmounts[personID]
		for i := 0; i < numDays; i++ {
			cumulative += dayAmounts[i]
			if targetAmount > 0 {
				percentages[i] = cumulative / targetAmount * 100
			}
		}
		result[personID] = percentages
	}

	return result
}

func interpolatedMilestones(numDays int, startDate time.Time) []float64 {
	values := make([]float64, numDays)
	for i := 0; i < numDays; i++ {
		date := startDate.AddDate(0, 0, i)
		ms := assistantMilestones
		if !date.After(ms[0].Date) {
			values[i] = ms[0].Milestone
			continue
		}
		if !date.Before(ms[len(ms)-1].Date) {
			values[i] = ms[len(ms)-1].Milestone
			continue
		}
		for j := 1; j < len(ms); j++ {
			if !date.After(ms[j].Date) {
				t := date.Sub(ms[j-1].Date).Hours() / ms[j].Date.Sub(ms[j-1].Date).Hours()
				values[i] = ms[j-1].Milestone + t*(ms[j].Milestone-ms[j-1].Milestone)
				break
			}
		}
	}
	return values
}

func writeProgressCSV(filename string, progress map[int][]float64, persons []Person) {
	startDate := assistantMilestones[0].Date.UTC().Truncate(24 * time.Hour)

	nameMap := make(map[int]string, len(persons))
	for _, p := range persons {
		nameMap[p.PersonID] = p.DisplayName
	}

	type personRow struct {
		id   int
		name string
	}
	rows := make([]personRow, 0, len(progress))
	for id := range progress {
		name := nameMap[id]
		if name == "" {
			name = strconv.Itoa(id)
		}
		rows = append(rows, personRow{id: id, name: name})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].name < rows[j].name
	})

	numDays := 0
	for _, p := range progress {
		numDays = len(p)
		break
	}

	header := make([]string, 1+numDays)
	header[0] = "Name"
	for i := 0; i < numDays; i++ {
		header[1+i] = startDate.AddDate(0, 0, i).Format("2006-01-02")
	}

	f, err := os.Create(filename)
	check(err)
	defer f.Close()

	w := csv.NewWriter(f)
	check(w.Write(header))

	milestoneRow := make([]string, 1+numDays)
	milestoneRow[0] = "Milestone"
	for i, v := range interpolatedMilestones(numDays, startDate) {
		milestoneRow[1+i] = strconv.FormatFloat(v, 'f', 2, 64)
	}
	check(w.Write(milestoneRow))

	for _, r := range rows {
		row := make([]string, 1+numDays)
		row[0] = r.name
		for i, pct := range progress[r.id] {
			row[1+i] = strconv.FormatFloat(pct, 'f', 2, 64)
		}
		check(w.Write(row))
	}

	w.Flush()
	check(w.Error())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
