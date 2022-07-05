package scheduler

import "time"

func getNextDateWithZeroTime() time.Time {
	today := time.Now()
	tmrw := today.AddDate(0, 0, 1)
	year, month, day := tmrw.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, tmrw.Location())

}
