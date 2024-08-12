package timer

import "time"

func ParseTime(timeStr string) (time.Time, error) {
	now := time.Now()
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location()), nil
}
