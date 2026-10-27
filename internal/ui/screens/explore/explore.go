// Package explore is the searchable, filterable list of every chip.
package explore

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/lazylist"
	"github.com/mehran9675/giode/components/segmented"

	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/store"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/env"
)

// headerKey identifies the filter header, the list's first row.
const headerKey = "explore-header"

// Screen is the Explore tab.
type Screen struct {
	env    *env.Env
	search *components.SearchBar
	sort   *segmented.Control
	list   *lazylist.List

	query   store.Query
	results []*model.Chip
	ranFor  store.Query
	ranOn   *model.Catalog
}

// New returns the Explore screen.
func New(e *env.Env) *Screen {
	s := &Screen{env: e}
	s.search = components.NewSearchBar("Search name, family, codename…", func(string) {})
	s.sort = giode.Input.Segmented(store.SortLabels...).OnChange(func(i int) { s.query.Sort = store.SortOrder(i) })
	s.sort.Styles(giode.Styles{Width: -1})
	s.list = giode.LazyList(
		func() int { return len(s.results) + 2 },
		s.key,
		s.row,
		giode.Styles{Gap: 10, ScrollBar: giode.ScrollBar{HiddenMaxWidth: 600}},
	)
	return s
}

// Apply replaces the filters with q and scrolls to the top.
func (s *Screen) Apply(q store.Query) {
	s.query = q
	s.search.Field.SetText(q.Text)
	s.sort.SetSelected(int(q.Sort))
	s.list.ScrollToKey(headerKey, giode.LazyAlignStart)
}

// View lays the screen out.
func (s *Screen) View() giode.Element {
	s.query.Text = s.search.Text()
	cat := s.env.Store.Catalog.Get()
	if s.ranOn != cat || s.ranFor != s.query {
		s.results = s.query.Run(cat)
		s.ranOn, s.ranFor = cat, s.query
	}
	return giode.Box(giode.Styles{Width: -1, Height: -1, Justify: giode.JustifyCenter},
		giode.Box(giode.Styles{Width: -1, Height: -1, MaxWidth: components.MaxContentWidth}, s.list))
}

func (s *Screen) key(i int) any {
	switch {
	case i == 0:
		return headerKey
	case i > len(s.results):
		return "explore-footer"
	}
	return s.results[i-1]
}

func (s *Screen) row(i int) giode.Element {
	pad := func(el giode.Element) giode.Element {
		return giode.Box(giode.Styles{Width: -1, Padding: giode.SymmetricInset(20, 0)}, el)
	}
	cat := s.env.Store.Catalog.Get()
	switch {
	case i == 0:
		return pad(s.header(cat))
	case i > len(s.results):
		if len(s.results) == 0 {
			return pad(components.EmptyState("search", "No chips found", "Try a different name or remove a filter.", nil))
		}
		return giode.Box(giode.Styles{Height: 24})
	}
	chip := s.results[i-1]
	id := chip.ID
	return pad(components.ChipRow("explore", cat, chip, func() { s.env.Routes.OpenChip(id) }))
}
