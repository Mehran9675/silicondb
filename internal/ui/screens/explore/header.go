package explore

import (
	"fmt"

	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/chip"

	"github.com/mehran9675/silicondb/internal/format"
	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/store"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/theme"
)

// header is the title, search field and filters above the results.
func (s *Screen) header(cat *model.Catalog) giode.Element {
	count := format.Number(float64(len(s.results)), 0) + " chips"
	if len(s.results) == 1 {
		count = "1 chip"
	}
	children := []giode.Element{
		components.AppBar(nil, "Explore", count),
		s.search.View(),
		s.categoryChips(cat),
		s.decadeChips(cat),
	}
	if s.query.Maker != "" {
		children = append(children, s.makerChip(cat))
	}
	s.sort.Styles(components.SegmentedStyles())
	children = append(children, s.sort, giode.Box(giode.Styles{Height: 2}))
	return giode.Stack(giode.Styles{Width: -1, Gap: 12}, children...)
}

// categoryChips filters by category.
func (s *Screen) categoryChips(cat *model.Catalog) giode.Element {
	counts := cat.CountByCategory()
	chips := []giode.Element{s.filterChip("cat", "", "All", s.query.Category == "", func() { s.query.Category = "" })}
	for _, c := range cat.Schema.Categories {
		if counts[c.ID] == 0 {
			continue
		}
		id := c.ID
		chips = append(chips, s.filterChip("cat", id, c.Name, s.query.Category == id, func() {
			if s.query.Category == id {
				s.query.Category = ""
			} else {
				s.query.Category = id
			}
		}))
	}
	return giode.Row(giode.Styles{Width: -1, Gap: 8}, chips...)
}

// decadeChips filters by decade of release.
func (s *Screen) decadeChips(cat *model.Catalog) giode.Element {
	var chips []giode.Element
	for _, d := range store.Decades(cat) {
		d := d
		chips = append(chips, s.filterChip("decade", fmt.Sprint(d), fmt.Sprintf("%ds", d), s.query.Decade == d, func() {
			if s.query.Decade == d {
				s.query.Decade = 0
			} else {
				s.query.Decade = d
			}
		}))
	}
	return giode.Row(giode.Styles{Width: -1, Gap: 8}, chips...)
}

// makerChip shows the manufacturer filter, removable.
func (s *Screen) makerChip(cat *model.Catalog) giode.Element {
	c := giode.Keyed("explore-maker", s.query.Maker, func() *chip.Chip {
		return giode.Chip(cat.MakerName(s.query.Maker)).
			Icon(theme.Icon("factory", 14, theme.OnPrimary)).
			Selected(true).
			OnRemove(func() { s.query.Maker = "" })
	})
	return giode.Row(giode.Styles{Gap: 8}, c)
}

// filterChip is a selectable chip, kept per (group, id).
func (s *Screen) filterChip(group, id, label string, selected bool, onClick func()) giode.Element {
	return components.FilterChip("explore-"+group, id, label, "", selected, onClick)
}
