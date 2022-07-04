package appointmentslotsdto

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Parse input dates string and get dates
func validateAndGetDates(dates []string) ([]primitive.DateTime, *map[string]string) {
	errs := map[string]string{}
	var dateList []primitive.DateTime = make([]primitive.DateTime, 0)

	visited := make(map[primitive.DateTime]bool)

	for i := 0; i < len(dates); i++ {
		t, err := time.Parse("2006-01-02", dates[i])
		if err != nil {
			errs[dates[i]] = strings.Replace(err.Error(), "time", "date", 1)
			continue
		}

		curr := time.Now()

		if t.Unix() < curr.Unix() {
			errs[dates[i]] = "Date should not be in the past or today"
			continue
		}

		date := primitive.NewDateTimeFromTime(t)

		// Check Duplicates
		if visited[date] {
			errs[dates[i]] = "Date is found as duplicate"
		} else {
			visited[date] = true
		}

		dateList = append(dateList, date)
	}

	return dateList, &errs

}

// Takes slot range and generate dates
func validateExcludeDaysAndGetDates(excludeDays []string, noOfDays uint16) ([]primitive.DateTime, *map[string]string) {
	errs := map[string]string{}
	var dateList []primitive.DateTime = make([]primitive.DateTime, 0)

	skipDays, err := skipDatesMap(excludeDays)
	if err != nil {
		errs["exclude_days"] = err.Error()
	}

	for i := 1; i < int(noOfDays+1); i++ {

		t := time.Now().AddDate(0, 0, i)

		nt, _ := time.Parse("2006-01-02", t.Format("2006-01-02"))

		if skipDays[t.Weekday()] {
			continue
		}

		date := primitive.NewDateTimeFromTime(nt)
		dateList = append(dateList, date)

	}

	if len(excludeDays) >= 7 {
		errs["exclude_days"] = "All days in a week cannot be excluded"
	}

	return dateList, &errs

}

// Convert Array to HashTable
func skipDatesMap(days []string) (map[time.Weekday]bool, error) {

	skipMap := map[time.Weekday]bool{}
	var err error

	for _, val := range days {
		switch val {
		case "Sunday", "Sun":
			skipMap[time.Sunday] = true

		case "Saturday", "Sat":
			skipMap[time.Saturday] = true

		case "Monday", "Mon":
			skipMap[time.Monday] = true

		case "Tuesday", "Tue":
			skipMap[time.Tuesday] = true

		case "Wednesday", "Wed":
			skipMap[time.Wednesday] = true

		case "Thursday", "Thu":
			skipMap[time.Thursday] = true

		case "Friday", "Fri":
			skipMap[time.Friday] = true

		default:
			err = errors.New("provided excluded day is invalid")
		}
	}

	return skipMap, err

}
