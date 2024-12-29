package utils

import "time"

func AdjustToEcuadorTime(t time.Time) time.Time {
	loc, err := time.LoadLocation("America/Guayaquil")
	if err != nil {
		return t
	}
	return t.In(loc)
}
