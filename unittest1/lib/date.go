package lib

import "time"

// This function converts "YYYY-MM-DD" to "20 October 2025".
func FormatDateLong(s string) (string, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return "", err
	}
	return t.Format("2 January 2006"), nil
}