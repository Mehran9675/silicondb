package components

import (
	"fmt"

	"github.com/mehran9675/giode"

	"github.com/mehran9675/silicondb/internal/model"
	"github.com/mehran9675/silicondb/internal/ui/theme"
)

// CategoryTile is a tappable tile for a category with its chip count.
func CategoryTile(c model.Category, count int, onClick func()) giode.Element {
	return giode.Keyed("categorytile", c.ID, func() *giode.PressableElement {
		return giode.Pressable(func(s giode.PressState) giode.Element {
			bg := theme.Surface
			if s.Pressed {
				bg = theme.Pressed()
			} else if s.Hovered {
				bg = theme.Raised()
			}
			return giode.Stack(giode.Styles{
				Width: -1, Gap: 12,
				Padding:      giode.UniformInset(16),
				Background:   bg,
				BorderRadius: theme.RadiusCard,
				BorderWidth:  1,
				BorderColor:  theme.Border,
			},
				giode.Row(giode.Styles{Width: -1, Justify: giode.JustifySpaceBetween, Align: giode.AlignStart, Wrap: giode.WrapNoWrap},
					IconBadge(c.Icon, c.Color, 42),
					giode.Text(fmt.Sprint(count), giode.Styles{FontFamily: theme.Family, FontSize: 13, FontWeight: giode.FontWeightSemiBold, Color: theme.Muted}),
				),
				giode.Stack(giode.Styles{Gap: 2},
					giode.Text(c.Name, theme.Strong().With(giode.Styles{FontSize: 16})),
					giode.Text(c.Title, theme.Caption().With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
				),
			)
		}).Label(c.Name).OnClick(func(giode.ClickEvent) { onClick() })
	})
}
