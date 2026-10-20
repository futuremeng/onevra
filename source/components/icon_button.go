package components

import "github.com/egoist/mygo/ui"

// IconButtonStyle selects an icon button's face.
type IconButtonStyle int

const (
	// IconButtonPlain tints a transparent face on hover and press.
	IconButtonPlain IconButtonStyle = iota
	// IconButtonToolbar raises the face, shadowed, as the chosen segment of a
	// segmented tab list is raised.
	IconButtonToolbar
	// IconButtonCompact is IconButtonPlain at a smaller square, for dense
	// rows such as a SidebarSectionWithActions heading.
	IconButtonCompact
	// IconButtonToolbarExpanding is IconButtonToolbar that widens to show its
	// label after the icon while hovered or keyboard-focused.
	IconButtonToolbarExpanding
)

// IconButton creates an icon-only button with an accessible label and a
// tooltip. It uses one theme-scaled square size. Applications own the label and
// handle Clicked on the result; MyGo owns its pointer, keyboard and focus
// behavior.
//
// IconButtonPlain keeps a transparent face and tints it on hover and press.
// IconButtonToolbar keeps the raised, shadowed face of ButtonToolbar, matching
// TabsSegmented's chosen segment. IconButtonCompact shrinks the plain face.
// IconButtonToolbarExpanding rests as the IconButtonToolbar square and, while
// hovered or keyboard-focused, springs wider to reveal the label beside the
// icon instead of showing a tooltip; neighbors in the same row move aside.
func IconButton(c *ui.Context, icon *ui.SVG, label string, style IconButtonStyle) ui.Element {
	t := c.Theme()
	size, glyph := max(32, t.Rem(2)), t.Rem(1.25)
	if style == IconButtonCompact {
		size, glyph = max(20, t.Rem(1.25)), t.Rem(.875)
	}
	var b ui.Element
	switch style {
	case IconButtonToolbar:
		b = Button(c, "", ButtonToolbar).Label(label).Size(size, size).Padding(0)
		b.Children(func() { Icon(c, icon).Size(t.Rem(1.125), t.Rem(1.125)) })
	case IconButtonToolbarExpanding:
		glyph = t.Rem(1.125)
		// The inset centers the icon in the resting square, within the border.
		b = Button(c, "", ButtonToolbar).Label(label).Height(size).
			Padding(0, (size-glyph)/2-1).Justify(ui.Start).Clip()
		// The layout follows the animated width rather than a Transition's
		// painted size, so an end-aligned row keeps the button's trailing edge
		// under the pointer and its neighbors move with it frame by frame.
		gap, trail := t.Space(1.5), t.Space(1)
		reveal := float32(0)
		if b.Hovered() || b.FocusVisible() {
			w, _ := c.MeasureText(0, ui.Span{Text: label})
			reveal = gap + w + trail
		}
		reveal = max(0, b.AnimateTo("reveal", reveal, ui.Snappy))
		b.Width(size + reveal)
		b.Children(func() {
			Icon(c, icon).Size(glyph, glyph).Shrink(0)
			if reveal > 0 {
				// The label keeps its natural width while the face clips it.
				ui.Row(c).PaddingStart(gap).PaddingEnd(trail).Shrink(0).
					Children(func() { Text(c, label).SingleLine().Unselectable().Shrink(0) })
			}
		})
		return b
	default:
		b = ui.ButtonBase(c).Size(size, size).Radius(t.Radius).TextColor(t.TextMuted).Label(label)
		// Resolve the face in the build pass's paint, after the caller applies
		// properties such as Disabled; later paint-only passes, such as an
		// exiting Dialog's copy, reuse it rather than the expired handle.
		var face ui.Color
		resolved := false
		b.Draw(func(p *ui.Painter, r ui.Rect) {
			if !resolved {
				resolved = true
				switch {
				case b.IsDisabled():
				case b.Pressed():
					face = t.Text.Alpha(0.14)
				case b.Hovered():
					face = t.Text.Alpha(0.08)
				}
			}
			if face.A > 0 {
				p.Fill(r, face, t.Radius)
			}
		}).Children(func() {
			Icon(c, icon).Size(glyph, glyph)
		})
	}
	Tooltip(c, b, label)
	return b
}
