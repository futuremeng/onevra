package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// support is the predicted state of a component on a mobile platform.
type support int

const (
	// supportReady works with touch as it is.
	supportReady support = iota
	// supportAdapt keeps the component with touch-specific changes.
	supportAdapt
	// supportNative defers to the platform's own control.
	supportNative
	// supportNew is a mobile component Onevra does not have yet.
	supportNew
	// supportDesktop has no phone equivalent; another pattern replaces it.
	supportDesktop
)

var supports = []support{supportReady, supportAdapt, supportNative, supportNew, supportDesktop}

// name and summary return the message keys of a level's badge and legend.
func (s support) name() string {
	return [...]string{"supportReady", "supportAdapt", "supportNative", "supportNew", "supportDesktop"}[s]
}

func (s support) summary() string {
	return [...]string{"summaryReady", "summaryAdapt", "summaryNative", "summaryNew", "summaryDesktop"}[s]
}

// badge returns the level's badge options: new components share the accent
// with native ones and add an icon.
func (s support) badge() components.BadgeOptions {
	o := components.BadgeOptions{Tone: [...]components.BadgeTone{components.BadgeSuccess, components.BadgeWarning,
		components.BadgeAccent, components.BadgeAccent, components.BadgeNeutral}[s]}
	if s == supportNew {
		o.Icon = lucide.Sparkles()
	}
	return o
}

// entry is one component's prediction. ID is the public API name, which is not
// translated; group and note are message keys.
type entry struct {
	id           string
	group        string
	ios, android support
	note         string
	// preview optionally shows the component live on the simulated screen.
	preview func(c *ui.Context, a *app)
}

func (e entry) on(p platform) support {
	if p == platformAndroid {
		return e.android
	}
	return e.ios
}

var groups = []string{"groupActions", "groupInput", "groupNavigation", "groupContent", "groupFeedback", "groupMobile"}

// catalog lists the predictions in group order. It covers the components a
// mobile application would reach for, not every Onevra API.
var catalog = []entry{
	{"Button", "groupActions", supportReady, supportReady, "noteButton", previewButton},
	{"IconButton", "groupActions", supportAdapt, supportAdapt, "noteIconButton", previewIconButton},
	{"FloatingButtonGroup", "groupActions", supportReady, supportReady, "noteFloating", previewFloating},
	{"HoldToConfirm", "groupActions", supportReady, supportReady, "noteHold", previewHold},

	{"TextBox", "groupInput", supportNative, supportNative, "noteTextBox", previewTextBox},
	{"Switch", "groupInput", supportReady, supportReady, "noteSwitch", previewSwitch},
	{"Checkbox", "groupInput", supportAdapt, supportReady, "noteCheckbox", previewCheckbox},
	{"Radio", "groupInput", supportAdapt, supportReady, "noteRadio", previewRadio},
	{"Slider", "groupInput", supportReady, supportReady, "noteSlider", previewSlider},
	{"NumberInput", "groupInput", supportAdapt, supportAdapt, "noteNumberInput", nil},
	{"Select", "groupInput", supportAdapt, supportAdapt, "noteSelect", nil},
	{"NativeSelect", "groupInput", supportNative, supportNative, "noteNativeSelect", nil},
	{"TagInput", "groupInput", supportAdapt, supportAdapt, "noteTagInput", nil},
	{"CalendarPicker", "groupInput", supportNative, supportNative, "noteDate", nil},
	{"TimePicker", "groupInput", supportNative, supportNative, "noteDate", nil},
	{"ColorPicker", "groupInput", supportNative, supportAdapt, "noteColor", nil},

	{"Tabs", "groupNavigation", supportAdapt, supportAdapt, "noteTabs", previewTabs},
	{"ScrollEdgeBar", "groupNavigation", supportReady, supportAdapt, "noteScrollEdge", nil},
	{"Sidebar", "groupNavigation", supportAdapt, supportDesktop, "noteSidebar", nil},
	{"Breadcrumbs", "groupNavigation", supportDesktop, supportDesktop, "noteBreadcrumbs", nil},
	{"Split", "groupNavigation", supportDesktop, supportDesktop, "noteSplit", nil},
	{"command.CommandView", "groupNavigation", supportAdapt, supportAdapt, "noteCommand", nil},

	{"List", "groupContent", supportReady, supportReady, "noteList", nil},
	{"Tree", "groupContent", supportAdapt, supportAdapt, "noteTree", nil},
	{"Table", "groupContent", supportAdapt, supportAdapt, "noteTable", nil},
	{"GridView", "groupContent", supportReady, supportReady, "noteGrid", nil},
	{"Carousel", "groupContent", supportReady, supportReady, "noteSwipeable", nil},
	{"CardStack", "groupContent", supportReady, supportReady, "noteSwipeable", nil},
	{"chart", "groupContent", supportAdapt, supportAdapt, "noteChart", nil},
	{"chat", "groupContent", supportReady, supportReady, "noteChat", nil},
	{"Avatar", "groupContent", supportReady, supportReady, "noteDisplay", previewAvatar},
	{"Badge", "groupContent", supportReady, supportReady, "noteDisplay", previewBadge},

	{"Dialog", "groupFeedback", supportAdapt, supportAdapt, "noteDialog", previewSheetButton},
	{"AlertDialog", "groupFeedback", supportNative, supportNative, "noteAlert", nil},
	{"Popover", "groupFeedback", supportAdapt, supportAdapt, "notePopover", nil},
	{"Tooltip", "groupFeedback", supportDesktop, supportDesktop, "noteTooltip", nil},
	{"Toasts", "groupFeedback", supportAdapt, supportNative, "noteToast", nil},
	{"Progress", "groupFeedback", supportReady, supportReady, "noteDisplay", previewProgress},
	{"Spinner", "groupFeedback", supportReady, supportReady, "noteDisplay", previewSpinner},

	{"Sheet", "groupMobile", supportNew, supportNew, "noteSheet", previewSheetButton},
	{"NavigationStack", "groupMobile", supportNew, supportNew, "noteNavigationStack", nil},
	{"TabBar", "groupMobile", supportNew, supportNew, "noteTabBar", nil},
	{"SafeArea", "groupMobile", supportNew, supportNew, "noteSafeArea", previewSafeArea},
	{"SwipeActions", "groupMobile", supportNew, supportNew, "noteSwipeActions", nil},
	{"PullToRefresh", "groupMobile", supportNew, supportNew, "notePullToRefresh", nil},
}

// find returns the catalog entry with id.
func find(id string) (entry, bool) {
	for _, e := range catalog {
		if e.id == id {
			return e, true
		}
	}
	return entry{}, false
}

// count returns how many entries have level s on p.
func count(p platform, s support) int {
	n := 0
	for _, e := range catalog {
		if e.on(p) == s {
			n++
		}
	}
	return n
}

// demo keeps the previews' values in memory, shared by every device.
type demo struct {
	on, checked, held bool
	choice            int
	tab               int
	value             float64
	text              string
}

func previewButton(c *ui.Context, a *app) {
	ui.Column(c).Gap(10).Children(func() {
		components.Button(c, a.text("sampleSave"), components.ButtonPrimary).Height(touchTarget(a))
		components.Button(c, a.text("sampleCancel"), components.ButtonSecondary).Height(touchTarget(a))
	})
}

func previewIconButton(c *ui.Context, a *app) {
	ui.Row(c).Gap(8).Children(func() {
		components.IconButton(c, lucide.Search(), a.text("sampleSearch"), components.IconButtonPlain)
		components.IconButton(c, lucide.Plus(), a.text("sampleAdd"), components.IconButtonToolbar)
	})
}

func previewFloating(c *ui.Context, a *app) {
	components.FloatingButtonGroup(c, a.text("sampleActions"), components.FloatingButtonGroupRaised, ui.Horizontal,
		[]components.FloatingButtonGroupItem{
			{ID: "search", Label: a.text("sampleSearch"), Icon: lucide.Search, IconOnly: true},
			{ID: "add", Label: a.text("sampleAdd"), Icon: lucide.Plus},
		})
}

func previewHold(c *ui.Context, a *app) {
	components.HoldToConfirm(c, &a.demo.held, a.text("sampleDelete"))
}

func previewTextBox(c *ui.Context, a *app) {
	// Focusing the field shows the simulated keyboard, as a tap would.
	components.TextBox(c, &a.demo.text, components.TextBoxOptions{}).Label(a.text("sampleName")).
		Placeholder(a.text("sampleName")).Bind(&a.fields[0])
}

func previewSwitch(c *ui.Context, a *app) {
	components.Switch(c, &a.demo.on, a.text("sampleWifi"))
}

func previewCheckbox(c *ui.Context, a *app) {
	components.Checkbox(c, &a.demo.checked, a.text("sampleRemember"))
}

func previewRadio(c *ui.Context, a *app) {
	ui.Column(c).Gap(10).Children(func() {
		for i, key := range []string{"sampleSmall", "sampleMedium", "sampleLarge"} {
			components.Radio(c, &a.demo.choice, i, a.text(key))
		}
	})
}

func previewSlider(c *ui.Context, a *app) {
	components.Slider(c, &a.demo.value, 0, 1)
}

func previewTabs(c *ui.Context, a *app) {
	components.Tabs(c, &a.demo.tab, components.TabsSegmented,
		a.text("sampleDay"), a.text("sampleWeek"), a.text("sampleMonth"))
}

func previewAvatar(c *ui.Context, a *app) {
	ui.Row(c).Gap(8).Children(func() {
		components.Avatar(c, "Ada Lovelace", components.AvatarOptions{})
		components.Avatar(c, "Grace Hopper", components.AvatarOptions{Status: components.AvatarOnline})
	})
}

func previewBadge(c *ui.Context, a *app) {
	ui.Row(c).Gap(8).Wrap().Children(func() {
		for _, s := range supports {
			components.Badge(c, a.text(s.name()), s.badge())
		}
	})
}

func previewProgress(c *ui.Context, a *app) {
	components.Progress(c, 0.6)
}

func previewSpinner(c *ui.Context, a *app) {
	components.Spinner(c, components.SpinnerRing)
}

func previewSheetButton(c *ui.Context, a *app) {
	if components.Button(c, a.text("sampleOpenSheet"), components.ButtonSecondary).Height(touchTarget(a)).Clicked() {
		a.sheet = true
	}
}

func previewSafeArea(c *ui.Context, a *app) {
	components.Switch(c, &a.overlay, a.text("overlay"))
}

// touchTarget is the minimum height of a touch control: 44 points on iOS and
// 48 dp on Android.
func touchTarget(a *app) float32 {
	if devices[a.device].platform == platformAndroid {
		return 48
	}
	return 44
}
