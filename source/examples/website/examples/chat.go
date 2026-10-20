package examples

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chat"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = registerWindow(&assistant{turns: []*assistantTurn{{
	id: "hello", prompt: "What is Onevra?", status: chat.Completed,
	reply: "A Go component library for **native desktop apps**, built on MyGo.",
}}})

// assistantReply is what the simulated model answers; a real application
// would receive it from a network request, off the UI thread.
const assistantReply = "Every Onevra control is drawn by **MyGo**, so a window needs:\n\n" +
	"1. A *builder* that runs every frame.\n2. State kept in your own values.\n3. No web view at all.\n\n" +
	"```go\ncomponents.Button(c, \"Save\", components.ButtonPrimary)\n```"

type assistantTurn struct {
	id, prompt, reply string
	status            chat.Status
	started           time.Time
	state             chat.TurnState
	search            chat.ActivityState
	group             chat.GroupState
}

type assistant struct {
	layout       chat.FloatingConversationState
	conversation chat.ConversationState
	turns        []*assistantTurn
	draft        string
	markdown     markdown.Options
}

func (a *assistant) View(c *ui.Context) {
	now := c.Now()
	running := false
	for _, turn := range a.turns {
		if turn.status == chat.Running {
			running = true
			// Reveal the reply over time, as a stream would deliver it.
			if n := int(now.Sub(turn.started)/(30*time.Millisecond)) * 4; n < len(assistantReply) {
				for n < len(assistantReply) && !utf8.RuneStart(assistantReply[n]) {
					n++
				}
				turn.reply = assistantReply[:n]
				c.After(30 * time.Millisecond)
			} else {
				turn.reply, turn.status = assistantReply, chat.Completed
			}
		}
	}
	chat.FloatingConversation(c, &a.layout, &a.conversation, "Latest", func(c *ui.Context) {
		for _, turn := range a.turns {
			a.turn(c, turn)
		}
	}, func(c *ui.Context) {
		result := chat.GlassComposer(c, &a.layout, &a.draft, chat.ComposerOptions{
			Label: "Message", Placeholder: "Ask anything", SendLabel: "Send", StopLabel: "Stop",
			Running: running,
		})
		if result.Send && strings.TrimSpace(a.draft) != "" {
			a.turns = append(a.turns, &assistantTurn{
				id: fmt.Sprint("turn-", len(a.turns)), prompt: strings.TrimSpace(a.draft),
				status: chat.Running, started: now,
			})
			a.draft = ""
			a.conversation.JumpToLatest()
		}
		if result.Stop {
			a.turns[len(a.turns)-1].status = chat.Cancelled
		}
	})
}

func (a *assistant) turn(c *ui.Context, turn *assistantTurn) {
	var steps []chat.Step
	var final chat.View
	if turn.started != (time.Time{}) {
		search := chat.Completed
		if turn.status == chat.Running && turn.reply == "" {
			search = chat.Running
		}
		steps = append(steps, chat.Step{ID: "work", Title: "Searched the docs", Group: &turn.group,
			Activities: []chat.Activity{{ID: "search", Icon: lucide.Search, Action: "Search", Detail: "docs",
				Status: search, State: &turn.search,
				Content: func(c *ui.Context) { components.Text(c, "Found 3 pages about builders.") }}}})
	}
	if turn.status == chat.Running {
		if turn.reply != "" {
			steps = append(steps, chat.Step{ID: "reply", Message: func(c *ui.Context) {
				markdown.Stream(c, turn.reply, &a.markdown).FillWidth()
			}})
		}
	} else if turn.reply != "" {
		final = func(c *ui.Context) { markdown.View(c, turn.reply, &a.markdown).FillWidth() }
	}
	chat.Turn(c, &turn.state, chat.TurnOptions{
		ID: turn.id, Status: turn.status, Summary: "Worked for 2s", Working: "Working",
		Prompt: func(c *ui.Context) { components.Text(c, turn.prompt) },
		Steps:  steps, Final: final, HoldOpen: !a.conversation.AtLatest(),
	})
}
