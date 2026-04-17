package appointment

import (
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
)

// MaxAppointmentsPerDay is the workshop daily capacity (non-cancelled appointments).
const MaxAppointmentsPerDay = 8

// validateWorkshopClock checks local wall-clock time is within:
// 09:30–12:30 or 14:00–17:30 (inclusive endpoints).
func validateWorkshopClock(t time.Time) error {
	loc := time.Local
	d := t.In(loc)
	m := d.Hour()*60 + d.Minute()
	const (
		morningStart   = 9*60 + 30
		morningEnd     = 12*60 + 30
		afternoonStart = 14 * 60
		afternoonEnd   = 17*60 + 30
	)
	if (m >= morningStart && m <= morningEnd) || (m >= afternoonStart && m <= afternoonEnd) {
		return nil
	}
	return domain.ErrAppointmentOutsideBusinessHours
}

// dayRangeUTC returns [start, end) in UTC for the calendar day of t in local timezone.
func dayRangeUTC(t time.Time) (startUTC, endUTC time.Time) {
	loc := time.Local
	d := t.In(loc)
	startLocal := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
	endLocal := startLocal.Add(24 * time.Hour)
	return startLocal.UTC(), endLocal.UTC()
}
