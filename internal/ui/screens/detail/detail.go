// Package detail shows every specification of one chip.
package detail

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/scroll"

	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/env"
	"github.com/mehran9675/silicondb/internal/ui/theme"
)

// Screen is a chip's specification sheet.
type Screen struct {
	env     *env.Env
	id      string
	scroll  *scroll.Scroll
	back    *components.RoundButton
	compare *components.RoundButton
}

// New returns the screen for the chip with id.
func New(e *env.Env, id string) *Screen {
	s := &Screen{env: e, id: id}
	s.scroll = giode.Scroll(s.content, giode.Styles{ScrollBar: giode.ScrollBar{HiddenMaxWidth: 600}})
	s.back = components.NewRoundButton("chevron-left", "Back", e.Routes.Back)
	s.compare = components.NewRoundButton("compare", "Compare", s.toggleCompare)
	e.Store.Viewed(id)
	return s
}

// View lays the screen out.
func (s *Screen) View() giode.Element { return s.scroll }

func (s *Screen) toggleCompare() {
	if s.env.Store.ToggleCompare(s.id) {
		s.env.Toasts.Show(giode.ToastSuccess, "Added to Compare")
	} else {
		s.env.Toasts.Show(giode.ToastInfo, "Removed from Compare")
	}
}

func (s *Screen) content() giode.Element {
	cat := s.env.Store.Catalog.Get()
	chip := cat.Chip(s.id)
	if chip == nil {
		return components.Column(16,
			components.AppBar(s.back, "Not found", ""),
			components.EmptyState("alert", "This chip is no longer in the database", "It may have been renamed in the latest update.", nil),
		)
	}
	if s.env.Store.InCompare(s.id) {
		s.compare.Icon("check", theme.Primary)
	} else {
		s.compare.Icon("plus", "")
	}
	children := []giode.Element{
		components.AppBar(s.back, "", "", s.compare),
		hero(s.env, cat, chip),
		highlights(cat, chip),
	}
	for _, sec := range sections(cat.Schema, chip) {
		children = append(children, components.SpecSection(sec.title, sec.lines))
	}
	children = append(children, related(s.env, cat, chip))
	return components.Column(16, children...)
}

// related lists other chips of the same family and manufacturer.
func related(e *env.Env, cat *model.Catalog, chip *model.Chip) giode.Element {
	rows := []giode.Element{components.SectionHeader("More in "+chip.Family, nil)}
	for _, c := range cat.ChipsBy(chip.Manufacturer) {
		if c == chip || c.Family != chip.Family || len(rows) > 5 {
			continue
		}
		id := c.ID
		rows = append(rows, components.ChipRow("related", cat, c, func() { e.Routes.OpenChip(id) }))
	}
	if len(rows) == 1 {
		return nil
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10}, rows...)
}
