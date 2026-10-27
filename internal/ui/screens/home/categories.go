package home

import (
	"github.com/mehran9675/giode"

	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/store"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/env"
)

// categories is the grid of chip categories with their counts.
func categories(e *env.Env, cat *model.Catalog) giode.Element {
	counts := cat.CountByCategory()
	var tiles []giode.Element
	for _, c := range cat.Schema.Categories {
		if counts[c.ID] == 0 {
			continue
		}
		id := c.ID
		tiles = append(tiles, components.CategoryTile(c, counts[c.ID], func() {
			e.Routes.Explore(store.Query{Category: id})
		}))
	}
	cols := 2
	if !giode.Compact() {
		cols = 4
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10},
		components.SectionHeader("Categories", nil),
		components.Grid(cols, 12, tiles...),
	)
}
