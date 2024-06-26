package datelib

import (
	"math"
	"time"
)

func Today() string {
	return time.Now().Format(time.DateOnly)
}

func Parse(date string) (time.Time, error) {
	return time.Parse(time.DateOnly, date)
}

func isYearOnly(
	frequencyYears int,
	frequencyMonths int,
	frequencyWeeks int,
	frequencyDays int,
) bool {
	return frequencyYears > 0 && frequencyMonths == 0 && frequencyWeeks == 0 && frequencyDays == 0
}

func isMonthOnly(
	frequencyYears int,
	frequencyMonths int,
	frequencyWeeks int,
	frequencyDays int,
) bool {
	return frequencyYears == 0 && frequencyMonths > 0 && frequencyWeeks == 0 && frequencyDays == 0
}

func PeriodicDates(
	periodicStart time.Time,
	targetStart time.Time,
	targetEnd time.Time,
	frequencyYears int,
	frequencyMonths int,
	frequencyWeeks int,
	frequencyDays int,
) []time.Time {
	diffHours := targetStart.Sub(periodicStart).Hours()
	diffDays := int(math.Floor(diffHours / 24))
	minimumDays := frequencyYears*365 + frequencyMonths*28 + frequencyWeeks*7 + frequencyDays
	offsetDays := (diffDays / minimumDays) * minimumDays

	loopDate := periodicStart.AddDate(0, 0, offsetDays)
	periodicStartNextDate := periodicStart.AddDate(frequencyYears, frequencyMonths, frequencyWeeks*7+frequencyDays)

	if isYearOnly(frequencyYears, frequencyMonths, frequencyWeeks, frequencyDays) {
		loopDate = time.Date(loopDate.Year(), periodicStart.Month(), periodicStart.Day(), 0, 0, 0, 0, loopDate.Location())
	}
	if isMonthOnly(frequencyYears, frequencyMonths, frequencyWeeks, frequencyDays) {
		loopDate = time.Date(loopDate.Year(), loopDate.Month(), periodicStart.Day(), 0, 0, 0, 0, loopDate.Location())
	}
	if loopDate.Before(periodicStartNextDate) {
		loopDate = periodicStart
	}

	dates := []time.Time{}
	for !targetEnd.Before(loopDate) {
		if !loopDate.Before(targetStart) {
			dates = append(dates, loopDate)
		}
		loopDate = loopDate.AddDate(frequencyYears, frequencyMonths, frequencyWeeks*7+frequencyDays)
	}

	return dates
}
