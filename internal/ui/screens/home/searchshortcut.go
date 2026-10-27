package home

import (
	"github.com/mehran9675/giode"

	"github.com/mehran9675/silicondb/internal/ui/theme"
)

// searchShortcut looks like a search field and opens Explore.
func searchShortcut(onClick func()) *giode.PressableElement {
	return giode.Pressable(func(s giode.PressState) giode.Element {
		bg := theme.Surface
		if s.Pressed {
			bg = theme.Pressed()
		}
		return giode.Row(giode.Styles{
			Width: -1, Align: giode.AlignCenter, Gap: 10, Wrap: giode.WrapNoWrap,
			Padding:      giode.SymmetricInset(16, 15),
			Background:   bg,
			BorderRadius: 16,
			BorderWidth:  1,
			BorderColor:  theme.Border,
		},
			theme.Icon("search", 18, theme.Muted),
			giode.Text("Search 8008, Zen 5, H100, M4…", theme.Body().With(giode.Styles{Color: theme.Muted})),
		)
	}).Label("Search").OnClick(func(giode.ClickEvent) { onClick() })
}
