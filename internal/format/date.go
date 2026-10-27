package format

import "strconv"

var months = [...]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// Date formats a release date of the form YYYY, YYYY-MM or YYYY-MM-DD as
// "1972", "Apr 1972" or "Apr 1, 1972".
func Date(s string) string {
	if len(s) < 7 {
		return s
	}
	m, err := strconv.Atoi(s[5:7])
	if err != nil || m < 1 || m > 12 {
		return s[:4]
	}
	if len(s) >= 10 {
		if d, err := strconv.Atoi(s[8:10]); err == nil {
			return months[m-1] + " " + strconv.Itoa(d) + ", " + s[:4]
		}
	}
	return months[m-1] + " " + s[:4]
}
