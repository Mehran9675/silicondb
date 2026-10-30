package store

import (
	"sort"
	"strings"

	"github.com/mehran9675/silicondb/internal/model"
)

// SortOrder orders query results.
type SortOrder int

const (
	NewestFirst SortOrder = iota
	OldestFirst
	ByName
	// Fastest orders by benchmark score; chips without one follow,
	// newest first.
	Fastest
)

// SortKeys are the message keys naming the sort orders, in order.
var SortKeys = []string{"sort.newest", "sort.oldest", "sort.name", "sort.fastest"}

// SpeedKey is the benchmark Fastest sorts by.
const SpeedKey = "geekbench6_multi"

// Query selects chips from the catalog. Zero fields match everything.
type Query struct {
	Text     string
	Category string
	Maker    string
	// Decade is the first year of a decade (1970, 1980, ...), or 0.
	Decade int
	Sort   SortOrder
}

// Run returns the chips matching q.
func (q Query) Run(c *model.Catalog) []*model.Chip {
	if c == nil {
		return nil
	}
	words := strings.Fields(strings.ToLower(q.Text))
	src := c.Chips
	if q.Maker != "" {
		src = c.ChipsBy(q.Maker)
	}
	out := make([]*model.Chip, 0, len(src))
	for _, chip := range src {
		if q.Category != "" && chip.Category != q.Category {
			continue
		}
		if q.Decade != 0 {
			if y := chip.Year(); y < q.Decade || y >= q.Decade+10 {
				continue
			}
		}
		if len(words) > 0 && !chip.Matches(words) {
			continue
		}
		out = append(out, chip)
	}
	switch q.Sort {
	case OldestFirst:
		sort.SliceStable(out, func(i, j int) bool { return out[i].Released < out[j].Released })
	case ByName:
		sort.SliceStable(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	case Fastest:
		sort.SliceStable(out, func(i, j int) bool {
			a, okA := out[i].Number(SpeedKey)
			b, okB := out[j].Number(SpeedKey)
			if okA != okB {
				return okA
			}
			return a > b
		})
	}
	return out
}

// Decades lists the decades the catalog has chips from.
func Decades(c *model.Catalog) []int {
	first, last := c.YearSpan()
	if first == 0 {
		return nil
	}
	var ds []int
	for d := first / 10 * 10; d <= last; d += 10 {
		ds = append(ds, d)
	}
	return ds
}
