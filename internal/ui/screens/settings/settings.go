// Package settings holds appearance, database and about options.
package settings

import (
	"fmt"
	"time"

	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/input"
	"github.com/mehran9675/giode/components/scroll"
	"github.com/mehran9675/giode/components/toggle"

	"github.com/mehran9675/silicondb/internal/config"
	"github.com/mehran9675/silicondb/internal/format"
	"github.com/mehran9675/silicondb/internal/store"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/env"
	"github.com/mehran9675/silicondb/internal/ui/theme"
)

// Screen is the Settings tab.
type Screen struct {
	env     *env.Env
	scroll  *scroll.Scroll
	dark    *toggle.Switch
	source  *input.Input
	refresh *giode.ButtonElement
	redo    *giode.ButtonElement
	save    *giode.ButtonElement
	reset   *giode.ButtonElement
}

// New returns the Settings screen.
func New(e *env.Env) *Screen {
	s := &Screen{env: e}
	s.scroll = giode.Scroll(s.content, giode.Styles{ScrollBar: giode.ScrollBar{HiddenMaxWidth: 600}})
	s.dark = giode.Input.Switch("Dark mode", giode.Styles{Width: -1, FontFamily: theme.Family, FontSize: 15}).OnChange(func(on bool) {
		e.Store.SetDark(on)
		theme.SetDark(on)
	})
	s.dark.SetChecked(e.Store.Dark())
	s.source = giode.Input.Text("https://…/data/", giode.Styles{
		FontFamily: theme.Family, FontSize: 13.5, Color: theme.Text, Background: theme.Raised(),
		BorderColor: theme.Border, BorderRadius: 12,
	}).Keyboard(giode.KeyboardURL)
	s.source.SetText(e.Store.DataURL())
	s.refresh = giode.Button("Check for updates", giode.Styles{BorderRadius: 12}).
		Icon(theme.Icon("refresh", 16, "")).OnClick(e.Store.Refresh)
	s.redo = giode.Button("Download again", giode.Styles{BorderRadius: 12, Background: theme.Primary, Color: theme.Primary}).
		Variant(giode.VariantOutlined).OnClick(e.Store.ClearCache)
	s.save = giode.Button("Use this source", giode.Styles{BorderRadius: 12}).OnClick(func() {
		e.Store.SetDataURL(s.source.Text())
		e.Toasts.Show(giode.ToastInfo, "Downloading from the new source…")
	})
	s.reset = giode.Button("Default", giode.Styles{BorderRadius: 12, Background: theme.Primary, Color: theme.Primary}).
		Variant(giode.VariantGhost).OnClick(func() {
		s.source.SetText(config.DataURL)
		e.Store.SetDataURL("")
	})
	return s
}

// View lays the screen out.
func (s *Screen) View() giode.Element { return s.scroll }

func (s *Screen) content() giode.Element {
	cat := s.env.Store.Catalog.Get()
	st := s.env.Store.Status.Get()
	syncing := st.Phase == store.Syncing
	s.refresh.Loading(syncing)
	s.redo.Disabled(syncing)

	return components.Column(16,
		components.AppBar(nil, "Settings", ""),
		components.Card(giode.Styles{Padding: giode.SymmetricInset(16, 8)},
			giode.Row(giode.Styles{Width: -1, Gap: 12, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				components.IconBadge("moon", theme.Primary, 36),
				giode.Box(giode.Styles{FlexGrow: 1}, s.dark),
			)),
		components.SpecSection("Database", []components.SpecLine{
			{Label: "Chips", Value: format.Number(float64(len(cat.Chips)), 0)},
			{Label: "Manufacturers", Value: fmt.Sprint(len(cat.Manufacturers))},
			{Label: "Data updated", Value: format.Date(cat.Updated)},
			{Label: "Last checked", Value: ago(st.LastSynced)},
			{Label: "Status", Value: status(st)},
		}),
		giode.Row(giode.Styles{Gap: 10}, s.refresh, s.redo),
		components.Card(giode.Styles{Gap: 12},
			giode.Text("Data source", theme.Overline().With(giode.Styles{Color: theme.Primary})),
			giode.Text("Silicon DB downloads its data from this folder and keeps an offline copy. Only files that changed are downloaded on each update.",
				theme.Caption().With(giode.Styles{FontSize: 13, LineHeight: 19})),
			s.source,
			giode.Row(giode.Styles{Gap: 10}, s.save, s.reset),
		),
		components.SpecSection("About", []components.SpecLine{
			{Label: "Version", Value: config.Version},
			{Label: "Typeface", Value: "Inter (SIL OFL)"},
			{Label: "Built with", Value: "giode"},
		}),
	)
}

func status(st store.SyncStatus) string {
	switch st.Phase {
	case store.Syncing:
		if st.Total > 0 {
			return fmt.Sprintf("Downloading %d of %d", st.Done, st.Total)
		}
		return "Checking…"
	case store.Synced:
		if st.Files == 0 {
			return "Up to date"
		}
		return fmt.Sprintf("Updated %d files (%s)", st.Files, size(st.Bytes))
	case store.Failed:
		return "Offline"
	}
	return "—"
}

func size(b int) string {
	if b >= 1<<20 {
		return format.Number(float64(b)/(1<<20), 1) + " MB"
	}
	return format.Number(float64(b)/(1<<10), 0) + " KB"
}

func ago(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "Just now"
	case d < time.Hour:
		return fmt.Sprintf("%d min ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d h ago", int(d.Hours()))
	}
	return t.Format("Jan 2, 2006")
}
