package validator_test

import (
	"testing"

	"github.com/podoru/podoru/pkg/validator"
)

type TestStruct struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required,min=2,max=50"`
	Slug     string `validate:"required,slug"`
	Password string `validate:"required,min=8"`
}

func TestValidator_Valid(t *testing.T) {
	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	input := TestStruct{
		Email:    "test@example.com",
		Name:     "John Doe",
		Slug:     "my-test-slug",
		Password: "password123",
	}

	err = v.Validate(&input)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidator_InvalidEmail(t *testing.T) {
	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	input := TestStruct{
		Email:    "invalid-email",
		Name:     "John Doe",
		Slug:     "my-slug",
		Password: "password123",
	}

	err = v.Validate(&input)
	if err == nil {
		t.Error("expected validation error for invalid email")
	}
}

func TestValidator_RequiredField(t *testing.T) {
	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	input := TestStruct{
		Email:    "test@example.com",
		Name:     "",
		Slug:     "my-slug",
		Password: "password123",
	}

	err = v.Validate(&input)
	if err == nil {
		t.Error("expected validation error for missing required field")
	}
}

func TestValidator_InvalidSlug(t *testing.T) {
	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	testCases := []struct {
		slug   string
		valid  bool
		reason string
	}{
		{"my-slug", true, "valid slug"},
		{"my-long-slug-123", true, "slug with numbers"},
		{"a", true, "single character"},
		{"My-Slug", false, "uppercase letters"},
		{"my_slug", false, "underscore not allowed"},
		{"my slug", false, "space not allowed"},
		{"-my-slug", false, "starts with hyphen"},
		{"my-slug-", false, "ends with hyphen"},
		{"my--slug", false, "double hyphen"},
	}

	for _, tc := range testCases {
		input := TestStruct{
			Email:    "test@example.com",
			Name:     "John Doe",
			Slug:     tc.slug,
			Password: "password123",
		}

		err := v.Validate(&input)
		if tc.valid && err != nil {
			t.Errorf("slug '%s' should be valid (%s), got error: %v", tc.slug, tc.reason, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("slug '%s' should be invalid (%s)", tc.slug, tc.reason)
		}
	}
}

func TestValidator_MinLength(t *testing.T) {
	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	input := TestStruct{
		Email:    "test@example.com",
		Name:     "A",
		Slug:     "my-slug",
		Password: "password123",
	}

	err = v.Validate(&input)
	if err == nil {
		t.Error("expected validation error for name too short")
	}
}

func TestValidator_PasswordTooShort(t *testing.T) {
	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	input := TestStruct{
		Email:    "test@example.com",
		Name:     "John Doe",
		Slug:     "my-slug",
		Password: "short",
	}

	err = v.Validate(&input)
	if err == nil {
		t.Error("expected validation error for password too short")
	}
}

func TestFormatValidationErrors(t *testing.T) {
	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	input := TestStruct{
		Email:    "invalid",
		Name:     "",
		Slug:     "INVALID",
		Password: "short",
	}

	err = v.Validate(&input)
	if err == nil {
		t.Fatal("expected validation errors")
	}

	errors := validator.FormatValidationErrors(err)

	if len(errors) == 0 {
		t.Error("expected formatted errors")
	}

	if _, ok := errors["email"]; !ok {
		t.Error("expected email error")
	}

	if _, ok := errors["name"]; !ok {
		t.Error("expected name error")
	}
}

func TestGenerateSlug(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"My Project Name", "my-project-name"},
		{"Test 123", "test-123"},
		{"  Spaces  ", "spaces"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"Special@#$Characters", "special-characters"},
		{"Already-Slug", "already-slug"},
		{"MixedCase123", "mixedcase123"},
	}

	for _, tc := range testCases {
		result := validator.GenerateSlug(tc.input)
		if result != tc.expected {
			t.Errorf("GenerateSlug(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

type PartialStruct struct {
	Email string `validate:"required,email"`
	Name  string `validate:"required,min=2"`
}

func TestValidatePartial(t *testing.T) {
	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	input := PartialStruct{
		Email: "test@example.com",
		Name:  "",
	}

	err = v.ValidatePartial(&input, "Email")
	if err != nil {
		t.Errorf("partial validation should pass for Email: %v", err)
	}
}
