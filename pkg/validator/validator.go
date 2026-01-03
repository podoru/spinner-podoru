package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

type Validator struct {
	validate *validator.Validate
}

func New() (*Validator, error) {
	v := validator.New()

	if err := v.RegisterValidation("slug", validateSlug); err != nil {
		return nil, err
	}

	return &Validator{validate: v}, nil
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func (v *Validator) ValidatePartial(i interface{}, fields ...string) error {
	return v.validate.StructPartial(i, fields...)
}

func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	return slugRegex.MatchString(slug)
}

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			errors[field] = formatErrorMessage(e)
		}
	}

	return errors
}

func formatErrorMessage(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + e.Param() + " characters"
	case "max":
		return field + " must be at most " + e.Param() + " characters"
	case "slug":
		return field + " must be a valid slug (lowercase letters, numbers, and hyphens)"
	case "url":
		return field + " must be a valid URL"
	case "fqdn":
		return field + " must be a valid domain name"
	case "cidr":
		return field + " must be a valid CIDR notation"
	case "ip":
		return field + " must be a valid IP address"
	case "oneof":
		return field + " must be one of: " + e.Param()
	default:
		return field + " is invalid"
	}
}

func GenerateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = regexp.MustCompile(`^-+|-+$`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	return slug
}
