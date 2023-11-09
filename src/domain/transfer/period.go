package transfer

import (
	"errors"
	"time"
)

var (
	ErrMoreThanYearInFuture    = errors.New("you cannot search more than a year in the future")
	ErrEndDateBeforeStart      = errors.New("end date cannot be before start date")
	ErrPeriodBiggerThan3Months = errors.New("you cannot search a period bigger than tree months")
)

type Period struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func NewSearchDate(s time.Time, e time.Time) (Period, error) {
	start := s.UTC()
	end := e.UTC()
	now := time.Now().UTC()
	year := now.AddDate(1, 0, 0).UTC()
	endMax3months := s.AddDate(0, 3, 0).UTC()

	if end.After(endMax3months) {
		return Period{}, ErrPeriodBiggerThan3Months
	}

	if start.After(year) || end.After(year) {
		return Period{}, ErrMoreThanYearInFuture
	}

	if end.Before(start) {
		return Period{}, ErrEndDateBeforeStart
	}

	return Period{
		Start: start,
		End:   end,
	}, nil
}
