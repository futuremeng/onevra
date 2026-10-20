package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/flow"
)

var _ = registerWindow(&projectPlan{})

type projectPlan struct {
	state    flow.State // layout and camera
	selected string
}

func (p *projectPlan) View(c *ui.Context) {
	card := func(title, detail string) func(*ui.Context) {
		return func(c *ui.Context) {
			flow.Card(c, flow.CardOptions{Title: title, Description: detail, Width: 200})
		}
	}
	data := flow.Data{
		Nodes: []flow.Node{
			{ID: "survey", Label: "Survey", Content: card("Survey", "Collect requirements")},
			{ID: "design", Label: "Design", Content: card("Design", "Sketch the approach")},
			{ID: "prototype", Label: "Prototype", Content: card("Prototype", "Try the riskiest part")},
			{ID: "review", Label: "Review", Content: card("Review", "Decide what ships")},
			{ID: "ship", Label: "Ship", Content: card("Ship", "Release and announce")},
		},
		Edges: []flow.Edge{
			{Source: "survey", Target: "design"},
			{Source: "survey", Target: "prototype", Dashed: true, Label: "informs"},
			{Source: "design", Target: "review"},
			{Source: "prototype", Target: "review"},
			{Source: "review", Target: "ship"},
		},
	}
	flow.View(c, &p.state, data, flow.Options{
		Label:    "Project plan",
		Selected: p.selected,
		Select:   func(id string) { p.selected = id },
	})
}
