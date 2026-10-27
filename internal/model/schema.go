package model

// Schema describes how to present the data: the categories, and the
// spec fields grouped into sections with labels and units. It lives in
// the remote database so new fields need no app update.
type Schema struct {
	Format     int                 `json:"format"`
	Categories []Category          `json:"categories"`
	Groups     []Group             `json:"groups"`
	Highlights map[string][]string `json:"highlights"`

	fields map[string]*Field
}

// Category is a kind of chip: CPU, GPU, SoC and so on.
type Category struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// Group is a section of related spec fields.
type Group struct {
	ID     string  `json:"id"`
	Label  string  `json:"label"`
	Fields []Field `json:"fields"`
}

// Field describes one spec key.
type Field struct {
	Key    string `json:"key"`
	Label  string `json:"label"`
	Unit   string `json:"unit"`
	Format string `json:"format"`
}

// Field returns the description of key, or nil if the schema doesn't
// know it.
func (s *Schema) Field(key string) *Field {
	if s.fields == nil {
		s.fields = map[string]*Field{}
		for gi := range s.Groups {
			for fi := range s.Groups[gi].Fields {
				f := &s.Groups[gi].Fields[fi]
				s.fields[f.Key] = f
			}
		}
	}
	return s.fields[key]
}

// Category returns the category with id, or a placeholder for an id the
// schema doesn't list.
func (s *Schema) Category(id string) Category {
	for _, c := range s.Categories {
		if c.ID == id {
			return c
		}
	}
	return Category{ID: id, Name: id, Title: id, Icon: "cpu", Color: "#8B90A0"}
}

// fallbackHighlights fill in when a chip lacks its category's highlights.
var fallbackHighlights = []string{
	"cores", "boost_clock_mhz", "base_clock_mhz", "process_nm", "transistors",
	"tdp_w", "data_width_bits", "memory_size_mb", "fp32_tflops", "int8_tops",
}

// HighlightsFor returns up to n spec keys worth showing first for chip:
// its category's highlights it has, then other common specs it has.
func (s *Schema) HighlightsFor(chip *Chip, n int) []string {
	var keys []string
	add := func(k string) {
		if len(keys) >= n {
			return
		}
		if _, ok := chip.Specs[k]; !ok {
			return
		}
		for _, have := range keys {
			if have == k {
				return
			}
		}
		keys = append(keys, k)
	}
	for _, k := range s.Highlights[chip.Category] {
		add(k)
	}
	for _, k := range fallbackHighlights {
		add(k)
	}
	return keys
}
