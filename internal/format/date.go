package format

import "strconv"

var months = [...]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// Date formats a release date of the form YYYY, YYYY-MM or YYYY-MM-DD as
// "1972", "Apr 1972" or "Apr 1, 1972", in the current language.
func Date(s string) string {
	if len(s) < 7 {
		return s
	}
	m, err := strconv.Atoi(s[5:7])
	if err != nil || m < 1 || m > 12 {
		return s[:4]
	}
	month := text("month."+strconv.Itoa(m), months[m-1])
	if len(s) >= 10 {
		if d, err := strconv.Atoi(s[8:10]); err == nil {
			return text("date.full", "{month} {day}, {year}", "month", month, "day", strconv.Itoa(d), "year", s[:4])
		}
	}
	return text("date.month_year", "{month} {year}", "month", month, "year", s[:4])
}
