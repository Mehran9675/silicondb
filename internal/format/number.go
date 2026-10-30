// Package format turns spec values and dates into display text.
package format

import (
	"strconv"
	"strings"
)

// Number formats v with at most decimals fraction digits, trailing zeros
// removed and thousands grouped.
func Number(v float64, decimals int) string {
	s := strconv.FormatFloat(v, 'f', decimals, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	whole, frac, hasFrac := strings.Cut(s, ".")
	group, decimal := separators()
	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteString(group)
		}
		b.WriteRune(r)
	}
	if hasFrac {
		b.WriteString(decimal)
		b.WriteString(frac)
	}
	return b.String()
}

// Count formats a large count in words beyond a million: 3,500,
// 1.2 million, 92 billion.
func Count(v float64) string {
	switch {
	case v >= 1e12:
		return Number(v/1e12, 2) + " " + text("number.trillion", "trillion")
	case v >= 1e9:
		return Number(v/1e9, 2) + " " + text("number.billion", "billion")
	case v >= 1e6:
		return Number(v/1e6, 2) + " " + text("number.million", "million")
	}
	return Number(v, 2)
}
