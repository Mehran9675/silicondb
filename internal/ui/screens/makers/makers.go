// Package makers lists the manufacturers in the database.
package makers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/scroll"
	"github.com/mehran9675/giode/components/segmented"

	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/env"
)

// Screen is the Manufacturers tab.
type Screen struct {
	env    *env.Env
	scroll *scroll.Scroll
	search *components.SearchBar
	order  *segmented.Control
}

// New returns the Manufacturers screen.
func New(e *env.Env) *Screen {
	s := &Screen{env: e}
	s.search = components.NewSearchBar("Search manufacturers", func(string) {})
	s.order = giode.Input.Segmented("A–Z", "Most chips", "Oldest")
	s.order.Styles(giode.Styles{Width: -1})
	s.scroll = giode.Scroll(s.content, giode.Styles{ScrollBar: giode.ScrollBar{HiddenMaxWidth: 600}})
	return s
}

// View lays the screen out.
func (s *Screen) View() giode.Element { return s.scroll }

func (s *Screen) content() giode.Element {
	cat := s.env.Store.Catalog.Get()
	makers := s.filtered(cat)
	s.order.Styles(components.SegmentedStyles())
	children := []giode.Element{
		components.AppBar(nil, "Manufacturers", fmt.Sprintf("%d companies since 1972", len(cat.Manufacturers))),
		s.search.View(),
		s.order,
	}
	for _, m := range makers {
		id := m.ID
		children = append(children, components.MakerRow(m, len(cat.ChipsBy(id)), func() { s.env.Routes.OpenMaker(id) }))
	}
	if len(makers) == 0 {
		children = append(children, components.EmptyState("factory", "No manufacturers found", "Try another name or country.", nil))
	}
	return components.Column(10, children...)
}

// filtered applies the search text and the order.
func (s *Screen) filtered(cat *model.Catalog) []*model.Manufacturer {
	q := strings.ToLower(strings.TrimSpace(s.search.Text()))
	var out []*model.Manufacturer
	for _, m := range cat.Manufacturers {
		if q == "" || strings.Contains(strings.ToLower(m.Name+" "+m.Country), q) {
			out = append(out, m)
		}
	}
	switch s.order.Selected() {
	case 1:
		sort.SliceStable(out, func(i, j int) bool { return len(cat.ChipsBy(out[i].ID)) > len(cat.ChipsBy(out[j].ID)) })
	case 2:
		sort.SliceStable(out, func(i, j int) bool { return founded(out[i]) < founded(out[j]) })
	}
	return out
}

func founded(m *model.Manufacturer) int {
	if m.Founded == 0 {
		return 1 << 30
	}
	return m.Founded
}
