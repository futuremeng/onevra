package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/graph"
)

var _ = registerWindow(&notesGraph{})

type notesGraph struct {
	state  graph.State // layout and camera, kept across frames
	opened string
}

var notesData = graph.Data{
	Nodes: []graph.Node{
		{ID: "index", Label: "Index", Openable: true, Size: 12},
		{ID: "ideas", Label: "Ideas", Group: "notes", Openable: true},
		{ID: "journal", Label: "Journal", Group: "notes", Openable: true},
		{ID: "reading", Label: "Reading list", Group: "notes"},
		{ID: "go", Label: "Go", Group: "topics"},
		{ID: "design", Label: "Design", Group: "topics"},
		{ID: "ui", Label: "Interfaces", Group: "topics"},
		{ID: "sketch", Label: "sketch.png", Group: "files", Size: 6},
		{ID: "talk", Label: "talk.pdf", Group: "files", Size: 6},
	},
	Edges: []graph.Edge{
		{Source: "index", Target: "ideas"}, {Source: "index", Target: "journal"},
		{Source: "index", Target: "reading"}, {Source: "ideas", Target: "go"},
		{Source: "ideas", Target: "design"}, {Source: "design", Target: "ui"},
		{Source: "go", Target: "ui"}, {Source: "ideas", Target: "sketch"},
		{Source: "reading", Target: "talk"}, {Source: "journal", Target: "design"},
	},
}

func (n *notesGraph) View(c *ui.Context) {
	// Drag nodes, wheel or pinch to zoom, and double-click an openable node.
	graph.View(c, &n.state, notesData, graph.Options{
		Label: "Notes",
		Open:  func(id string) { n.opened = id },
	})
}
