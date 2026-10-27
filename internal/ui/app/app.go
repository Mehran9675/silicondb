// Package app assembles the screens into the running application.
package app

import (
	"image"
	"time"

	"path/filepath"

	gioapp "gioui.org/app"
	"github.com/mehran9675/giode"

	"github.com/mehran9675/silicondb/assets"
	"github.com/mehran9675/silicondb/internal/store"
	"github.com/mehran9675/silicondb/internal/ui/components"
	"github.com/mehran9675/silicondb/internal/ui/env"
	"github.com/mehran9675/silicondb/internal/ui/nav"
	"github.com/mehran9675/silicondb/internal/ui/screens/compare"
	"github.com/mehran9675/silicondb/internal/ui/screens/detail"
	"github.com/mehran9675/silicondb/internal/ui/screens/explore"
	"github.com/mehran9675/silicondb/internal/ui/screens/home"
	"github.com/mehran9675/silicondb/internal/ui/screens/maker"
	"github.com/mehran9675/silicondb/internal/ui/screens/makers"
	"github.com/mehran9675/silicondb/internal/ui/screens/settings"
	"github.com/mehran9675/silicondb/internal/ui/theme"
)

// Tab indexes.
const (
	tabHome = iota
	tabExplore
	tabMakers
	tabCompare
	tabSettings
)

// App is Silicon DB.
type App struct {
	env     *env.Env
	nav     *nav.Navigator
	explore *explore.Screen
	bar     *components.BottomNav
	retry   *giode.ButtonElement
	window  *giode.App
}

// New prepares the app: fonts, icons, theme, state and screens.
func New() *App {
	theme.LoadFonts()
	theme.LoadIcons()

	dir := ""
	if d, err := gioapp.DataDir(); err == nil {
		dir = filepath.Join(d, "silicondb")
	}
	st := store.New(dir)
	theme.Apply(st.Dark())

	a := &App{}
	a.env = &env.Env{
		Store:  st,
		Routes: a,
		Toasts: giode.Toaster().Placement(giode.ToastBottomCenter).Duration(2500 * time.Millisecond).Max(2),
	}
	a.explore = explore.New(a.env)
	a.nav = nav.New(home.New(a.env), a.explore, makers.New(a.env), compare.New(a.env), settings.New(a.env))
	a.bar = components.NewBottomNav([]components.NavTab{
		{Icon: "home", Label: "Home"},
		{Icon: "explore", Label: "Explore"},
		{Icon: "factory", Label: "Makers"},
		{Icon: "compare", Label: "Compare"},
		{Icon: "settings", Label: "Settings"},
	}, a.nav.Tab, a.nav.SetTab)
	a.retry = giode.Button("Try again", giode.Styles{BorderRadius: 14, Padding: giode.SymmetricInset(20, 12)}).OnClick(st.Refresh)

	a.window = giode.New(giode.Config{
		Title:        "Silicon DB",
		Size:         image.Pt(420, 880),
		Resizable:    true,
		MinWidth:     340,
		MinHeight:    560,
		Background:   theme.Background,
		Icon:         assets.AppIcon,
		OnBackground: st.Save,
		OnStartup:    st.Start,
	})
	return a
}

// Run shows the window until it closes.
func (a *App) Run() error { return a.window.Run(a.view) }

func (a *App) view() giode.Element {
	page := func() giode.Element { return a.page() }
	return giode.Stack(giode.Styles{Width: -1, Height: -1},
		giode.Responsive(
			giode.Breakpoint{MaxWidth: 600, Name: "compact", View: page},
			giode.Breakpoint{Name: "wide", View: page},
		),
		a.env.Toasts,
	)
}

// page is the current screen above the tab bar, or the first-run states.
func (a *App) page() giode.Element {
	st := a.env.Store.Status.Get()
	if a.env.Store.Catalog.Get() == nil {
		if st.Phase == store.Failed {
			return components.EmptyState("cloud", "Can't reach the database",
				"Silicon DB needs a connection the first time it starts to download the chip data. "+st.Err, a.retry)
		}
		return components.Splash(st)
	}
	return giode.Stack(giode.Styles{Width: -1, Height: -1},
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1}, a.nav.View()),
		a.bar.View(),
	)
}

// OpenChip implements env.Routes.
func (a *App) OpenChip(id string) { a.nav.Push(detail.New(a.env, id)) }

// OpenMaker implements env.Routes.
func (a *App) OpenMaker(id string) { a.nav.Push(maker.New(a.env, id)) }

// Explore implements env.Routes.
func (a *App) Explore(q store.Query) {
	a.explore.Apply(q)
	a.nav.SetTab(tabExplore)
}

// Makers implements env.Routes.
func (a *App) Makers() { a.nav.SetTab(tabMakers) }

// Back implements env.Routes.
func (a *App) Back() { a.nav.Back() }
