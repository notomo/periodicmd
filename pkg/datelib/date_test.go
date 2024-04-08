package datelib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func dates(ts []time.Time) []string {
	strs := []string{}
	for _, t := range ts {
		strs = append(strs, t.Format(time.DateOnly))
	}
	return strs
}

func TestPeriodicDates(t *testing.T) {
	t.Run("returns empty if time range is zero", func(t *testing.T) {
		now := time.Now()
		periodicStart := now
		targetStart := now
		targetEnd := now.Add(-1 * time.Second)

		got := PeriodicDates(
			periodicStart,
			targetStart,
			targetEnd,
			0,
			0,
			1,
			0,
		)

		want := []time.Time{}
		assert.Equal(t, want, got)
	})

	t.Run("returns times included by time range", func(t *testing.T) {
		periodicStart := time.Now()
		targetStart := periodicStart
		targetEnd := targetStart.Add(14 * 24 * time.Hour)

		got := PeriodicDates(
			periodicStart,
			targetStart,
			targetEnd,
			0,
			0,
			1,
			0,
		)

		want := []time.Time{
			targetStart,
			targetStart.Add(7 * 24 * time.Hour),
			targetStart.Add(14 * 24 * time.Hour),
		}
		assert.Equal(t, dates(want), dates(got))
	})
}
