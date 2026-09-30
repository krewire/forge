package widget

import (
	"html/template"
	"strings"
)

// Card represents a structured container panel.
type Card struct {
	Title       string
	Description string
	Body        Widget
	Footer      Widget
	Actions     Widget
}

// NewCard creates a new card with title.
func NewCard(title string) *Card {
	return &Card{Title: title}
}

// WithDesc sets a card description.
func (c *Card) WithDesc(desc string) *Card {
	c.Description = desc
	return c
}

// WithBody sets the card's inner content.
func (c *Card) WithBody(body Widget) *Card {
	c.Body = body
	return c
}

// WithFooter sets the card footer.
func (c *Card) WithFooter(footer Widget) *Card {
	c.Footer = footer
	return c
}

// WithActions sets header action widgets.
func (c *Card) WithActions(actions Widget) *Card {
	c.Actions = actions
	return c
}

// Render produces the card HTML.
func (c *Card) Render() template.HTML {
	var buf strings.Builder
	buf.WriteString(`<div class="forge-card">`)

	if c.Title != "" || c.Description != "" || c.Actions != nil {
		buf.WriteString(`<div class="forge-card-header" style="display:flex; justify-content:space-between; align-items:flex-start;">`)
		buf.WriteString(`<div>`)
		if c.Title != "" {
			buf.WriteString(`<h3 class="forge-card-title">` + template.HTMLEscapeString(c.Title) + `</h3>`)
		}
		if c.Description != "" {
			buf.WriteString(`<p class="forge-card-desc">` + template.HTMLEscapeString(c.Description) + `</p>`)
		}
		buf.WriteString(`</div>`)
		if c.Actions != nil {
			buf.WriteString(`<div>` + string(c.Actions.Render()) + `</div>`)
		}
		buf.WriteString(`</div>`)
	}

	if c.Body != nil {
		buf.WriteString(`<div class="forge-card-body">` + string(c.Body.Render()) + `</div>`)
	}

	if c.Footer != nil {
		buf.WriteString(`<div class="forge-card-footer" style="border-top:1.5px solid var(--forge-border); padding-top:0.75rem; margin-top:1rem;">` + string(c.Footer.Render()) + `</div>`)
	}

	buf.WriteString(`</div>`)
	return template.HTML(buf.String())
}

// Alert renders an informational or status callout box.
type Alert struct {
	Message string
	Variant string // info, success, warning, danger
	Icon    string
}

// NewAlert creates a new alert box.
func NewAlert(message string) *Alert {
	return &Alert{Message: message, Variant: "info", Icon: "ℹ"}
}

// Success sets variant to success.
func (a *Alert) Success() *Alert { a.Variant = "success"; a.Icon = "✓"; return a }

// Warning sets variant to warning.
func (a *Alert) Warning() *Alert { a.Variant = "warning"; a.Icon = "⚠"; return a }

// Danger sets variant to danger.
func (a *Alert) Danger() *Alert { a.Variant = "danger"; a.Icon = "✕"; return a }

// Render produces the alert HTML.
func (a *Alert) Render() template.HTML {
	return template.HTML(`<div class="forge-alert forge-alert-` + template.HTMLEscapeString(a.Variant) + `">` +
		`<span class="alert-icon" style="font-weight:900;">` + template.HTMLEscapeString(a.Icon) + `</span>` +
		`<span>` + template.HTMLEscapeString(a.Message) + `</span>` +
		`</div>`)
}

// Stack represents a layout container that arranges widgets horizontally or vertically.
type Stack struct {
	Direction string // vertical or horizontal
	Widgets   []Widget
	Gap       string
	Align     string
}

// VStack creates a vertical stack.
func VStack(widgets ...Widget) *Stack {
	return &Stack{Direction: "vertical", Widgets: widgets, Gap: "1rem"}
}

// HStack creates a horizontal stack.
func HStack(widgets ...Widget) *Stack {
	return &Stack{Direction: "horizontal", Widgets: widgets, Gap: "1rem", Align: "center"}
}

// WithGap sets the gap size.
func (s *Stack) WithGap(gap string) *Stack {
	s.Gap = gap
	return s
}

// Render produces the stack layout.
func (s *Stack) Render() template.HTML {
	var buf strings.Builder
	dirCls := "forge-stack-v"
	if s.Direction == "horizontal" {
		dirCls = "forge-stack-h"
	}
	style := ""
	if s.Gap != "" {
		style = ` style="gap:` + template.HTMLEscapeString(s.Gap) + `;"`
	}
	buf.WriteString(`<div class="` + dirCls + `"` + style + `>`)
	for _, w := range s.Widgets {
		if w != nil {
			buf.WriteString(string(w.Render()))
		}
	}
	buf.WriteString(`</div>`)
	return template.HTML(buf.String())
}

// Grid arranges items in a responsive grid.
type Grid struct {
	Columns int
	Widgets []Widget
}

// NewGrid creates a grid with column count (2, 3, 4).
func NewGrid(columns int, widgets ...Widget) *Grid {
	if columns < 1 || columns > 4 {
		columns = 2
	}
	return &Grid{Columns: columns, Widgets: widgets}
}

// Render produces the grid layout.
func (g *Grid) Render() template.HTML {
	var buf strings.Builder
	cls := "forge-grid-2"
	if g.Columns == 3 {
		cls = "forge-grid-3"
	} else if g.Columns == 4 {
		cls = "forge-grid-4"
	}
	buf.WriteString(`<div class="` + cls + `">`)
	for _, w := range g.Widgets {
		if w != nil {
			buf.WriteString(string(w.Render()))
		}
	}
	buf.WriteString(`</div>`)
	return template.HTML(buf.String())
}
