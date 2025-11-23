package validator

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, e := range ve {
		messages = append(messages, e.Field+": "+e.Message)
	}
	return strings.Join(messages, "; ")
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

type Validator struct {
	errors ValidationErrors
}

func New() *Validator {
	return &Validator{
		errors: make(ValidationErrors, 0),
	}
}

func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{Field: field, Message: message})
}

func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "is required")
	}
	return v
}

func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.AddError(field, "must be a valid email address")
	}
	return v
}

func (v *Validator) MinLength(field, value string, min int) *Validator {
	if value == "" {
		return v
	}
	if len(value) < min {
		v.AddError(field, "must be at least "+string(rune('0'+min))+" characters")
	}
	return v
}

func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if value == "" {
		return v
	}
	if len(value) > max {
		v.AddError(field, "must be at most "+string(rune('0'+max))+" characters")
	}
	return v
}

func (v *Validator) Password(field, value string) *Validator {
	if value == "" {
		return v
	}
	if len(value) < 8 {
		v.AddError(field, "must be at least 8 characters")
		return v
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		v.AddError(field, "must contain uppercase, lowercase, and digit")
	}
	return v
}

func (v *Validator) UUID(field, value string) *Validator {
	if value == "" {
		return v
	}
	if _, err := uuid.Parse(value); err != nil {
		v.AddError(field, "must be a valid UUID")
	}
	return v
}

func (v *Validator) DateAfter(field string, value, after time.Time) *Validator {
	if value.IsZero() {
		return v
	}
	if !value.After(after) {
		v.AddError(field, "must be after "+after.Format(time.RFC3339))
	}
	return v
}

func (v *Validator) DateBefore(field string, value, before time.Time) *Validator {
	if value.IsZero() {
		return v
	}
	if !value.Before(before) {
		v.AddError(field, "must be before "+before.Format(time.RFC3339))
	}
	return v
}

func (v *Validator) PositiveNumber(field string, value float64) *Validator {
	if value <= 0 {
		v.AddError(field, "must be a positive number")
	}
	return v
}

func (v *Validator) InList(field, value string, allowed []string) *Validator {
	if value == "" {
		return v
	}
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.AddError(field, "must be one of: "+strings.Join(allowed, ", "))
	return v
}
