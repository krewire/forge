package theme

import (
	"strings"
	"testing"
)

func TestThemeScriptDefaults(t *testing.T) {
	s := string(Theme{}.Script())
	for _, want := range []string{
		"<script>",
		"krewire-theme",
		"localStorage",
		"(prefers-color-scheme: dark)",
		"data-theme-toggle",
		"auto",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("default script missing %q", want)
		}
	}
}

func TestThemeScriptCustom(t *testing.T) {
	s := string(Theme{StorageKey: "site-theme", Default: "dark"}.Script())
	if !strings.Contains(s, `k="site-theme"`) {
		t.Errorf("script missing custom storage key: %q", s)
	}
	if !strings.Contains(s, `d="dark"`) {
		t.Errorf("script missing custom default: %q", s)
	}
}

func TestThemeButton(t *testing.T) {
	b := string(Theme{}.Button())
	for _, want := range []string{
		`data-theme-toggle`,
		`aria-label="Toggle light/dark theme"`,
		`icon-sun`,
		`icon-moon`,
	} {
		if !strings.Contains(b, want) {
			t.Errorf("button missing %q", want)
		}
	}
}

func TestThemeStorageKeyAndDefault(t *testing.T) {
	if got := (Theme{}).StorageKeyOrDefault(); got != "krewire-theme" {
		t.Errorf("default storage key = %q", got)
	}
	if got := (Theme{Default: "bogus"}).DefaultTheme(); got != "auto" {
		t.Errorf("invalid default = %q, want auto", got)
	}
	if got := (Theme{Default: "light"}).DefaultTheme(); got != "light" {
		t.Errorf("default = %q, want light", got)
	}
}

func TestThemeCSS(t *testing.T) {
	th := Default()
	css := string(th.CSS())
	if !strings.Contains(css, "--forge-primary") {
		t.Errorf("CSS missing --forge-primary")
	}
	if !strings.Contains(css, ".forge-btn") {
		t.Errorf("CSS missing .forge-btn")
	}
}

func TestThemeTailwindConfigScript(t *testing.T) {
	th := Default()
	s := string(th.TailwindConfigScript())
	if !strings.Contains(s, "tailwind.config") {
		t.Errorf("missing tailwind.config in %s", s)
	}
	if !strings.Contains(s, "var(--forge-primary") {
		t.Errorf("missing forge color token mapping in %s", s)
	}
}
