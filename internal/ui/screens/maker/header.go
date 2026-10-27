package maker

import (
	"strconv"

	"github.com/mehran9675/giode"

	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/theme"
)

// header is the company card, its numbers and the category filter.
func (s *Screen) header(cat *model.Catalog, m *model.Manufacturer) giode.Element {
	all := cat.ChipsBy(m.ID)
	counts := map[string]int{}
	first, last := 0, 0
	for _, c := range all {
		counts[c.Category]++
		if y := c.Year(); y != 0 {
			if first == 0 || y < first {
				first = y
			}
			last = max(last, y)
		}
	}
	span := "—"
	if first != 0 {
		span = strconv.Itoa(first) + "–" + strconv.Itoa(last)
	}
	status := "Active"
	if !m.Active() {
		status = "Until " + strconv.Itoa(m.Defunct)
	}

	chips := []giode.Element{s.filterChip("", "All", len(all))}
	for _, c := range cat.Schema.Categories {
		if counts[c.ID] > 0 {
			chips = append(chips, s.filterChip(c.ID, c.Name, counts[c.ID]))
		}
	}

	return giode.Stack(giode.Styles{Width: -1, Gap: 14, Padding: giode.Inset{Bottom: 4}},
		components.AppBar(s.back, "", ""),
		components.Card(giode.Styles{Gap: 16, Padding: giode.UniformInset(20)},
			giode.Row(giode.Styles{Gap: 14, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				giode.Avatar(m.Name, giode.Styles{Width: 60}),
				giode.Stack(giode.Styles{FlexGrow: 1, Gap: 3},
					giode.Text(m.Name, theme.Title()),
					giode.Text(components.MakerSubtitle(m), theme.Caption().With(giode.Styles{FontSize: 13.5})),
				),
			),
			components.Grid(3, 10,
				components.StatTile("Chips", strconv.Itoa(len(all)), theme.Primary),
				components.StatTile("Releases", span, theme.Accent),
				components.StatTile("Status", status, theme.Success),
			),
		),
		giode.Row(giode.Styles{Width: -1, Gap: 8}, chips...),
	)
}

// filterChip selects a category of this manufacturer's chips.
func (s *Screen) filterChip(id, label string, n int) giode.Element {
	return components.FilterChip("maker/"+s.id, id, label, strconv.Itoa(n), s.category == id, func() { s.category = id })
}
