package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/mehran9675/silicondb/internal/model"
)

// completeness tracks, against schema.json's "expected" lists, which
// specs each chip is missing: neither given nor marked not applicable.
type completeness struct {
	schema  *model.Schema
	prefix  string
	missing map[string][]string
	filled  map[string][2]int // per category: expected keys covered, total
}

func newCompleteness(schema *model.Schema, prefix string) *completeness {
	return &completeness{schema: schema, prefix: prefix, missing: map[string][]string{}, filled: map[string][2]int{}}
}

func (c *completeness) add(id, category string, specs map[string]json.RawMessage, na []string) {
	if c.prefix != "" && !strings.HasPrefix(id, c.prefix) {
		return
	}
	n := c.filled[category]
	for _, k := range c.schema.Expected[category] {
		n[1]++
		if _, ok := specs[k]; ok || slices.Contains(na, k) {
			n[0]++
			continue
		}
		c.missing[id] = append(c.missing[id], k)
	}
	c.filled[category] = n
}

func (c *completeness) print() {
	ids := make([]string, 0, len(c.missing))
	for id := range c.missing {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		fmt.Printf("%s: %s\n", id, strings.Join(c.missing[id], ", "))
	}
	cats := make([]string, 0, len(c.filled))
	for cat := range c.filled {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	for _, cat := range cats {
		n := c.filled[cat]
		if n[1] > 0 {
			fmt.Printf("%-5s %5.1f%% complete\n", cat, 100*float64(n[0])/float64(n[1]))
		}
	}
	fmt.Printf("%d chips with gaps\n", len(ids))
}
