package widget

import (
	"html/template"
	"strings"
)

// Button represents a clickable button or anchor.
type Button struct {
	Text     string
	URL      string
	Variant  string // primary, secondary, danger, ghost, default
	Size     string // sm, md, lg
	Icon     string
	Type     string // button, submit, reset
	Disabled bool
	OnClick  string
	Class    string
}

// NewButton creates a new button with label.
func NewButton(text string) *Button {
	return &Button{
		Text:    text,
		Variant: "default",
		Size:    "md",
		Type:    "button",
	}
}

// Primary sets variant to primary.
func (b *Button) Primary() *Button {
	b.Variant = "primary"
	return b
}

// Secondary sets variant to secondary.
func (b *Button) Secondary() *Button {
	b.Variant = "secondary"
	return b
}

// Danger sets variant to danger.
func (b *Button) Danger() *Button {
	b.Variant = "danger"
	return b
}

// Small sets size to sm.
func (b *Button) Small() *Button {
	b.Size = "sm"
	return b
}

// Large sets size to lg.
func (b *Button) Large() *Button {
	b.Size = "lg"
	return b
}

// Link turns the button into an anchor link.
func (b *Button) Link(url string) *Button {
	b.URL = url
	return b
}

// WithIcon sets a leading icon.
func (b *Button) WithIcon(icon string) *Button {
	b.Icon = icon
	return b
}

// Submit marks button as form submit type.
func (b *Button) Submit() *Button {
	b.Type = "submit"
	return b
}

// Render produces the button HTML.
func (b *Button) Render() template.HTML {
	cls := "forge-btn forge-btn-" + b.Variant
	if b.Size != "md" && b.Size != "" {
		cls += " forge-btn-" + b.Size
	}
	if b.Class != "" {
		cls += " " + b.Class
	}

	content := template.HTMLEscapeString(b.Text)
	if b.Icon != "" {
		content = `<span class="btn-icon">` + template.HTMLEscapeString(b.Icon) + `</span> ` + content
	}

	if b.URL != "" {
		return template.HTML(`<a href="` + template.HTMLEscapeString(b.URL) + `" class="` + template.HTMLEscapeString(cls) + `">` + content + `</a>`)
	}

	dis := ""
	if b.Disabled {
		dis = " disabled"
	}
	clk := ""
	if b.OnClick != "" {
		clk = ` onclick="` + template.HTMLEscapeString(b.OnClick) + `"`
	}

	btnType := b.Type
	if btnType == "" {
		btnType = "button"
	}

	return template.HTML(`<button type="` + template.HTMLEscapeString(btnType) + `" class="` + template.HTMLEscapeString(cls) + `"` + clk + dis + `>` + content + `</button>`)
}

// Badge represents an inline status or label tag.
type Badge struct {
	Text    string
	Variant string // primary, secondary, success, warning, danger, default
	Icon    string
}

// NewBadge creates a new badge with text.
func NewBadge(text string) *Badge {
	return &Badge{Text: text, Variant: "default"}
}

// Primary sets badge variant to primary.
func (bg *Badge) Primary() *Badge { bg.Variant = "primary"; return bg }

// Secondary sets badge variant to secondary.
func (bg *Badge) Secondary() *Badge { bg.Variant = "secondary"; return bg }

// Success sets badge variant to success.
func (bg *Badge) Success() *Badge { bg.Variant = "success"; return bg }

// Warning sets badge variant to warning.
func (bg *Badge) Warning() *Badge { bg.Variant = "warning"; return bg }

// Danger sets badge variant to danger.
func (bg *Badge) Danger() *Badge { bg.Variant = "danger"; return bg }

// WithIcon adds an icon to the badge.
func (bg *Badge) WithIcon(icon string) *Badge { bg.Icon = icon; return bg }

// Render produces the badge markup.
func (bg *Badge) Render() template.HTML {
	var buf strings.Builder
	buf.WriteString(`<span class="forge-badge forge-badge-` + template.HTMLEscapeString(bg.Variant) + `">`)
	if bg.Icon != "" {
		buf.WriteString(`<span class="badge-icon">` + template.HTMLEscapeString(bg.Icon) + `</span>`)
	}
	buf.WriteString(template.HTMLEscapeString(bg.Text))
	buf.WriteString(`</span>`)
	return template.HTML(buf.String())
}
