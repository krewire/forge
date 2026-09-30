package forge

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/krewire/forge/form"
	"github.com/krewire/forge/panel"
	"github.com/krewire/forge/theme"
	"github.com/krewire/forge/widget"
)

// NavItem represents a link in the app navigation bar.
type NavItem struct {
	Label string
	URL   string
}

// App represents a programmatic web application built with Forge.
type App struct {
	Title    string
	Theme    *theme.Theme
	Nav      []NavItem
	mux      *http.ServeMux
	routes   map[string]bool
}

// New creates a new Forge programmatic application builder.
func New(title string) *App {
	return &App{
		Title:  title,
		Theme:  theme.Default(),
		mux:    http.NewServeMux(),
		routes: make(map[string]bool),
	}
}

// WithTheme sets custom theme tokens.
func (a *App) WithTheme(t *theme.Theme) *App {
	a.Theme = t
	return a
}

// AddNav adds a navigation item.
func (a *App) AddNav(label, url string) *App {
	a.Nav = append(a.Nav, NavItem{Label: label, URL: url})
	return a
}

// Page registers a simple static or dynamic content page.
func (a *App) Page(path, title string, content widget.Widget) *App {
	a.routes[path] = true
	a.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		body := content.Render()
		fmt.Fprint(w, a.RenderShell(title, body))
	})
	return a
}

// Panel registers a panel page.
func (a *App) Panel(path string, p *panel.Panel) *App {
	return a.Page(path, p.Title, p)
}

// Dashboard registers a dashboard page.
func (a *App) Dashboard(path string, d *panel.Dashboard) *App {
	return a.Page(path, d.Title, d)
}

// Form registers an interactive form route with GET (render) and POST (submit) handlers.
func (a *App) Form(path string, f *form.Form, onSubmit func(values map[string]string) (widget.Widget, error)) *App {
	a.routes[path] = true
	a.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := f.BindRequest(r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if f.Validate() {
				if onSubmit != nil {
					res, err := onSubmit(f.Values())
					if err != nil {
						alert := widget.NewAlert(err.Error()).Danger()
						body := template.HTML(string(alert.Render()) + string(f.Render()))
						w.Header().Set("Content-Type", "text/html; charset=utf-8")
						fmt.Fprint(w, a.RenderShell("Error", body))
						return
					}
					if res != nil {
						w.Header().Set("Content-Type", "text/html; charset=utf-8")
						fmt.Fprint(w, a.RenderShell("Success", res.Render()))
						return
					}
				}
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, a.RenderShell(f.Action, f.Render()))
	})
	return a
}

// RenderShell renders the complete HTML envelope with theme and navigation.
func (a *App) RenderShell(pageTitle string, content template.HTML) template.HTML {
	var b strings.Builder
	title := a.Title
	if pageTitle != "" {
		title = pageTitle + " — " + a.Title
	}

	b.WriteString("<!doctype html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString("  <meta charset=\"utf-8\">\n")
	b.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	b.WriteString("  <title>" + template.HTMLEscapeString(title) + "</title>\n")
	b.WriteString("  <style>\n" + string(a.Theme.CSS()) + "\n</style>\n")
	b.WriteString("</head>\n<body class=\"forge-app\">\n")

	// Top Navigation
	b.WriteString(`<header style="background:var(--forge-surface); border-bottom:var(--forge-pop-border); padding:0.85rem 1.5rem; display:flex; align-items:center; justify-content:space-between;">`)
	b.WriteString(`<div style="display:flex; align-items:center; gap:0.75rem;">`)
	b.WriteString(`<span style="background:var(--forge-primary); color:var(--forge-primary-content); width:28px; height:28px; display:inline-grid; place-items:center; border-radius:6px; font-weight:900; border:var(--forge-pop-border); box-shadow:var(--forge-pop-shadow-sm);">◈</span>`)
	b.WriteString(`<span style="font-weight:900; font-size:1.1rem; letter-spacing:-0.03em;">` + template.HTMLEscapeString(a.Title) + `</span>`)
	b.WriteString(`</div>`)

	if len(a.Nav) > 0 {
		b.WriteString(`<nav style="display:flex; align-items:center; gap:0.5rem;">`)
		for _, item := range a.Nav {
			b.WriteString(`<a href="` + template.HTMLEscapeString(item.URL) + `" style="font-size:0.85rem; font-weight:700; padding:0.4rem 0.75rem; border-radius:6px; color:var(--forge-fg); text-decoration:none;">` + template.HTMLEscapeString(item.Label) + `</a>`)
		}
		b.WriteString(`</nav>`)
	}
	b.WriteString(`</header>\n`)

	// Main Canvas
	b.WriteString(`<main style="max-width:1200px; margin:0 auto; padding:2rem 1.5rem;">\n`)
	b.WriteString(string(content))
	b.WriteString("\n</main>\n")

	b.WriteString("</body>\n</html>")
	return template.HTML(b.String())
}

// ServeHTTP dispatches requests to registered routes.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
