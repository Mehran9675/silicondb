package settings

import (
	"fmt"
	"time"

	"github.com/mehran9675/silicondb/internal/format"
	"github.com/mehran9675/silicondb/internal/store"
	"github.com/mehran9675/silicondb/internal/ui/tr"
)

func status(st store.SyncStatus) string {
	switch st.Phase {
	case store.Syncing:
		if st.Total > 0 {
			return tr.T("status.downloading", "done", fmt.Sprint(st.Done), "total", fmt.Sprint(st.Total))
		}
		return tr.T("status.checking")
	case store.Synced:
		if st.Files == 0 {
			return tr.T("status.up_to_date")
		}
		return tr.T("status.updated", "count", fmt.Sprint(st.Files), "n", st.Files, "size", size(st.Bytes))
	case store.Failed:
		return tr.T("status.offline")
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
		return tr.T("status.never")
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return tr.T("status.just_now")
	case d < time.Hour:
		return tr.Count("status.minutes_ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return tr.Count("status.hours_ago", int(d.Hours()))
	}
	return format.Date(t.Format("2006-01-02"))
}
