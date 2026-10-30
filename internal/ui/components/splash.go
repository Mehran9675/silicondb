package components

import (
	"github.com/mehran9675/giode"

	"github.com/mehran9675/silicondb/internal/store"
	"github.com/mehran9675/silicondb/internal/ui/theme"
	"github.com/mehran9675/silicondb/internal/ui/tr"
)

// Splash fills the screen while the first download runs, before there is
// anything to show.
func Splash(st store.SyncStatus) giode.Element {
	label := tr.T("splash.connecting")
	var progress giode.Element = giode.Spinner(giode.Styles{Width: 26, Height: 26, Color: theme.Primary})
	if st.Total > 0 {
		label = tr.T("splash.downloading")
		progress = giode.Box(giode.Styles{Width: 220},
			giode.Progress(float32(st.Done)/float32(st.Total), giode.Styles{Height: 6, Color: theme.Primary, Background: theme.Faint()}))
	}
	return giode.Stack(giode.Styles{Width: -1, Height: -1, Align: giode.AlignCenter, Justify: giode.JustifyCenter, Gap: 18},
		giode.Box(giode.Styles{
			Width: 84, Height: 84, Align: giode.AlignCenter, Justify: giode.JustifyCenter,
			Background: "#7C3AED", BackgroundGradient: theme.BrandGradient(135), BorderRadius: 26,
			BoxShadow: giode.BoxShadow{Y: 14, Blur: 36, Color: "#7C3AED55"},
		}, theme.Icon("cpu", 44, "#FFFFFF")),
		giode.Text(tr.T("app.name"), theme.Title()),
		giode.Text(label, theme.Caption().With(giode.Styles{FontSize: 14})),
		progress,
	)
}
