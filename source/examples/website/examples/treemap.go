package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/treemap"
)

var _ = registerWindow(&diskUsage{})

type diskUsage struct {
	state treemap.State // the drilled branch and active tile
}

var diskNodes = []treemap.Node{
	{ID: "photos", Label: "Photos", Children: []treemap.Node{
		{ID: "photos/2025", Label: "2025", Value: 962},
		{ID: "photos/2026", Label: "2026", Value: 382},
		{ID: "photos/raw", Label: "RAW", Value: 1480},
	}},
	{ID: "music", Label: "Music", Value: 640},
	{ID: "video", Label: "Video", Children: []treemap.Node{
		{ID: "video/films", Label: "Films", Value: 2210},
		{ID: "video/clips", Label: "Clips", Value: 410},
	}},
	{ID: "apps", Label: "Applications", Value: 880},
	{ID: "docs", Label: "Documents", Value: 120},
}

func (d *diskUsage) View(c *ui.Context) {
	ui.Column(c).Fill().Padding(c.Theme().Space(4)).Children(func() {
		// Click a tile to drill into it; the breadcrumbs lead back.
		treemap.View(c, &d.state, diskNodes, treemap.Options{
			Label:       "Disk usage",
			FormatValue: func(v float64) string { return fmt.Sprintf("%.1f GB", v/1000) },
		})
	})
}
