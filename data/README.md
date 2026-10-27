# Silicon DB data

The app downloads everything it shows from this folder at runtime. Nothing about chips is
compiled into the app, so editing these files and pushing them updates every installed copy
on its next sync.

```
data/
  index.json                         generated: lists every shard with its revision
  schema.json                        categories, spec groups, labels, units, highlights
  manufacturers/<id>/
    manufacturer.json                the company
    <category>.json                  one shard per category the company has chips in
```

The app fetches `index.json` first, then only the shards (and `schema.json`) whose revision
differs from its cached copy. A typo fix in `manufacturers/zilog/cpu.json` therefore costs
users one small download, not the whole database.

## Updating

1. Edit or add files under `manufacturers/`.
2. Run `go run ./cmd/indexgen` from the repository root. It validates every file and rewrites
   `index.json` with fresh revisions and counts.
3. Commit and push.

## manufacturer.json

```json
{
  "id": "intel",
  "name": "Intel",
  "country": "United States",
  "founded": 1968,
  "defunct": 0,
  "website": "https://www.intel.com"
}
```

`id` matches the folder name: lowercase ASCII letters, digits and hyphens. `defunct` is the
year the company ceased to exist or was absorbed, 0 if it still exists.

## Shards: manufacturers/&lt;id&gt;/&lt;category&gt;.json

`<category>` is one of the category ids in `schema.json`: `cpu`, `gpu`, `apu`, `soc`, `ai`
(TPUs, NPUs and other AI accelerators), `mcu`, `dsp`, `fpga`.

```json
{
  "manufacturer": "intel",
  "category": "cpu",
  "chips": [
    {
      "id": "intel-8008",
      "name": "8008",
      "family": "MCS-8",
      "released": "1972-04",
      "specs": {
        "process_nm": 10000,
        "transistors": 3500,
        "data_width_bits": 8,
        "address_width_bits": 14,
        "base_clock_mhz": 0.5,
        "cores": 1,
        "package": "18-pin DIP"
      }
    }
  ]
}
```

| Field | Required | Notes |
| --- | --- | --- |
| `id` | yes | Unique across the whole database: `<manufacturer>-<model>`, lowercase, hyphens. |
| `name` | yes | Model name without the manufacturer (`"Core i9-14900K"`, `"GeForce RTX 4090"`). |
| `family` | yes | Product line used for grouping (`"Core i9"`, `"GeForce 40"`, `"Snapdragon 8"`). |
| `released` | yes | `YYYY`, `YYYY-MM` or `YYYY-MM-DD`. |
| `specs` | yes | Flat object of spec keys to values; see below. |

### Spec values

Keys come from `schema.json` (`groups[].fields[].key`). Values are numbers, strings, or
arrays of strings. Use the unit the key names:

- `*_mhz` in MHz (a 3.2 GHz clock is `3200`), `*_kb` in KiB, `*_mb` in MiB (a 24 GB card is
  `24576`), `process_nm` in nanometres (10 µm is `10000`), `tdp_w` in watts,
  `transistors` as a plain integer, `die_size_mm2` in mm², `*_tflops` in TFLOPS (dense, no
  sparsity), `*_tops` in TOPS, `launch_price_usd` in US dollars.
- `l1_cache_kb`, `l2_cache_kb`, `l3_cache_kb` are totals for the chip.
- `instruction_sets`, `graphics_api` and `connectivity` are arrays of strings.
- Leave a spec out rather than guess. Omitted specs simply aren't shown.

Keys not in `schema.json` are still shown, under "Other", with a label derived from the key;
adding them to `schema.json` later gives them a proper label, unit and group without an app
update.
