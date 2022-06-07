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
	label.Pos = l.pos
	label.HAlign = ge.AlignCenterHorizontal
	label.VAlign = ge.AlignCenter
	scene.AddGraphics(label)
}

func (l *label) Update(delta float64) {}

func (l *label) IsDisposed() bool { return false }

type checkboxButton struct {
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
	sprite := scene.NewSprite(ImageMenuCheckboxButton)
	sprite.Pos = b.pos
	scene.AddGraphics(sprite)

	label := scene.NewLabel(FontDescription)
	label.Text = b.Text
	label.Pos = b.pos.WithOffset(32, 0)
	label.HAlign = ge.AlignCenterHorizontal
	label.VAlign = ge.AlignCenter
	scene.AddGraphics(label)
	b.label = label

	b.status = scene.NewLabel(FontDescription)
	b.status.HAlign = ge.AlignCenterHorizontal
	b.status.VAlign = ge.AlignCenter
	b.status.Pos = b.pos.WithOffset(-120, 0)
	scene.AddGraphics(b.status)

	b.updateColor()
}

func (b *checkboxButton) Update(delta float64) {
	b.label.Text = b.Text
	if b.checked {
		b.status.Text = "ON"
	} else {
		b.status.Text = "OFF"
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
	sprite := scene.NewSprite(ImageMenuSelectButton)
	sprite.Pos = b.pos
	scene.AddGraphics(sprite)

	label := scene.NewLabel(FontBig)
	label.Text = b.options[b.selected.Value()]
	label.Pos = b.pos
	label.HAlign = ge.AlignCenterHorizontal
	label.VAlign = ge.AlignCenter
	scene.AddGraphics(label)
	b.label = label
	b.updateColor()
}

func (b *selectButton) Update(delta float64) {
	b.label.Text = b.options[b.selected.Value()]
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
	label.Text = b.Text
	label.Pos = b.pos
	label.HAlign = ge.AlignCenterHorizontal
	label.VAlign = ge.AlignCenter
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
	d.ironPrice.HAlign = ge.AlignCenterHorizontal
	d.ironPrice.Pos = d.Pos.WithOffset(14, 38)
	scene.AddGraphics(d.ironPrice)

	d.goldPrice = scene.NewLabel(FontSmall)
	d.goldPrice.HAlign = ge.AlignCenterHorizontal
	d.goldPrice.Pos = d.Pos.WithOffset(14+40, 38)
	scene.AddGraphics(d.goldPrice)

	d.oilPrice = scene.NewLabel(FontSmall)
	d.oilPrice.HAlign = ge.AlignCenterHorizontal
	d.oilPrice.Pos = d.Pos.WithOffset(14+80, 38)
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
