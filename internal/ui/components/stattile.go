package components

import (
	"github.com/mehran9675/giode"

	"github.com/mehran9675/silicondb/internal/ui/theme"
)

// StatTile shows one headline figure with its label, tinted by color.
func StatTile(label, value, color string) giode.Element {
	return giode.Stack(giode.Styles{
		Width: -1, Gap: 6,
		Padding:      giode.UniformInset(14),
		Background:   theme.Tint(color, 0.08),
		BorderRadius: theme.RadiusTile,
		BorderWidth:  1,
		BorderColor:  theme.Tint(color, 0.22),
	},
		giode.Text(label, theme.Overline().With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
		giode.Text(value, theme.Value().With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
	)
}
