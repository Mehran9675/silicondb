package model

// Manufacturer is a company that designed or made chips.
type Manufacturer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Founded int    `json:"founded"`
	Defunct int    `json:"defunct"`
	Website string `json:"website"`
}

// Active reports whether the company still exists.
func (m *Manufacturer) Active() bool { return m.Defunct == 0 }
