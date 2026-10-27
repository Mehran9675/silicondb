# Silicon DB

A mobile app with the specifications of processing chips — CPUs, GPUs, APUs, SoCs,
TPUs/NPUs and other AI accelerators, microcontrollers, DSPs and FPGAs — from every
manufacturer, from the Intel 8008 (1972) to today. Built in Go with
[giode](https://github.com/mehran9675/giode) and packaged with the
[giode CLI](https://github.com/mehran9675/giode-cli).

## Features

- **Home**: the database at a glance, categories, latest releases, the earliest chips,
  top manufacturers and recently viewed chips.
- **Explore**: search by name, family, codename or architecture; filter by category,
  decade and manufacturer; sort by date or name.
- **Chip pages**: headline figures and every spec, grouped (fabrication, cores and clocks,
  cache, graphics, memory, performance, power and I/O).
- **Manufacturers**: every company with its history and its chips by category.
- **Compare**: up to four chips side by side, with the best value in each row highlighted.
- **Offline first**: the last download is kept on the device, and updates fetch only the
  files that changed.
- Dark and light themes, phone and tablet layouts.

## Data

Nothing about chips is compiled into the app. It downloads the database from
[`data/`](data/README.md) at runtime:

```
data/index.json                          table of contents with a revision per file
data/schema.json                         categories, spec groups, labels and units
data/manufacturers/<id>/manufacturer.json
data/manufacturers/<id>/<category>.json  one shard per manufacturer and category
```

The app compares each file's revision in `index.json` with its offline copy and downloads
only what changed. Spec labels, units and groups live in `schema.json`, so new fields
need no app update either. See [data/README.md](data/README.md) for the format and how to
edit it; after editing, run:

```sh
go run ./cmd/indexgen        # validates everything and rewrites data/index.json
```

By default the app reads `https://raw.githubusercontent.com/Mehran9675/silicondb/main/data/`.
Change it at build time with
`-ldflags "-X github.com/mehran9675/silicondb/internal/config.DataURL=<url>"` (or
`app.ldflags` in `giode.toml`), or at runtime in Settings → Data source.

## Building

Requires Go 1.26+ and the giode CLI:

```sh
go install github.com/mehran9675/giode-cli/cmd/giode@latest
giode doctor                 # what each target needs
giode run android            # build, install and launch on a connected phone
giode build android          # dist/android/Silicon DB.apk
giode build ios              # on macOS, with a provisioning profile in giode.toml
giode build web              # dist/web
giode run                    # on this computer
```

## Project layout

```
main.go                     entry point
assets/                     embedded fonts (Inter), SVG icons and the app icon
cmd/indexgen/               data validator and index generator
data/                       the database (downloaded by the app, not embedded)
internal/
  config/                   build-time settings (data URL, version)
  model/                    chips, manufacturers, shards, index, schema, catalog
  data/                     HTTP client, offline cache, incremental sync
  format/                   spec values, units and dates for display
  store/                    app state: catalog, sync status, queries, preferences
  ui/
    theme/                  palettes, type scale, fonts, icons
    nav/                    tab and stack navigation
    env/                    what screens share
    components/             reusable widgets, one per file
    screens/                home, explore, detail, makers, maker, compare, settings
    app/                    assembles the screens into the app
```

## Tests

```sh
go test ./internal/...
```

## License

Apache License 2.0; see [LICENSE](LICENSE). Inter is under the SIL Open Font License
([assets/fonts/OFL.txt](assets/fonts/OFL.txt)).
