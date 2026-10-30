// Package maker shows one manufacturer and its chips.
package maker

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/lazylist"

	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/env"
	"github.com/mehran9675/silicondb/internal/ui/tr"
)

// Screen is a manufacturer's page.
type Screen struct {
	env      *env.Env
	id       string
	back     *components.RoundButton
	list     *lazylist.List
	category string
	chips    []*model.Chip
}

// New returns the page of the manufacturer with id.
func New(e *env.Env, id string) *Screen {
	s := &Screen{env: e, id: id}
	s.back = components.NewRoundButton("chevron-left", tr.T("common.back"), e.Routes.Back)
	s.list = giode.LazyList(
		func() int { return len(s.chips) + 1 },
		func(i int) any {
			if i == 0 {
				return "maker-header"
			}
			return s.chips[i-1]
		},
		s.row,
		giode.Styles{Gap: 10, ScrollBar: giode.ScrollBar{HiddenMaxWidth: 600}},
	)
	return s
}

// View lays the screen out.
func (s *Screen) View() giode.Element {
	cat := s.env.Store.Catalog.Get()
	s.chips = s.chips[:0]
	for _, c := range cat.ChipsBy(s.id) {
		if s.category == "" || c.Category == s.category {
			s.chips = append(s.chips, c)
		}
	}
	return giode.Box(giode.Styles{Width: -1, Height: -1, Justify: giode.JustifyCenter},
		giode.Box(giode.Styles{Width: -1, Height: -1, MaxWidth: components.MaxContentWidth}, s.list))
}

func (s *Screen) row(i int) giode.Element {
	cat := s.env.Store.Catalog.Get()
	pad := func(el giode.Element) giode.Element {
		return giode.Box(giode.Styles{Width: -1, Padding: giode.SymmetricInset(20, 0)}, el)
	}
	if i == 0 {
		m := cat.Manufacturer(s.id)
		if m == nil {
			return pad(giode.Stack(giode.Styles{Width: -1},
				components.AppBar(s.back, tr.T("common.not_found"), ""),
				components.EmptyState("factory", tr.T("maker.missing.title"), tr.T("maker.missing.message"), nil)))
		}
		return pad(s.header(cat, m))
	}
	chip := s.chips[i-1]
	id := chip.ID
	return pad(components.ChipRow("maker", cat, chip, func() { s.env.Routes.OpenChip(id) }))
}
