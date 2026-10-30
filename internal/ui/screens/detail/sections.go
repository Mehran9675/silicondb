package detail

import (
	"sort"

	"github.com/mehran9675/silicondb/internal/format"
	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/tr"
)

// benchmarkGroup is the schema group shown as bars instead of rows.
const benchmarkGroup = "benchmarks"

// section is a spec group with the lines the chip has values for.
type section struct {
	title string
	lines []components.SpecLine
}

// sections groups the chip's specs as the schema describes, with keys the
// schema doesn't know under "Other". Benchmarks are left to benchmarks.
func sections(schema *model.Schema, chip *model.Chip) []section {
	seen := map[string]bool{}
	var out []section
	for gi := range schema.Groups {
		g := &schema.Groups[gi]
		if g.ID == benchmarkGroup {
			for _, f := range g.Fields {
				seen[f.Key] = true
			}
			continue
		}
		sec := section{title: tr.Group(g)}
		for i := range g.Fields {
			f := &g.Fields[i]
			v, ok := chip.Specs[f.Key]
			if !ok {
				continue
			}
			seen[f.Key] = true
			if s := tr.Value(f, f.Key, v); s != "" {
				sec.lines = append(sec.lines, components.SpecLine{Label: tr.Label(f, f.Key), Value: s})
			}
		}
		if len(sec.lines) > 0 {
			out = append(out, sec)
		}
	}
	var rest []string
	for k := range chip.Specs {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	other := section{title: tr.T("common.other")}
	for _, k := range rest {
		if s := format.Value(nil, chip.Specs[k]); s != "" {
			other.lines = append(other.lines, components.SpecLine{Label: tr.Label(nil, k), Value: s})
		}
	}
	if len(other.lines) > 0 {
		out = append(out, other)
	}
	return out
}

// benchmarks returns the chip's scores with their share of the best score
// in the database.
func benchmarks(cat *model.Catalog, chip *model.Chip) (title string, scores []components.Benchmark) {
	for gi := range cat.Schema.Groups {
		g := &cat.Schema.Groups[gi]
		if g.ID != benchmarkGroup {
			continue
		}
		title = tr.Group(g)
		for i := range g.Fields {
			f := &g.Fields[i]
			v, ok := chip.Number(f.Key)
			if !ok {
				continue
			}
			frac := float32(1)
			if m := cat.Max(f.Key); m > 0 {
				frac = float32(v / m)
			}
			scores = append(scores, components.Benchmark{Label: tr.Label(f, f.Key), Score: format.Value(f, v), Fraction: frac})
		}
	}
	return title, scores
}
