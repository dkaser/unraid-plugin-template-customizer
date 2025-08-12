package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

// Validator for plugin/package name
var nameRegex = regexp.MustCompile(`^[a-z0-9_-]+$`)

func nameValidator(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid input")
	}
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("value cannot be empty")
	}
	if !nameRegex.MatchString(str) {
		return fmt.Errorf("must be lowercase, no spaces, no special characters except hyphen or underscore (a-z, 0-9, - or _)")
	}
	return nil
}

// Validator for URL
func urlValidator(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid input")
	}
	if strings.TrimSpace(str) == "" {
		return nil // allow blank (will default later)
	}
	parsed, err := url.ParseRequestURI(str)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("must be a valid URL (http/https)")
	}
	return nil
}

// Validator for plugin description (max 6 lines after wrapping at 70 chars)
func descriptionValidator(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid input")
	}
	wrapped := wordwrap.WrapString(str, 70)
	descLines := strings.Split(wrapped, "\n")
	if len(descLines) > 6 {
		return fmt.Errorf("description is too long (%d lines after wrapping, must be 6 or less)", len(descLines))
	}
	return nil
}
