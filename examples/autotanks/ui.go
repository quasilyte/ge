package main

import (
	"strconv"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/xslices"
)

type focusToggler interface {
	ToggleFocus()
}

type label struct {
	text string
	pos  ge.Pos
}

func newLabel(text string, pos ge.Pos) *label {
	return &label{text: text, pos: pos}
}

func (l *label) Init(scene *ge.Scene) {
	bg := ge.NewRect(128, 64)
	bg.Pos = l.pos
	bg.ColorScale.SetRGBA(0x3a, 0x44, 0x66, 200)
	scene.AddGraphics(bg)

	label := scene.NewLabel(FontBig)
	label.Text = l.text
	label.Pos = bg.AnchorPos()
	label.Width = bg.Width
	label.Height = bg.Height
	label.AlignHorizontal = ge.AlignHorizontalCenter
	label.AlignVertical = ge.AlignVerticalCenter
	scene.AddGraphics(label)
}

func (l *label) Update(delta float64) {}

func (l *label) IsDisposed() bool { return false }

type checkboxButton struct {
	scene   *ge.Scene
	status  *ge.Label
	label   *ge.Label
	checked bool
	Text    string
	Focused bool
	pos     ge.Pos
}

func newCheckboxButton(text string, checked bool, pos ge.Pos) *checkboxButton {
	return &checkboxButton{Text: text, checked: checked, pos: pos}
}

func (b *checkboxButton) Init(scene *ge.Scene) {
	b.scene = scene

	sprite := scene.NewSprite(ImageMenuCheckboxButton)
	sprite.Pos = b.pos
	scene.AddGraphics(sprite)

	label := scene.NewLabel(FontDescription)
	label.Text = b.scene.Dict().Get(b.Text)
	label.Pos = sprite.AnchorPos().WithOffset(76, 2)
	label.AlignHorizontal = ge.AlignHorizontalCenter
	label.AlignVertical = ge.AlignVerticalCenter
	label.Width = 242
	label.Height = 58
	scene.AddGraphics(label)
	b.label = label

	b.status = scene.NewLabel(FontDescription)
	b.status.AlignHorizontal = ge.AlignHorizontalCenter
	b.status.AlignVertical = ge.AlignVerticalCenter
	b.status.Width = 68
	b.status.Height = 34
	b.status.Pos = sprite.AnchorPos().WithOffset(5, 14)
	scene.AddGraphics(b.status)

	b.updateColor()
}

func (b *checkboxButton) Update(delta float64) {
	b.label.Text = b.scene.Dict().Get(b.Text)
	if b.checked {
		b.status.Text = b.scene.Dict().Get("ui.on")
	} else {
		b.status.Text = b.scene.Dict().Get("ui.off")
	}
	b.updateColor()
}

func (b *checkboxButton) updateColor() {
	if b.Focused {
		b.label.ColorScale = ge.ColorScale{R: 0.8, G: 1, B: 0.8, A: 1}
	} else {
		b.label.ColorScale = ge.ColorScale{R: 0.6, G: 0.6, B: 0.6, A: 1}
	}
}

func (b *checkboxButton) IsDisposed() bool { return false }

func (b *checkboxButton) ToggleChecked() { b.checked = !b.checked }

func (b *checkboxButton) SetChecked(checked bool) { b.checked = checked }

func (b *checkboxButton) ToggleFocus() { b.Focused = !b.Focused }

type selectButton struct {
	scene    *ge.Scene
	label    *ge.Label
	options  []string
	Focused  bool
	selected gemath.Slider
	pos      ge.Pos
}

func newSelectButton(options []string, pos ge.Pos) *selectButton {
	b := &selectButton{
		options: options,
		pos:     pos,
	}
	b.selected.SetBounds(0, len(b.options)-1)
	return b
}

func (b *selectButton) Init(scene *ge.Scene) {
	b.scene = scene

	sprite := scene.NewSprite(ImageMenuSelectButton)
	sprite.Pos = b.pos
	scene.AddGraphics(sprite)

	label := scene.NewLabel(FontBig)
	label.Text = scene.Dict().Get(b.options[b.selected.Value()])
	label.Pos = sprite.AnchorPos()
	label.AlignHorizontal = ge.AlignHorizontalCenter
	label.AlignVertical = ge.AlignVerticalCenter
	label.Width = sprite.FrameWidth
	label.Height = sprite.FrameHeight
	scene.AddGraphics(label)
	b.label = label
	b.updateColor()
}

func (b *selectButton) Update(delta float64) {
	b.label.Text = b.scene.Dict().Get(b.options[b.selected.Value()])
	b.updateColor()
}

func (b *selectButton) updateColor() {
	if b.Focused {
		b.label.ColorScale = ge.ColorScale{R: 0.8, G: 1, B: 0.8, A: 1}
	} else {
		b.label.ColorScale = ge.ColorScale{R: 0.6, G: 0.6, B: 0.6, A: 1}
	}
}

func (b *selectButton) IsDisposed() bool { return false }

func (b *selectButton) Select(option string) {
	b.selected.TrySetValue(xslices.Index(b.options, option))
}

func (b *selectButton) NextOption() {
	b.selected.Inc()
}

func (b *selectButton) PrevOption() {
	b.selected.Dec()
}

func (b *selectButton) SelectedOption() string {
	return b.options[b.selected.Value()]
}

func (b *selectButton) ToggleFocus() { b.Focused = !b.Focused }

type button struct {
	label *ge.Label

	Focused bool

	Text string

	pos ge.Pos
}

func newButton(text string, pos ge.Pos) *button {
	return &button{
		Text: text,
		pos:  pos,
	}
}

func (b *button) Init(scene *ge.Scene) {
	sprite := scene.NewSprite(ImageMenuButton)
	sprite.Pos = b.pos
	scene.AddGraphics(sprite)

	label := scene.NewLabel(FontBig)
	label.Text = scene.Dict().Get(b.Text)
	label.Pos = sprite.AnchorPos()
	label.AlignHorizontal = ge.AlignHorizontalCenter
	label.AlignVertical = ge.AlignVerticalCenter
	label.Width = sprite.FrameWidth
	label.Height = sprite.FrameHeight
	scene.AddGraphics(label)
	b.label = label
	b.updateColor()
}

func (b *button) Update(delta float64) {
	b.updateColor()
}

func (b *button) updateColor() {
	if b.Focused {
		b.label.ColorScale = ge.ColorScale{R: 0.8, G: 1, B: 0.8, A: 1}
	} else {
		b.label.ColorScale = ge.ColorScale{R: 0.6, G: 0.6, B: 0.6, A: 1}
	}
}

func (b *button) IsDisposed() bool { return false }

func (b *button) ToggleFocus() { b.Focused = !b.Focused }

type priceDisplay struct {
	Pos ge.Pos

	sprite *ge.Sprite

	ironPrice *ge.Label
	goldPrice *ge.Label
	oilPrice  *ge.Label

	price resourceContainer
}

func newPriceDisplay() *priceDisplay {
	return &priceDisplay{}
}

func (d *priceDisplay) Init(scene *ge.Scene) {
	d.sprite = scene.NewSprite(ImageResourceRow)
	d.sprite.Pos = d.Pos
	d.sprite.Centered = false
	scene.AddGraphics(d.sprite)

	d.ironPrice = scene.NewLabel(FontSmall)
	d.ironPrice.AlignHorizontal = ge.AlignHorizontalCenter
	d.ironPrice.Pos = d.Pos.WithOffset(0, 38)
	d.ironPrice.Width = 28
	scene.AddGraphics(d.ironPrice)

	d.goldPrice = scene.NewLabel(FontSmall)
	d.goldPrice.AlignHorizontal = ge.AlignHorizontalCenter
	d.goldPrice.Pos = d.Pos.WithOffset(40, 38)
	d.goldPrice.Width = 28
	scene.AddGraphics(d.goldPrice)

	d.oilPrice = scene.NewLabel(FontSmall)
	d.oilPrice.AlignHorizontal = ge.AlignHorizontalCenter
	d.oilPrice.Pos = d.Pos.WithOffset(80, 38)
	d.oilPrice.Width = 28
	scene.AddGraphics(d.oilPrice)
}

func (d *priceDisplay) IsDisposed() bool { return false }

func (d *priceDisplay) Update(delta float64) {}

func (d *priceDisplay) SetVisibility(v bool) {
	d.sprite.Visible = v
	d.ironPrice.Visible = v
	d.goldPrice.Visible = v
	d.oilPrice.Visible = v
}

func (d *priceDisplay) SetPrice(total resourceContainer) {
	d.ironPrice.Text = strconv.Itoa(total.Iron)
	d.goldPrice.Text = strconv.Itoa(total.Gold)
	d.oilPrice.Text = strconv.Itoa(total.Oil)
}

func (d *priceDisplay) SetAvailable(iron, gold, oil bool) {
	d.setAvailable(d.ironPrice, iron)
	d.setAvailable(d.goldPrice, gold)
	d.setAvailable(d.oilPrice, oil)
}

func (d *priceDisplay) setAvailable(l *ge.Label, available bool) {
	if available {
		l.ColorScale.SetRGBA(255, 255, 255, 255)
	} else {
		l.ColorScale.SetRGBA(255, 180, 180, 255)
	}
}
