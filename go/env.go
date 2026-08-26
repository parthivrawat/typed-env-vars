// Package env provides type-safe environment variable access for Go applications.
//
// This package offers automatic type conversion and validation for environment variables,
// making it easy to load configuration with proper type safety.
//
// Example:
//
//	import "github.com/parthivrawat/typed-env-vars"
//
//	databaseURL := env.URL("DATABASE_URL")
//	maxConnections := env.Int("MAX_CONNECTIONS", 10)
//	debug := env.Bool("DEBUG", false)
package env

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Error types for environment variable operations
type (
	// EnvVarError is the base error type for environment variable errors
	EnvVarError struct {
		Key     string
		Message string
	}

	// EnvVarNotFoundError is returned when a required environment variable is not set
	EnvVarNotFoundError struct {
		EnvVarError
	}

	// EnvVarTypeError is returned when an environment variable cannot be converted to the expected type
	EnvVarTypeError struct {
		EnvVarError
		Value string
	}
)

func (e *EnvVarError) Error() string {
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

func (e *EnvVarNotFoundError) Error() string {
	return fmt.Sprintf("required environment variable '%s' is not set", e.Key)
}

func (e *EnvVarTypeError) Error() string {
	return fmt.Sprintf("environment variable '%s' has value '%s' which %s", e.Key, e.Value, e.Message)
}

// getRaw gets the raw environment variable value
func getRaw(key string, defaultValue *string, required bool) (string, error) {
	value, exists := os.LookupEnv(key)

	if !exists {
		if required && defaultValue == nil {
			return "", &EnvVarNotFoundError{
				EnvVarError: EnvVarError{Key: key},
			}
		}
		if defaultValue != nil {
			return *defaultValue, nil
		}
		return "", nil
	}

	return value, nil
}

// String gets a string environment variable
//
// Example:
//
//	apiKey := env.String("API_KEY")
//	appName := env.String("APP_NAME", "MyApp")
func String(key string, defaultValue ...string) string {
	var def *string
	if len(defaultValue) > 0 {
		def = &defaultValue[0]
	}

	value, err := getRaw(key, def, def == nil)
	if err != nil {
		panic(err)
	}

	return value
}

// Int gets an integer environment variable
//
// Example:
//
//	port := env.Int("PORT")
//	maxRetries := env.Int("MAX_RETRIES", 3)
func Int(key string, defaultValue ...int) int {
	var def *string
	if len(defaultValue) > 0 {
		defStr := strconv.Itoa(defaultValue[0])
		def = &defStr
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		panic(err)
	}

	if rawValue == "" && def != nil {
		return defaultValue[0]
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		panic(&EnvVarTypeError{
			EnvVarError: EnvVarError{Key: key, Message: "cannot be converted to int"},
			Value:       rawValue,
		})
	}

	return value
}

// Float gets a float64 environment variable
//
// Example:
//
//	timeout := env.Float("TIMEOUT")
//	rateLimit := env.Float("RATE_LIMIT", 1.5)
func Float(key string, defaultValue ...float64) float64 {
	var def *string
	if len(defaultValue) > 0 {
		defStr := strconv.FormatFloat(defaultValue[0], 'f', -1, 64)
		def = &defStr
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		panic(err)
	}

	if rawValue == "" && def != nil {
		return defaultValue[0]
	}

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		panic(&EnvVarTypeError{
			EnvVarError: EnvVarError{Key: key, Message: "cannot be converted to float"},
			Value:       rawValue,
		})
	}

	return value
}

// Bool gets a boolean environment variable
//
// Accepts: true/false, yes/no, 1/0, on/off (case-insensitive)
//
// Example:
//
//	debug := env.Bool("DEBUG")
//	enableCache := env.Bool("ENABLE_CACHE", true)
func Bool(key string, defaultValue ...bool) bool {
	var def *string
	if len(defaultValue) > 0 {
		defStr := strconv.FormatBool(defaultValue[0])
		def = &defStr
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		panic(err)
	}

	if rawValue == "" && def != nil {
		return defaultValue[0]
	}

	trueValues := map[string]bool{
		"true": true, "yes": true, "1": true, "on": true, "t": true, "y": true,
	}
	falseValues := map[string]bool{
		"false": true, "no": true, "0": true, "off": true, "f": true, "n": true,
	}

	normalized := strings.ToLower(strings.TrimSpace(rawValue))

	if trueValues[normalized] {
		return true
	}
	if falseValues[normalized] {
		return false
	}

	panic(&EnvVarTypeError{
		EnvVarError: EnvVarError{
			Key:     key,
			Message: "cannot be converted to bool. Valid values: true, false, yes, no, 1, 0, on, off",
		},
		Value: rawValue,
	})
}

// List gets a list environment variable
//
// Example:
//
//	hosts := env.List("ALLOWED_HOSTS")
//	tags := env.List("TAGS", ",")
//	defaultHosts := env.List("HOSTS", ",", []string{"localhost"})
func List(key string, args ...interface{}) []string {
	separator := ","
	var defaultValue []string

	// Parse arguments
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			separator = v
		case []string:
			defaultValue = v
		}
	}

	var def *string
	if defaultValue != nil {
		defStr := strings.Join(defaultValue, separator)
		def = &defStr
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		panic(err)
	}

	if rawValue == "" {
		if defaultValue != nil {
			return defaultValue
		}
		return []string{}
	}

	parts := strings.Split(rawValue, separator)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// Map gets a map environment variable
//
// Format: KEY1=VALUE1,KEY2=VALUE2
//
// Example:
//
//	flags := env.Map("FEATURE_FLAGS")
//	settings := env.Map("SETTINGS", ";", ":")
func Map(key string, args ...interface{}) map[string]string {
	itemSeparator := ","
	keyValueSeparator := "="
	var defaultValue map[string]string

	// Parse arguments
	for i, arg := range args {
		switch v := arg.(type) {
		case string:
			if i == 0 {
				itemSeparator = v
			} else if i == 1 {
				keyValueSeparator = v
			}
		case map[string]string:
			defaultValue = v
		}
	}

	var def *string
	if defaultValue != nil {
		pairs := make([]string, 0, len(defaultValue))
		for k, v := range defaultValue {
			pairs = append(pairs, k+keyValueSeparator+v)
		}
		defStr := strings.Join(pairs, itemSeparator)
		def = &defStr
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		panic(err)
	}

	if rawValue == "" {
		if defaultValue != nil {
			return defaultValue
		}
		return map[string]string{}
	}

	result := make(map[string]string)
	items := strings.Split(rawValue, itemSeparator)

	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}

		parts := strings.SplitN(trimmed, keyValueSeparator, 2)
		if len(parts) != 2 {
			panic(&EnvVarTypeError{
				EnvVarError: EnvVarError{
					Key:     key,
					Message: fmt.Sprintf("cannot be parsed as map. Expected format: KEY1%sVALUE1%sKEY2%sVALUE2", keyValueSeparator, itemSeparator, keyValueSeparator),
				},
				Value: rawValue,
			})
		}

		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return result
}

// URL gets a URL environment variable with validation
//
// Example:
//
//	dbURL := env.URL("DATABASE_URL")
//	apiEndpoint := env.URL("API_ENDPOINT", "https://api.example.com")
func URL(key string, defaultValue ...string) string {
	var def *string
	if len(defaultValue) > 0 {
		def = &defaultValue[0]
	}

	value, err := getRaw(key, def, def == nil)
	if err != nil {
		panic(err)
	}

	if value == "" {
		return value
	}

	// Validate URL
	parsedURL, err := url.Parse(value)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		panic(&EnvVarTypeError{
			EnvVarError: EnvVarError{
				Key:     key,
				Message: "is not a valid URL",
			},
			Value: value,
		})
	}

	return value
}

// Custom gets an environment variable with a custom converter function
//
// Example:
//
//	releaseDate := env.Custom("RELEASE_DATE", func(s string) (interface{}, error) {
//	    return time.Parse("2006-01-02", s)
//	})
func Custom(key string, converter func(string) (interface{}, error), defaultValue ...interface{}) interface{} {
	var def *string
	if len(defaultValue) > 0 {
		def = new(string)
		*def = fmt.Sprintf("%v", defaultValue[0])
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		panic(err)
	}

	if rawValue == "" && def != nil {
		return defaultValue[0]
	}

	value, err := converter(rawValue)
	if err != nil {
		panic(&EnvVarTypeError{
			EnvVarError: EnvVarError{
				Key:     key,
				Message: fmt.Sprintf("conversion failed: %v", err),
			},
			Value: rawValue,
		})
	}

	return value
}

// MustString is like String but panics if the variable is not set
func MustString(key string) string {
	return String(key)
}

// MustInt is like Int but panics if the variable is not set
func MustInt(key string) int {
	return Int(key)
}

// MustFloat is like Float but panics if the variable is not set
func MustFloat(key string) float64 {
	return Float(key)
}

// MustBool is like Bool but panics if the variable is not set
func MustBool(key string) bool {
	return Bool(key)
}

// MustURL is like URL but panics if the variable is not set
func MustURL(key string) string {
	return URL(key)
}
