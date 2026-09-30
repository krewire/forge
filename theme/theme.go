package theme

import "html/template"

// Theme contains design tokens and configuration for Forge components.
type Theme struct {
	Name        string `json:"name"`
	Primary     string `json:"primary"`
	Secondary   string `json:"secondary"`
	Accent      string `json:"accent"`
	Background  string `json:"bg"`
	Surface     string `json:"surface"`
	Foreground  string `json:"fg"`
	Muted       string `json:"muted"`
	Border      string `json:"border"`
	Success     string `json:"success"`
	Warning     string `json:"warning"`
	Error       string `json:"error"`
	Radius      string `json:"radius"`
}

// Default returns the default pop-brutalist theme matching Krewire palette.
func Default() *Theme {
	return &Theme{
		Name:       "krewire",
		Primary:    "#39D353",
		Secondary:  "#00D1C1",
		Accent:     "#FF3B2E",
		Background: "#FAF8F4",
		Surface:    "#f0ede5",
		Foreground: "#0B1F3B",
		Muted:      "#475569",
		Border:     "#0B1F3B",
		Success:    "#39D353",
		Warning:    "#ffb703",
		Error:      "#FF3B2E",
		Radius:     "10px",
	}
}

// CSS returns the complete CSS styles for Forge UI, Form, and Panel components.
func (t *Theme) CSS() template.CSS {
	if t == nil {
		t = Default()
	}
	return template.CSS(`
:root {
  --forge-primary: ` + t.Primary + `;
  --forge-primary-content: #0B1F3B;
  --forge-secondary: ` + t.Secondary + `;
  --forge-secondary-content: #0B1F3B;
  --forge-accent: ` + t.Accent + `;
  --forge-bg: ` + t.Background + `;
  --forge-surface: ` + t.Surface + `;
  --forge-fg: ` + t.Foreground + `;
  --forge-muted: ` + t.Muted + `;
  --forge-border: ` + t.Border + `;
  --forge-success: ` + t.Success + `;
  --forge-warning: ` + t.Warning + `;
  --forge-error: ` + t.Error + `;
  --forge-radius: ` + t.Radius + `;
  --forge-pop-border: 2px solid var(--forge-border);
  --forge-pop-shadow: 3px 3px 0px var(--forge-border);
  --forge-pop-shadow-sm: 2px 2px 0px var(--forge-border);
  --forge-font-sans: Inter, system-ui, -apple-system, sans-serif;
  --forge-font-mono: ui-monospace, SFMono-Regular, Menlo, monospace;
}

[data-theme="dark"], .dark {
  --forge-primary: #39D353;
  --forge-primary-content: #0B1F3B;
  --forge-secondary: #00D1C1;
  --forge-secondary-content: #0B1F3B;
  --forge-bg: #0B1F3B;
  --forge-surface: #112646;
  --forge-fg: #FAF8F4;
  --forge-muted: #8fa2bc;
  --forge-border: #FAF8F4;
  --forge-pop-border: 2px solid var(--forge-border);
  --forge-pop-shadow: 3px 3px 0px var(--forge-border);
  --forge-pop-shadow-sm: 2px 2px 0px var(--forge-border);
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
