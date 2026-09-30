package widget

import (
	"fmt"
	"html/template"
	"strings"
)

// Widget represents any renderable Forge UI component.
type Widget interface {
	Render() template.HTML
}

// HTML converts a raw HTML string into a Widget.
type HTML string

// Render returns the HTML content.
func (h HTML) Render() template.HTML {
	return template.HTML(h)
}

// Text renders safe escaped text.
type Text struct {
	Content string
	Bold    bool
	Muted   bool
}

// NewText returns a new Text widget.
func NewText(content string) *Text {
	return &Text{Content: content}
}

// Strong marks the text as bold.
func (t *Text) Strong() *Text {
	t.Bold = true
	return t
}

// Subtle marks the text as muted.
func (t *Text) Subtle() *Text {
	t.Muted = true
	return t
}

// Render produces the text markup.
func (t *Text) Render() template.HTML {
	escaped := template.HTMLEscapeString(t.Content)
	if t.Bold {
		escaped = "<strong>" + escaped + "</strong>"
	}
	if t.Muted {
		return template.HTML(`<span style="color:var(--forge-muted)">` + escaped + `</span>`)
	}
	return template.HTML(escaped)
}

// Heading renders an h1-h6 tag.
type Heading struct {
	Level   int
	Content string
}

// NewHeading creates a heading of specified level (1..6).
func NewHeading(level int, content string) *Heading {
	if level < 1 || level > 6 {
		level = 2
	}
	return &Heading{Level: level, Content: content}
}

// Render produces the <hN> element.
func (h *Heading) Render() template.HTML {
	return template.HTML(fmt.Sprintf(`<h%d class="forge-heading" style="margin:0.5rem 0; font-weight:800">%s</h%d>`,
		h.Level, template.HTMLEscapeString(h.Content), h.Level))
}

// Divider renders a horizontal or vertical separator.
type Divider struct {
	Label string
}

// NewDivider creates a new divider.
func NewDivider(label ...string) *Divider {
	l := ""
	if len(label) > 0 {
		l = label[0]
	}
	return &Divider{Label: l}
}

// Render produces the divider.
func (d *Divider) Render() template.HTML {
	if d.Label == "" {
		return `<hr style="border:none; border-top:var(--forge-pop-border); margin:1.25rem 0;" />`
	}
	return template.HTML(fmt.Sprintf(`<div style="display:flex; align-items:center; gap:0.75rem; margin:1.25rem 0;"><hr style="flex:1; border:none; border-top:var(--forge-pop-border);" /><span style="font-size:0.75rem; font-weight:700; color:var(--forge-muted); text-transform:uppercase;">%s</span><hr style="flex:1; border:none; border-top:var(--forge-pop-border);" /></div>`, template.HTMLEscapeString(d.Label)))
}

// Fragment renders multiple widgets sequentially.
func Fragment(widgets ...Widget) template.HTML {
	var b strings.Builder
	for _, w := range widgets {
		if w != nil {
			b.WriteString(string(w.Render()))
		}
	}
	return template.HTML(b.String())
}
