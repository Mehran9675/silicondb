package home

import (
	"sort"

	"github.com/mehran9675/giode"

	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/store"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/env"
)

const sectionSize = 5

// recent lists the chips the user opened lately.
func recent(e *env.Env, cat *model.Catalog) giode.Element {
	var chips []*model.Chip
	for _, id := range e.Store.Recent() {
		if c := cat.Chip(id); c != nil && len(chips) < 3 {
			chips = append(chips, c)
		}
	}
	return chipList(e, "recent", "Recently viewed", cat, chips, nil)
}

// latest lists the newest chips.
func latest(e *env.Env, cat *model.Catalog) giode.Element {
	chips := cat.Chips[:min(sectionSize, len(cat.Chips))]
	return chipList(e, "latest", "Latest releases", cat, chips,
		components.TextLink("home-latest", "See all", func() { e.Routes.Explore(store.Query{}) }))
}

// pioneers lists the oldest chips: where it all began.
func pioneers(e *env.Env, cat *model.Catalog) giode.Element {
	n := len(cat.Chips)
	var chips []*model.Chip
	for i := n - 1; i >= 0 && len(chips) < sectionSize; i-- {
		chips = append(chips, cat.Chips[i])
	}
	return chipList(e, "pioneers", "Where it began", cat, chips,
		components.TextLink("home-pioneers", "See all", func() { e.Routes.Explore(store.Query{Sort: store.OldestFirst}) }))
}

// topMakers lists the manufacturers with the most chips.
func topMakers(e *env.Env, cat *model.Catalog) giode.Element {
	makers := append([]*model.Manufacturer(nil), cat.Manufacturers...)
	sort.SliceStable(makers, func(i, j int) bool {
		return len(cat.ChipsBy(makers[i].ID)) > len(cat.ChipsBy(makers[j].ID))
	})
	makers = makers[:min(sectionSize, len(makers))]
	rows := []giode.Element{components.SectionHeader("Top manufacturers",
		components.TextLink("home-makers", "See all", func() { e.Routes.Makers() }))}
	for _, m := range makers {
		id := m.ID
		rows = append(rows, components.MakerRow(m, len(cat.ChipsBy(id)), func() { e.Routes.OpenMaker(id) }))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10}, rows...)
}
