package detail

import (
	"sort"

	"github.com/mehran9675/silicondb/internal/format"
	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/ui/components"
)

// section is a spec group with the lines the chip has values for.
type section struct {
	title string
	lines []components.SpecLine
}

// sections groups the chip's specs as the schema describes, with keys the
// schema doesn't know under "Other".
func sections(schema *model.Schema, chip *model.Chip) []section {
	seen := map[string]bool{}
	var out []section
	for _, g := range schema.Groups {
		sec := section{title: g.Label}
		for i := range g.Fields {
			f := &g.Fields[i]
			v, ok := chip.Specs[f.Key]
			if !ok {
				continue
			}
			seen[f.Key] = true
			if s := format.Value(f, v); s != "" {
				sec.lines = append(sec.lines, components.SpecLine{Label: format.Label(f, f.Key), Value: s})
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
	other := section{title: "Other"}
	for _, k := range rest {
		if s := format.Value(nil, chip.Specs[k]); s != "" {
			other.lines = append(other.lines, components.SpecLine{Label: format.Label(nil, k), Value: s})
		}
	}
	if len(other.lines) > 0 {
		out = append(out, other)
	}
	return out
}
