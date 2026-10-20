package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = registerWindow(&preferences{sounds: true, scale: 100})

type preferences struct {
	state           components.SettingsState
	sounds, banners bool
	launch          bool
	scale           float64
	name            string
}

func (p *preferences) View(c *ui.Context) {
	categories := []components.SettingsCategory{
		{ID: "general", Title: "General", Icon: lucide.Settings(), Groups: []components.SettingsGroup{
			{ID: "startup", Title: "Startup", Items: []components.SettingsItem{
				{ID: "launch", Title: "Open at login", Description: "Start when you sign in.",
					Build: func(c *ui.Context) {
						components.Switch(c, &p.launch, "Open at login", components.SwitchOptions{HideLabel: true})
					}},
			}},
			{ID: "profile", Title: "Profile", Items: []components.SettingsItem{
				{ID: "name", Title: "Display name", Stacked: true,
					Build: func(c *ui.Context) {
						components.TextBox(c, &p.name, components.TextBoxOptions{}).Label("Display name").FillWidth()
					}},
			}},
		}},
		{ID: "notifications", Title: "Notifications", Icon: lucide.Bell(), Groups: []components.SettingsGroup{
			{ID: "alerts", Items: []components.SettingsItem{
				{ID: "sounds", Title: "Play sounds", Keywords: []string{"audio"},
					Build: func(c *ui.Context) {
						components.Switch(c, &p.sounds, "Play sounds", components.SwitchOptions{HideLabel: true})
					}},
				{ID: "banners", Title: "Show banners",
					Build: func(c *ui.Context) {
						components.Switch(c, &p.banners, "Show banners", components.SwitchOptions{HideLabel: true})
					}},
			}},
		}},
		{ID: "appearance", Title: "Appearance", Icon: lucide.Palette(), Groups: []components.SettingsGroup{
			{ID: "text", Items: []components.SettingsItem{
				{ID: "scale", Title: "Text size", Keywords: []string{"font", "zoom"}, Stacked: true,
					Build: func(c *ui.Context) {
						components.Slider(c, &p.scale, 80, 140, components.SliderOptions{Step: 10}).Label("Text size").FillWidth()
					}},
			}},
		}},
	}
	// The host chooses the narrow layout from the space it gives Settings.
	width, _ := c.Size()
	components.Settings(c, &p.state, categories, components.SettingsOptions{Narrow: width < 560})
}
