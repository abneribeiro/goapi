package validator

import (
	"testing"
)

func TestValidator_Required(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"valid value", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Required("field", tt.value)

			if tt.wantError && !v.Errors().HasErrors() {
				t.Error("expected validation error")
			}
			if !tt.wantError && v.Errors().HasErrors() {
				t.Errorf("unexpected validation error: %v", v.Errors())
			}
		})
	}
}

func TestValidator_Email(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"valid email", "test@example.com", false},
		{"valid email with subdomain", "test@mail.example.com", false},
		{"invalid email no @", "testexample.com", true},
		{"invalid email no domain", "test@", true},
		{"invalid email no user", "@example.com", true},
		{"empty string skips validation", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Email("email", tt.value)

			if tt.wantError && !v.Errors().HasErrors() {
				t.Error("expected validation error")
			}
			if !tt.wantError && v.Errors().HasErrors() {
				t.Errorf("unexpected validation error: %v", v.Errors())
			}
		})
	}
}

func TestValidator_Password(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"valid password", "Password1", false},
		{"valid password complex", "MyP@ssw0rd!", false},
		{"too short", "Pass1", true},
		{"no uppercase", "password1", true},
		{"no lowercase", "PASSWORD1", true},
		{"no digit", "Password", true},
		{"empty string skips validation", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Password("password", tt.value)

			if tt.wantError && !v.Errors().HasErrors() {
				t.Error("expected validation error")
			}
			if !tt.wantError && v.Errors().HasErrors() {
				t.Errorf("unexpected validation error: %v", v.Errors())
			}
		})
	}
}

func TestValidator_UUID(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"valid UUID", "123e4567-e89b-12d3-a456-426614174000", false},
		{"invalid UUID", "not-a-uuid", true},
		{"empty string skips validation", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.UUID("id", tt.value)

			if tt.wantError && !v.Errors().HasErrors() {
				t.Error("expected validation error")
			}
			if !tt.wantError && v.Errors().HasErrors() {
				t.Errorf("unexpected validation error: %v", v.Errors())
			}
		})
	}
}

func TestValidator_InList(t *testing.T) {
	allowed := []string{"apple", "banana", "cherry"}

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"valid value", "apple", false},
		{"another valid value", "banana", false},
		{"invalid value", "orange", true},
		{"empty string skips validation", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.InList("fruit", tt.value, allowed)

			if tt.wantError && !v.Errors().HasErrors() {
				t.Error("expected validation error")
			}
			if !tt.wantError && v.Errors().HasErrors() {
				t.Errorf("unexpected validation error: %v", v.Errors())
			}
		})
	}
}

func TestValidator_PositiveNumber(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		wantError bool
	}{
		{"positive number", 10.5, false},
		{"zero", 0, true},
		{"negative number", -5.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.PositiveNumber("amount", tt.value)

			if tt.wantError && !v.Errors().HasErrors() {
				t.Error("expected validation error")
			}
			if !tt.wantError && v.Errors().HasErrors() {
				t.Errorf("unexpected validation error: %v", v.Errors())
			}
		})
	}
}

func TestValidator_ChainedValidations(t *testing.T) {
	v := New()
	v.Required("email", "test@example.com").
		Email("email", "test@example.com").
		Required("password", "Password123").
		Password("password", "Password123")

	if v.Errors().HasErrors() {
		t.Errorf("unexpected validation errors: %v", v.Errors())
	}
}

func TestValidator_MultipleErrors(t *testing.T) {
	v := New()
	v.Required("email", "")
	v.Required("password", "")
	v.Required("name", "")

	errors := v.Errors()
	if len(errors) != 3 {
		t.Errorf("expected 3 errors, got %d", len(errors))
	}
}
