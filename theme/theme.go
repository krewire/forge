package theme

import (
	"html/template"
	"strconv"
	"strings"
)

// defaultThemeKey is the localStorage key used when Theme.StorageKey is empty.
const defaultThemeKey = "krewire-theme"

// Theme configures the light/dark theming system for an application or site.
// It follows the system color scheme by default, persists the user's choice,
// and switches between light and dark without server involvement.
type Theme struct {
	// Name is an optional identifier for the theme.
	Name string `json:"name,omitempty"`
	// StorageKey is the localStorage key persisting the user's choice.
	// Defaults to "krewire-theme".
	StorageKey string `json:"storage_key,omitempty"`
	// Default is the initial preference when nothing is stored:
	// "auto", "light", or "dark". Defaults to "auto".
	Default string `json:"default,omitempty"`
	// Light is the color palette for light mode. Empty fields fall back to
	// DefaultLightPalette.
	Light Palette `json:"light,omitempty"`
	// Dark is the color palette for dark mode. Empty fields fall back to
	// DefaultDarkPalette.
	Dark Palette `json:"dark,omitempty"`
}

// Default returns a new Theme initialized with default light and dark palettes.
func Default() *Theme {
	return &Theme{
		Name:    "krewire",
		Default: "auto",
		Light:   DefaultLightPalette,
		Dark:    DefaultDarkPalette,
	}
}

// New returns a new Theme with default settings.
func New() *Theme {
	return Default()
}

// WithStorageKey sets a custom localStorage key.
func (t *Theme) WithStorageKey(key string) *Theme {
	t.StorageKey = key
	return t
}

// WithDefault sets the default mode ("auto", "light", or "dark").
func (t *Theme) WithDefault(def string) *Theme {
	t.Default = def
	return t
}

// StorageKeyOrDefault returns the storage key or the default "krewire-theme".
func (t Theme) StorageKeyOrDefault() string {
	if t.StorageKey == "" {
		return defaultThemeKey
	}
	return t.StorageKey
}

// DefaultTheme returns the configured default theme ("light", "dark", or "auto").
func (t Theme) DefaultTheme() string {
	switch t.Default {
	case "light", "dark":
		return t.Default
	default:
		return "auto"
	}
}

// Style returns the <style> block declaring the color palette as CSS custom
// properties on :root (light) and :root[data-theme="dark"] (dark), plus the
// system-scheme fallback for pages without the attribute before the script
// runs. Consumers may reference the tokens as var(--base-1),
// var(--primary-content), and so on.
func (t Theme) Style() template.CSS {
	light := ":root{" + t.Light.CSSVars(DefaultLightPalette) + "}"
	dark := ":root[data-theme=\"dark\"]{" + t.Dark.CSSVars(DefaultDarkPalette) + "}"
	fallback := "@media (prefers-color-scheme: dark){:root:not([data-theme]){" +
		t.Dark.CSSVars(DefaultDarkPalette) + "}}"
	return template.CSS("<style>" + light + dark + fallback + "</style>")
}

// Script returns the inline <script> block for the document head, preceded by
// the palette style. It applies the stored or system theme before first paint
// (avoiding a flash of the wrong theme), keeps the color-scheme meta in sync,
// and wires any element carrying the data-theme-toggle attribute to cycle
// light and dark modes.
func (t Theme) Script() template.HTML {
	return template.HTML(string(t.Style()) + "<script>" + themeJS(t.StorageKeyOrDefault(), t.DefaultTheme()) + "</script>")
}

// Button returns the theme-switcher button markup: a sun/moon toggle that the
// theme script wires automatically. Place it anywhere in the page.
func (t Theme) Button() template.HTML {
	return template.HTML(`<button type="button" class="theme-toggle" data-theme-toggle aria-label="Toggle light/dark theme" title="Toggle theme"><svg class="icon-sun" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg><svg class="icon-moon" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg></button>`)
}

// TailwindConfigScript returns the inline <script> block configuring TailwindCSS
// with Krewire/Forge design tokens and dark-mode mappings.
func (t Theme) TailwindConfigScript() template.HTML {
	return template.HTML(`<script>
  tailwind.config = {
    darkMode: ['class', '[data-theme="dark"]'],
    theme: {
      extend: {
        colors: {
          primary: 'var(--forge-primary, #39D353)',
          'primary-content': 'var(--forge-primary-content, #0B1F3B)',
          secondary: 'var(--forge-secondary, #00D1C1)',
          'secondary-content': 'var(--forge-secondary-content, #0B1F3B)',
          accent: 'var(--forge-accent, #FF3B2E)',
          surface: 'var(--forge-surface, #f0ede5)',
          muted: 'var(--forge-muted, #475569)',
        }
      }
    }
  }
</script>`)
}

// CSS returns the complete CSS styles for Forge UI, Form, and Panel components
// including base design tokens, mode vars, and pop-brutalist styling.
func (t *Theme) CSS() template.CSS {
	if t == nil {
		t = Default()
	}
	return template.CSS(`
` + ThemeModeVarsCSS + `
` + ThemeToggleCSS + `

:root {
  --forge-primary: var(--primary, ` + string(DefaultLightPalette.Primary) + `);
  --forge-primary-content: var(--primary-content, ` + string(DefaultLightPalette.PrimaryContent) + `);
  --forge-secondary: var(--secondary, ` + string(DefaultLightPalette.Secondary) + `);
  --forge-secondary-content: var(--secondary-content, ` + string(DefaultLightPalette.SecondaryContent) + `);
  --forge-accent: var(--accent, ` + string(DefaultLightPalette.Accent) + `);
  --forge-bg: var(--base-1, ` + string(DefaultLightPalette.Base1) + `);
  --forge-surface: var(--base-2, ` + string(DefaultLightPalette.Base2) + `);
  --forge-fg: var(--base-1-content, ` + string(DefaultLightPalette.Base1Content) + `);
  --forge-muted: var(--neutral-content, ` + string(DefaultLightPalette.NeutralContent) + `);
  --forge-border: var(--base-1-content, ` + string(DefaultLightPalette.Base1Content) + `);
  --forge-success: var(--success, ` + string(DefaultLightPalette.Success) + `);
  --forge-warning: var(--warning, ` + string(DefaultLightPalette.Warning) + `);
  --forge-error: var(--error, ` + string(DefaultLightPalette.Error) + `);
  --forge-radius: 10px;
  --forge-pop-border: 2px solid var(--forge-border);
  --forge-pop-shadow: 3px 3px 0px var(--forge-border);
  --forge-pop-shadow-sm: 2px 2px 0px var(--forge-border);
  --forge-font-sans: Inter, system-ui, -apple-system, sans-serif;
  --forge-font-mono: ui-monospace, SFMono-Regular, Menlo, monospace;
}

[data-theme="dark"], .dark {
  --forge-primary: var(--primary, ` + string(DefaultDarkPalette.Primary) + `);
  --forge-primary-content: var(--primary-content, ` + string(DefaultDarkPalette.PrimaryContent) + `);
  --forge-secondary: var(--secondary, ` + string(DefaultDarkPalette.Secondary) + `);
  --forge-secondary-content: var(--secondary-content, ` + string(DefaultDarkPalette.SecondaryContent) + `);
  --forge-accent: var(--accent, ` + string(DefaultDarkPalette.Accent) + `);
  --forge-bg: var(--base-1, ` + string(DefaultDarkPalette.Base1) + `);
  --forge-surface: var(--base-2, ` + string(DefaultDarkPalette.Base2) + `);
  --forge-fg: var(--base-1-content, ` + string(DefaultDarkPalette.Base1Content) + `);
  --forge-muted: var(--neutral-content, ` + string(DefaultDarkPalette.NeutralContent) + `);
  --forge-border: var(--base-1-content, ` + string(DefaultDarkPalette.Base1Content) + `);
  --forge-success: var(--success, ` + string(DefaultDarkPalette.Success) + `);
  --forge-warning: var(--warning, ` + string(DefaultDarkPalette.Warning) + `);
  --forge-error: var(--error, ` + string(DefaultDarkPalette.Error) + `);
}

.forge-app {
  font-family: var(--forge-font-sans);
  background: var(--forge-bg);
  color: var(--forge-fg);
  line-height: 1.5;
  min-height: 100vh;
}

/* Widgets */
.forge-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.55rem 1.15rem;
  font-size: 0.875rem;
  font-weight: 700;
  border-radius: var(--forge-radius);
  border: var(--forge-pop-border);
  box-shadow: var(--forge-pop-shadow-sm);
  background: var(--forge-surface);
  color: var(--forge-fg);
  cursor: pointer;
  text-decoration: none;
  transition: transform 0.12s ease, box-shadow 0.12s ease;
}
.forge-btn:hover {
  transform: translate(-1px, -1px);
  box-shadow: var(--forge-pop-shadow);
}
.forge-btn-primary {
  background: var(--forge-primary);
  color: var(--forge-primary-content);
}
.forge-btn-secondary {
  background: var(--forge-secondary);
  color: var(--forge-secondary-content);
}
.forge-btn-danger {
  background: var(--forge-error);
  color: #ffffff;
}
.forge-btn-sm { padding: 0.35rem 0.75rem; font-size: 0.78rem; }
.forge-btn-lg { padding: 0.75rem 1.6rem; font-size: 1rem; }

.forge-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.2rem 0.6rem;
  font-size: 0.75rem;
  font-weight: 800;
  border-radius: 999px;
  border: 1.5px solid var(--forge-border);
  background: var(--forge-surface);
  color: var(--forge-fg);
  font-family: var(--forge-font-mono);
}
.forge-badge-primary { background: var(--forge-primary); color: var(--forge-primary-content); }
.forge-badge-secondary { background: var(--forge-secondary); color: var(--forge-secondary-content); }
.forge-badge-success { background: var(--forge-success); color: var(--forge-primary-content); }
.forge-badge-warning { background: var(--forge-warning); color: #000; }
.forge-badge-danger { background: var(--forge-error); color: #fff; }

.forge-card {
  background: var(--forge-bg);
  border: var(--forge-pop-border);
  border-radius: var(--forge-radius);
  box-shadow: var(--forge-pop-shadow);
  padding: 1.25rem;
}
.forge-card-header {
  border-bottom: 1.5px solid var(--forge-border);
  padding-bottom: 0.75rem;
  margin-bottom: 1rem;
}
.forge-card-title {
  font-size: 1.15rem;
  font-weight: 800;
  margin: 0;
}
.forge-card-desc {
  font-size: 0.85rem;
  color: var(--forge-muted);
  margin: 0.25rem 0 0;
}

.forge-alert {
  padding: 0.9rem 1.2rem;
  border-radius: var(--forge-radius);
  border: var(--forge-pop-border);
  box-shadow: var(--forge-pop-shadow-sm);
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  font-weight: 600;
}
.forge-alert-info { background: color-mix(in srgb, var(--forge-secondary) 25%, var(--forge-bg)); }
.forge-alert-success { background: color-mix(in srgb, var(--forge-success) 25%, var(--forge-bg)); }
.forge-alert-warning { background: color-mix(in srgb, var(--forge-warning) 25%, var(--forge-bg)); }
.forge-alert-danger { background: color-mix(in srgb, var(--forge-error) 25%, var(--forge-bg)); }

/* Forms */
.forge-form {
  display: flex;
  flex-direction: column;
  gap: 1.15rem;
}
.forge-field-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}
.forge-label {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--forge-fg);
}
.forge-required {
  color: var(--forge-error);
}
.forge-input, .forge-select, .forge-textarea {
  width: 100%;
  padding: 0.6rem 0.85rem;
  font-size: 0.9rem;
  font-family: inherit;
  border-radius: calc(var(--forge-radius) - 2px);
  border: var(--forge-pop-border);
  background: var(--forge-surface);
  color: var(--forge-fg);
  box-shadow: inset 1px 1px 0px rgba(0,0,0,0.1);
  outline: none;
  box-sizing: border-box;
}
.forge-input:focus, .forge-select:focus, .forge-textarea:focus {
  border-color: var(--forge-primary);
  box-shadow: 0 0 0 2px var(--forge-primary);
}
.forge-field-error {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--forge-error);
  margin-top: 0.2rem;
}
.forge-field-help {
  font-size: 0.78rem;
  color: var(--forge-muted);
}

/* Panels */
.forge-panel {
  border: var(--forge-pop-border);
  border-radius: var(--forge-radius);
  background: var(--forge-bg);
  box-shadow: var(--forge-pop-shadow);
  overflow: hidden;
}
.forge-panel-head {
  padding: 1rem 1.25rem;
  background: var(--forge-surface);
  border-bottom: var(--forge-pop-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}
.forge-panel-title {
  font-size: 1.1rem;
  font-weight: 800;
  margin: 0;
}
.forge-panel-body {
  padding: 1.25rem;
}

/* Stats */
.forge-stat {
  padding: 1.25rem;
  border: var(--forge-pop-border);
  border-radius: var(--forge-radius);
  background: var(--forge-bg);
  box-shadow: var(--forge-pop-shadow-sm);
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}
.forge-stat-label {
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--forge-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.forge-stat-value {
  font-size: 1.75rem;
  font-weight: 900;
  color: var(--forge-fg);
}
.forge-stat-change {
  font-size: 0.8rem;
  font-weight: 700;
}
.forge-stat-change.up { color: var(--forge-success); }
.forge-stat-change.down { color: var(--forge-error); }

/* Tables */
.forge-table-wrap {
  width: 100%;
  overflow-x: auto;
  border: var(--forge-pop-border);
  border-radius: var(--forge-radius);
  background: var(--forge-bg);
}
.forge-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.875rem;
}
.forge-table th {
  padding: 0.75rem 1rem;
  background: var(--forge-surface);
  font-weight: 800;
  border-bottom: var(--forge-pop-border);
  color: var(--forge-fg);
}
.forge-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--forge-border);
}
.forge-table tr:hover {
  background: var(--forge-surface);
}

/* Layouts */
.forge-stack-v { display: flex; flex-direction: column; gap: 1rem; }
.forge-stack-h { display: flex; flex-direction: row; align-items: center; gap: 1rem; flex-wrap: wrap; }
.forge-grid-2 { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 1rem; }
.forge-grid-3 { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1rem; }
.forge-grid-4 { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; }
`)
}

// themeJS generates the plain-JavaScript theming runtime.
func themeJS(key, def string) string {
	var b strings.Builder
	b.WriteString(`(function(){var k=` + strconv.Quote(key) + `,d=` + strconv.Quote(def) + `,t=d;try{t=localStorage.getItem(k)||d}catch(e){}var m=window.matchMedia('(prefers-color-scheme: dark)');function apply(){var dark=t==='dark'||(t==='auto'&&m.matches);var r=document.documentElement;r.dataset.theme=dark?'dark':'light';r.style.colorScheme=dark?'dark':'light';var meta=document.querySelector('meta[name="color-scheme"]');if(meta)meta.content=dark?'dark':'light'}apply();if(m.addEventListener){m.addEventListener('change',function(){if(t==='auto')apply()})}window.krewireTheme={get:function(){return t},set:function(x){t=x;try{localStorage.setItem(k,x)}catch(e){}apply()},toggle:function(){var dark=t==='dark'||(t==='auto'&&m.matches);window.krewireTheme.set(dark?'light':'dark')}};document.addEventListener('click',function(e){var el=e.target&&e.target.closest?e.target.closest('[data-theme-toggle]'):null;if(el)window.krewireTheme.toggle()})})();`)
	return b.String()
}
