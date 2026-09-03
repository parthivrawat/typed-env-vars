// Package env provides type-safe environment variable access for Go applications.
//
// This package offers automatic type conversion and validation for environment variables,
// making it easy to load configuration with proper type safety.
//
// Example:
//
//	import "github.com/parthivrawat/typed-env-vars/go"
//
//	databaseURL := env.URL("DATABASE_URL")
//	maxConnections := env.Int("MAX_CONNECTIONS", 10)
//	debug := env.Bool("DEBUG", false)
package env

import (
	"fmt"
	"net/url"
	"os"
	"sort"
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

// ListOpts holds optional configuration for List and ListE.
type ListOpts struct {
	Separator string
	Default   []string
}

// MapOpts holds optional configuration for Map and MapE.
type MapOpts struct {
	ItemSeparator     string
	KeyValueSeparator string
	Default           map[string]string
}

// DictOpts is an alias for MapOpts for use with Dict/DictE.
type DictOpts = MapOpts

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

// StringE gets a string environment variable, returning an error instead of panicking.
func StringE(key string, defaultValue ...string) (string, error) {
	var def *string
	if len(defaultValue) > 0 {
		def = &defaultValue[0]
	}

	return getRaw(key, def, def == nil)
}

// String gets a string environment variable
//
// Example:
//
//	apiKey := env.String("API_KEY")
//	appName := env.String("APP_NAME", "MyApp")
func String(key string, defaultValue ...string) string {
	v, err := StringE(key, defaultValue...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustString is like String but panics if the variable is not set.
func MustString(key string) string {
	return String(key)
}

// IntE gets an integer environment variable, returning an error instead of panicking.
func IntE(key string, defaultValue ...int) (int, error) {
	var def *string
	if len(defaultValue) > 0 {
		defStr := strconv.Itoa(defaultValue[0])
		def = &defStr
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		return 0, err
	}

	if rawValue == "" && def != nil {
		return defaultValue[0], nil
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, &EnvVarTypeError{
			EnvVarError: EnvVarError{Key: key, Message: "cannot be converted to int"},
			Value:       rawValue,
		}
	}

	return value, nil
}

// Int gets an integer environment variable
//
// Example:
//
//	port := env.Int("PORT")
//	maxRetries := env.Int("MAX_RETRIES", 3)
func Int(key string, defaultValue ...int) int {
	v, err := IntE(key, defaultValue...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustInt is like Int but panics if the variable is not set or invalid.
func MustInt(key string) int {
	return Int(key)
}

// FloatE gets a float64 environment variable, returning an error instead of panicking.
func FloatE(key string, defaultValue ...float64) (float64, error) {
	var def *string
	if len(defaultValue) > 0 {
		defStr := strconv.FormatFloat(defaultValue[0], 'f', -1, 64)
		def = &defStr
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		return 0, err
	}

	if rawValue == "" && def != nil {
		return defaultValue[0], nil
	}

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		return 0, &EnvVarTypeError{
			EnvVarError: EnvVarError{Key: key, Message: "cannot be converted to float"},
			Value:       rawValue,
		}
	}

	return value, nil
}

// Float gets a float64 environment variable
//
// Example:
//
//	timeout := env.Float("TIMEOUT")
//	rateLimit := env.Float("RATE_LIMIT", 1.5)
func Float(key string, defaultValue ...float64) float64 {
	v, err := FloatE(key, defaultValue...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustFloat is like Float but panics if the variable is not set or invalid.
func MustFloat(key string) float64 {
	return Float(key)
}

// BoolE gets a boolean environment variable, returning an error instead of panicking.
//
// Accepts: true/false, yes/no, 1/0, on/off (case-insensitive)
func BoolE(key string, defaultValue ...bool) (bool, error) {
	var def *string
	if len(defaultValue) > 0 {
		defStr := strconv.FormatBool(defaultValue[0])
		def = &defStr
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		return false, err
	}

	if rawValue == "" && def != nil {
		return defaultValue[0], nil
	}

	trueValues := map[string]bool{
		"true": true, "yes": true, "1": true, "on": true, "t": true, "y": true,
	}
	falseValues := map[string]bool{
		"false": true, "no": true, "0": true, "off": true, "f": true, "n": true,
	}

	normalized := strings.ToLower(strings.TrimSpace(rawValue))

	if trueValues[normalized] {
		return true, nil
	}
	if falseValues[normalized] {
		return false, nil
	}

	return false, &EnvVarTypeError{
		EnvVarError: EnvVarError{
			Key:     key,
			Message: "cannot be converted to bool. Valid values: true, false, yes, no, 1, 0, on, off",
		},
		Value: rawValue,
	}
}

// Bool gets a boolean environment variable
//
// Example:
//
//	debug := env.Bool("DEBUG")
//	enableCache := env.Bool("ENABLE_CACHE", true)
func Bool(key string, defaultValue ...bool) bool {
	v, err := BoolE(key, defaultValue...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustBool is like Bool but panics if the variable is not set or invalid.
func MustBool(key string) bool {
	return Bool(key)
}

// ListE gets a list environment variable, returning an error instead of panicking.
//
// Example:
//
//	hosts, err := env.ListE("ALLOWED_HOSTS")
//	tags, err := env.ListE("TAGS", &env.ListOpts{Separator: ":"})
//	defaultHosts, err := env.ListE("HOSTS", &env.ListOpts{Default: []string{"localhost"}})
func ListE(key string, opts ...*ListOpts) ([]string, error) {
	cfg := &ListOpts{Separator: ","}
	if len(opts) > 0 && opts[0] != nil {
		cfg = opts[0]
	}
	separator := cfg.Separator
	if separator == "" {
		separator = ","
	}
	defaultValue := cfg.Default

	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue != nil {
			out := make([]string, len(defaultValue))
			copy(out, defaultValue)
			return out, nil
		}
		return nil, &EnvVarNotFoundError{
			EnvVarError: EnvVarError{Key: key},
		}
	}

	if value == "" {
		if defaultValue != nil {
			out := make([]string, len(defaultValue))
			copy(out, defaultValue)
			return out, nil
		}
		return []string{}, nil
	}

	parts := strings.Split(value, separator)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result, nil
}

// List gets a list environment variable
//
// Example:
//
//	hosts := env.List("ALLOWED_HOSTS")
//	tags := env.List("TAGS", &env.ListOpts{Separator: ":"})
//	defaultHosts := env.List("HOSTS", &env.ListOpts{Default: []string{"localhost"}})
func List(key string, opts ...*ListOpts) []string {
	v, err := ListE(key, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustList is like List but panics if the variable is not set or invalid.
func MustList(key string) []string {
	return List(key)
}

// MapE gets a map environment variable, returning an error instead of panicking.
//
// Format: KEY1=VALUE1,KEY2=VALUE2
//
// Example:
//
//	flags, err := env.MapE("FEATURE_FLAGS")
//	settings, err := env.MapE("SETTINGS", &env.MapOpts{ItemSeparator: ";", KeyValueSeparator: ":"})
func MapE(key string, opts ...*MapOpts) (map[string]string, error) {
	cfg := &MapOpts{ItemSeparator: ",", KeyValueSeparator: "="}
	if len(opts) > 0 && opts[0] != nil {
		cfg = opts[0]
	}
	itemSeparator := cfg.ItemSeparator
	if itemSeparator == "" {
		itemSeparator = ","
	}
	keyValueSeparator := cfg.KeyValueSeparator
	if keyValueSeparator == "" {
		keyValueSeparator = "="
	}
	defaultValue := cfg.Default

	copyDefault := func() map[string]string {
		out := make(map[string]string, len(defaultValue))
		for k, v := range defaultValue {
			out[k] = v
		}
		return out
	}

	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue != nil {
			return copyDefault(), nil
		}
		return nil, &EnvVarNotFoundError{
			EnvVarError: EnvVarError{Key: key},
		}
	}

	if value == "" {
		if defaultValue != nil {
			return copyDefault(), nil
		}
		return map[string]string{}, nil
	}

	result := make(map[string]string)
	items := strings.Split(value, itemSeparator)

	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}

		parts := strings.SplitN(trimmed, keyValueSeparator, 2)
		if len(parts) != 2 {
			return nil, &EnvVarTypeError{
				EnvVarError: EnvVarError{
					Key:     key,
					Message: fmt.Sprintf("cannot be parsed as map. Expected format: KEY1%sVALUE1%sKEY2%sVALUE2", keyValueSeparator, itemSeparator, keyValueSeparator),
				},
				Value: value,
			}
		}

		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return result, nil
}

// Map gets a map environment variable
//
// Format: KEY1=VALUE1,KEY2=VALUE2
//
// Example:
//
//	flags := env.Map("FEATURE_FLAGS")
//	settings := env.Map("SETTINGS", &env.MapOpts{ItemSeparator: ";", KeyValueSeparator: ":"})
func Map(key string, opts ...*MapOpts) map[string]string {
	v, err := MapE(key, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustMap is like Map but panics if the variable is not set or invalid.
func MustMap(key string) map[string]string {
	return Map(key)
}

func mapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// EnumE gets an enum environment variable with validation
//
// The values map maps accepted string values to their typed enum values.
// Matching is case-insensitive.
//
// Example:
//
//	type Environment string
//	const (
//	    EnvDevelopment Environment = "development"
//	    EnvStaging     Environment = "staging"
//	)
//	envMap := map[string]Environment{
//	    "development": EnvDevelopment,
//	    "staging":     EnvStaging,
//	}
//	environment, err := env.EnumE("ENVIRONMENT", envMap, EnvDevelopment)
func EnumE[T any](key string, values map[string]T, defaultValue ...T) (T, error) {
	var zero T
	rawValue, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return zero, &EnvVarNotFoundError{
			EnvVarError: EnvVarError{Key: key},
		}
	}

	normalized := strings.ToLower(strings.TrimSpace(rawValue))
	for k, v := range values {
		if strings.ToLower(strings.TrimSpace(k)) == normalized {
			return v, nil
		}
	}

	return zero, &EnvVarTypeError{
		EnvVarError: EnvVarError{
			Key:     key,
			Message: fmt.Sprintf("cannot be matched to a valid enum value. Valid values: %s", strings.Join(mapKeys(values), ", ")),
		},
		Value: rawValue,
	}
}

// Enum gets an enum environment variable
//
// Example:
//
//	environment := env.Enum("ENVIRONMENT", map[string]string{
//	    "development": "development",
//	    "staging":     "staging",
//	}, "development")
func Enum[T any](key string, values map[string]T, defaultValue ...T) T {
	v, err := EnumE(key, values, defaultValue...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustEnum is like Enum but panics if the variable is not set or invalid
func MustEnum[T any](key string, values map[string]T) T {
	return Enum(key, values)
}

// URLE gets a URL environment variable with validation, returning an error instead of panicking.
//
// Example:
//
//	dbURL, err := env.URLE("DATABASE_URL")
//	apiEndpoint, err := env.URLE("API_ENDPOINT", "https://api.example.com")
func URLE(key string, defaultValue ...string) (string, error) {
	var def *string
	if len(defaultValue) > 0 {
		def = &defaultValue[0]
	}

	value, err := getRaw(key, def, def == nil)
	if err != nil {
		return "", err
	}

	if value == "" {
		return value, nil
	}

	// Validate URL
	parsedURL, err := url.Parse(value)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", &EnvVarTypeError{
			EnvVarError: EnvVarError{
				Key:     key,
				Message: "is not a valid URL",
			},
			Value: value,
		}
	}

	return value, nil
}

// URL gets a URL environment variable with validation
//
// Example:
//
//	dbURL := env.URL("DATABASE_URL")
//	apiEndpoint := env.URL("API_ENDPOINT", "https://api.example.com")
func URL(key string, defaultValue ...string) string {
	v, err := URLE(key, defaultValue...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustURL is like URL but panics if the variable is not set or invalid.
func MustURL(key string) string {
	return URL(key)
}

// CustomE gets an environment variable with a custom converter function, returning an error instead of panicking.
//
// Example:
//
//	releaseDate, err := env.CustomE("RELEASE_DATE", func(s string) (time.Time, error) {
//	    return time.Parse("2006-01-02", s)
//	})
func CustomE[T any](key string, converter func(string) (T, error), defaultValue ...T) (T, error) {
	var zero T
	var def *string
	if len(defaultValue) > 0 {
		def = new(string)
		*def = fmt.Sprintf("%v", defaultValue[0])
	}

	rawValue, err := getRaw(key, def, def == nil)
	if err != nil {
		return zero, err
	}

	if rawValue == "" && def != nil {
		return defaultValue[0], nil
	}

	value, err := converter(rawValue)
	if err != nil {
		return zero, &EnvVarTypeError{
			EnvVarError: EnvVarError{
				Key:     key,
				Message: fmt.Sprintf("conversion failed: %v", err),
			},
			Value: rawValue,
		}
	}

	return value, nil
}

// Custom gets an environment variable with a custom converter function
//
// Example:
//
//	releaseDate := env.Custom("RELEASE_DATE", func(s string) (time.Time, error) {
//	    return time.Parse("2006-01-02", s)
//	})
func Custom[T any](key string, converter func(string) (T, error), defaultValue ...T) T {
	v, err := CustomE(key, converter, defaultValue...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustCustom is like Custom but panics if the variable is not set or conversion fails.
func MustCustom[T any](key string, converter func(string) (T, error)) T {
	return Custom(key, converter)
}

// Str is an alias for String for a common cross-language vocabulary.
func Str(key string, defaultValue ...string) string {
	return String(key, defaultValue...)
}

// StrE is an alias for StringE.
func StrE(key string, defaultValue ...string) (string, error) {
	return StringE(key, defaultValue...)
}

// Dict is an alias for Map for a common cross-language vocabulary.
func Dict(key string, opts ...*DictOpts) map[string]string {
	return Map(key, opts...)
}

// DictE is an alias for MapE.
func DictE(key string, opts ...*DictOpts) (map[string]string, error) {
	return MapE(key, opts...)
}
