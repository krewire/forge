package form

import (
	"fmt"
	"html/template"
	"net/mail"
	"strconv"
	"strings"
)

// Field represents a form input field.
type Field interface {
	Name() string
	Label() string
	Validate(value string) []string
	Render(value string, errors []string) template.HTML
}

// BaseField provides shared attributes for inputs.
type BaseField struct {
	FieldName   string
	FieldLabel  string
	Help        string
	Required    bool
	Placeholder string
	Validators  []func(string) string
}

func (b *BaseField) Name() string  { return b.FieldName }
func (b *BaseField) Label() string { return b.FieldLabel }

func (b *BaseField) Validate(value string) []string {
	var errs []string
	if b.Required && strings.TrimSpace(value) == "" {
		errs = append(errs, b.FieldLabel+" is required")
		return errs
	}
	for _, v := range b.Validators {
		if msg := v(value); msg != "" {
			errs = append(errs, msg)
		}
	}
	return errs
}

// TextField represents a single-line text input.
type TextField struct {
	BaseField
	InputType string // text, email, password, etc.
}

// NewText creates a new text input field.
func NewText(name, label string) *TextField {
	return &TextField{
		BaseField: BaseField{FieldName: name, FieldLabel: label},
		InputType: "text",
	}
}

// MakeRequired marks the field as mandatory.
func (f *TextField) MakeRequired() *TextField {
	f.Required = true
	return f
}

// WithPlaceholder sets placeholder text.
func (f *TextField) WithPlaceholder(ph string) *TextField {
	f.Placeholder = ph
	return f
}

// WithHelp sets field help text.
func (f *TextField) WithHelp(help string) *TextField {
	f.Help = help
	return f
}

// AsEmail configures the field as email type.
func (f *TextField) AsEmail() *TextField {
	f.InputType = "email"
	f.Validators = append(f.Validators, func(val string) string {
		if val == "" {
			return ""
		}
		if _, err := mail.ParseAddress(val); err != nil {
			return "Invalid email address"
		}
		return ""
	})
	return f
}

// AsPassword configures the field as password type.
func (f *TextField) AsPassword() *TextField {
	f.InputType = "password"
	return f
}

// MinLength enforces minimum character length.
func (f *TextField) MinLength(min int) *TextField {
	f.Validators = append(f.Validators, func(val string) string {
		if val != "" && len(val) < min {
			return fmt.Sprintf("Must be at least %d characters", min)
		}
		return ""
	})
	return f
}

// Render produces the field markup.
func (f *TextField) Render(value string, errors []string) template.HTML {
	var b strings.Builder
	b.WriteString(`<div class="forge-field-group">`)
	b.WriteString(`<label class="forge-label" for="field-` + template.HTMLEscapeString(f.FieldName) + `">` + template.HTMLEscapeString(f.FieldLabel))
	if f.Required {
		b.WriteString(` <span class="forge-required">*</span>`)
	}
	b.WriteString(`</label>`)

	b.WriteString(`<input class="forge-input" id="field-` + template.HTMLEscapeString(f.FieldName) + `" name="` + template.HTMLEscapeString(f.FieldName) + `" type="` + template.HTMLEscapeString(f.InputType) + `" value="` + template.HTMLEscapeString(value) + `"`)
	if f.Placeholder != "" {
		b.WriteString(` placeholder="` + template.HTMLEscapeString(f.Placeholder) + `"`)
	}
	if f.Required {
		b.WriteString(` required`)
	}
	b.WriteString(` />`)

	if f.Help != "" {
		b.WriteString(`<span class="forge-field-help">` + template.HTMLEscapeString(f.Help) + `</span>`)
	}
	for _, e := range errors {
		b.WriteString(`<span class="forge-field-error">` + template.HTMLEscapeString(e) + `</span>`)
	}
	b.WriteString(`</div>`)
	return template.HTML(b.String())
}

// NumberField represents a numeric input.
type NumberField struct {
	BaseField
	Min *float64
	Max *float64
}

// NewNumber creates a numeric input field.
func NewNumber(name, label string) *NumberField {
	nf := &NumberField{BaseField: BaseField{FieldName: name, FieldLabel: label}}
	nf.Validators = append(nf.Validators, func(val string) string {
		if val == "" {
			return ""
		}
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return "Must be a valid number"
		}
		if nf.Min != nil && num < *nf.Min {
			return fmt.Sprintf("Must be at least %v", *nf.Min)
		}
		if nf.Max != nil && num > *nf.Max {
			return fmt.Sprintf("Must be at most %v", *nf.Max)
		}
		return ""
	})
	return nf
}

// SetMin sets the minimum value.
func (f *NumberField) SetMin(min float64) *NumberField {
	f.Min = &min
	return f
}

// SetMax sets the maximum value.
func (f *NumberField) SetMax(max float64) *NumberField {
	f.Max = &max
	return f
}

// Render produces the numeric input markup.
func (f *NumberField) Render(value string, errors []string) template.HTML {
	var b strings.Builder
	b.WriteString(`<div class="forge-field-group">`)
	b.WriteString(`<label class="forge-label" for="field-` + template.HTMLEscapeString(f.FieldName) + `">` + template.HTMLEscapeString(f.FieldLabel))
	if f.Required {
		b.WriteString(` <span class="forge-required">*</span>`)
	}
	b.WriteString(`</label>`)
	b.WriteString(`<input class="forge-input" id="field-` + template.HTMLEscapeString(f.FieldName) + `" name="` + template.HTMLEscapeString(f.FieldName) + `" type="number" value="` + template.HTMLEscapeString(value) + `"`)
	if f.Min != nil {
		b.WriteString(fmt.Sprintf(` min="%v"`, *f.Min))
	}
	if f.Max != nil {
		b.WriteString(fmt.Sprintf(` max="%v"`, *f.Max))
	}
	b.WriteString(` />`)
	for _, e := range errors {
		b.WriteString(`<span class="forge-field-error">` + template.HTMLEscapeString(e) + `</span>`)
	}
	b.WriteString(`</div>`)
	return template.HTML(b.String())
}

// TextareaField represents a multi-line text input.
type TextareaField struct {
	BaseField
	Rows int
}

// NewTextarea creates a textarea field.
func NewTextarea(name, label string) *TextareaField {
	return &TextareaField{BaseField: BaseField{FieldName: name, FieldLabel: label}, Rows: 4}
}

// Render produces textarea markup.
func (f *TextareaField) Render(value string, errors []string) template.HTML {
	var b strings.Builder
	b.WriteString(`<div class="forge-field-group">`)
	b.WriteString(`<label class="forge-label" for="field-` + template.HTMLEscapeString(f.FieldName) + `">` + template.HTMLEscapeString(f.FieldLabel) + `</label>`)
	b.WriteString(fmt.Sprintf(`<textarea class="forge-textarea" id="field-%s" name="%s" rows="%d">%s</textarea>`,
		template.HTMLEscapeString(f.FieldName), template.HTMLEscapeString(f.FieldName), f.Rows, template.HTMLEscapeString(value)))
	for _, e := range errors {
		b.WriteString(`<span class="forge-field-error">` + template.HTMLEscapeString(e) + `</span>`)
	}
	b.WriteString(`</div>`)
	return template.HTML(b.String())
}

// SelectOption represents a key-value choice.
type SelectOption struct {
	Value string
	Label string
}

// SelectField represents a dropdown choice input.
type SelectField struct {
	BaseField
	Options []SelectOption
}

// NewSelect creates a dropdown select field.
func NewSelect(name, label string, options ...SelectOption) *SelectField {
	return &SelectField{BaseField: BaseField{FieldName: name, FieldLabel: label}, Options: options}
}

// AddOption adds a choice.
func (f *SelectField) AddOption(value, label string) *SelectField {
	f.Options = append(f.Options, SelectOption{Value: value, Label: label})
	return f
}

// Render produces select markup.
func (f *SelectField) Render(value string, errors []string) template.HTML {
	var b strings.Builder
	b.WriteString(`<div class="forge-field-group">`)
	b.WriteString(`<label class="forge-label" for="field-` + template.HTMLEscapeString(f.FieldName) + `">` + template.HTMLEscapeString(f.FieldLabel) + `</label>`)
	b.WriteString(`<select class="forge-select" id="field-` + template.HTMLEscapeString(f.FieldName) + `" name="` + template.HTMLEscapeString(f.FieldName) + `">`)
	for _, opt := range f.Options {
		sel := ""
		if opt.Value == value {
			sel = " selected"
		}
		b.WriteString(`<option value="` + template.HTMLEscapeString(opt.Value) + `"` + sel + `>` + template.HTMLEscapeString(opt.Label) + `</option>`)
	}
	b.WriteString(`</select>`)
	for _, e := range errors {
		b.WriteString(`<span class="forge-field-error">` + template.HTMLEscapeString(e) + `</span>`)
	}
	b.WriteString(`</div>`)
	return template.HTML(b.String())
}

// CheckboxField represents a boolean toggle checkbox.
type CheckboxField struct {
	BaseField
}

// NewCheckbox creates a checkbox field.
func NewCheckbox(name, label string) *CheckboxField {
	return &CheckboxField{BaseField: BaseField{FieldName: name, FieldLabel: label}}
}

// Render produces checkbox markup.
func (f *CheckboxField) Render(value string, errors []string) template.HTML {
	chk := ""
	if value == "true" || value == "on" || value == "1" {
		chk = " checked"
	}
	return template.HTML(`<div class="forge-field-group" style="flex-direction:row; align-items:center; gap:0.5rem;">` +
		`<input type="checkbox" id="field-` + template.HTMLEscapeString(f.FieldName) + `" name="` + template.HTMLEscapeString(f.FieldName) + `" value="true"` + chk + ` />` +
		`<label class="forge-label" for="field-` + template.HTMLEscapeString(f.FieldName) + `">` + template.HTMLEscapeString(f.FieldLabel) + `</label>` +
		`</div>`)
}
