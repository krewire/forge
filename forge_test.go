package forge

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/krewire/forge/form"
	"github.com/krewire/forge/panel"
	"github.com/krewire/forge/widget"
)

func TestAppPageRouting(t *testing.T) {
	app := New("Forge Admin").
		AddNav("Dashboard", "/").
		AddNav("Settings", "/settings").
		Page("/", "Dashboard", widget.NewCard("Welcome to Admin")).
		Page("/settings", "Settings", widget.NewAlert("System settings").Success())

	// Test GET /
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Dashboard — Forge Admin") {
		t.Errorf("expected title, got %s", body)
	}
	if !strings.Contains(body, "Welcome to Admin") {
		t.Errorf("expected card content, got %s", body)
	}
	if !strings.Contains(body, "Settings") {
		t.Errorf("expected nav link, got %s", body)
	}

	// Test 404
	req404 := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec404 := httptest.NewRecorder()
	app.ServeHTTP(rec404, req404)
	if rec404.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec404.Code)
	}
}

func TestAppFormRouting(t *testing.T) {
	contactForm := form.New("/contact", "POST").
		Add(form.NewText("name", "Your Name").MakeRequired()).
		Add(form.NewText("message", "Message").MakeRequired())

	submitted := false
	app := New("Portal").
		Form("/contact", contactForm, func(vals map[string]string) (widget.Widget, error) {
			submitted = true
			return widget.NewAlert("Thank you, " + vals["name"]).Success(), nil
		})

	// GET form
	reqGet := httptest.NewRequest(http.MethodGet, "/contact", nil)
	recGet := httptest.NewRecorder()
	app.ServeHTTP(recGet, reqGet)
	if recGet.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", recGet.Code)
	}
	if !strings.Contains(recGet.Body.String(), "Your Name") {
		t.Errorf("expected form markup, got %s", recGet.Body.String())
	}

	// POST form valid
	formData := url.Values{
		"name":    {"Budi"},
		"message": {"Hello Forge!"},
	}
	reqPost := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(formData.Encode()))
	reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recPost := httptest.NewRecorder()
	app.ServeHTTP(recPost, reqPost)

	if !submitted {
		t.Error("expected onSubmit handler to have been invoked")
	}
	if !strings.Contains(recPost.Body.String(), "Thank you, Budi") {
		t.Errorf("expected success alert, got %s", recPost.Body.String())
	}
}

func TestAppDashboard(t *testing.T) {
	dash := panel.NewDashboard("Overview").
		AddStat(panel.NewStat("Users", "250")).
		AddPanel(panel.New("System Health").WithBody(widget.NewBadge("Normal").Success()))

	app := New("Monitor").Dashboard("/overview", dash)

	req := httptest.NewRequest(http.MethodGet, "/overview", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Overview — Monitor") {
		t.Errorf("missing title in %s", body)
	}
	if !strings.Contains(body, "Users") || !strings.Contains(body, "250") {
		t.Errorf("missing stat in %s", body)
	}
}

func TestAppTailwindDefaultAndOverride(t *testing.T) {
	appDefault := New("Default App")
	shellDefault := string(appDefault.RenderShell("Home", ""))
	if !strings.Contains(shellDefault, "https://cdn.tailwindcss.com") {
		t.Errorf("expected default Tailwind CDN injection, got %s", shellDefault)
	}
	if !strings.Contains(shellDefault, "tailwind.config") {
		t.Errorf("expected default tailwind config, got %s", shellDefault)
	}

	appCustom := New("Custom App").WithTailwindURL("/custom/tailwind.css")
	shellCustom := string(appCustom.RenderShell("Home", ""))
	if !strings.Contains(shellCustom, `<link rel="stylesheet" href="/custom/tailwind.css">`) {
		t.Errorf("expected custom stylesheet link, got %s", shellCustom)
	}

	appDisabled := New("No Tailwind").DisableTailwind()
	shellDisabled := string(appDisabled.RenderShell("Home", ""))
	if strings.Contains(shellDisabled, "tailwindcss") {
		t.Errorf("expected no Tailwind when disabled, got %s", shellDisabled)
	}
}
