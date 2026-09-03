package env

import (
	"os"
	"testing"
)

func TestString(t *testing.T) {
	t.Run("returns value when present", func(t *testing.T) {
		os.Setenv("TEST_VAR", "hello")
		defer os.Unsetenv("TEST_VAR")

		result := String("TEST_VAR")
		if result != "hello" {
			t.Errorf("expected 'hello', got '%s'", result)
		}
	})

	t.Run("panics when required variable is missing", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for missing required variable")
			}
		}()
		String("MISSING_VAR")
	})

	t.Run("returns default when variable is missing", func(t *testing.T) {
		result := String("MISSING_VAR", "default")
		if result != "default" {
			t.Errorf("expected 'default', got '%s'", result)
		}
	})
}

func TestInt(t *testing.T) {
	t.Run("parses valid integer", func(t *testing.T) {
		os.Setenv("TEST_VAR", "42")
		defer os.Unsetenv("TEST_VAR")

		result := Int("TEST_VAR")
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("parses negative integer", func(t *testing.T) {
		os.Setenv("TEST_VAR", "-10")
		defer os.Unsetenv("TEST_VAR")

		result := Int("TEST_VAR")
		if result != -10 {
			t.Errorf("expected -10, got %d", result)
		}
	})

	t.Run("panics on invalid integer", func(t *testing.T) {
		os.Setenv("TEST_VAR", "not_a_number")
		defer os.Unsetenv("TEST_VAR")

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid integer")
			}
		}()
		Int("TEST_VAR")
	})

	t.Run("returns default value", func(t *testing.T) {
		result := Int("MISSING_VAR", 100)
		if result != 100 {
			t.Errorf("expected 100, got %d", result)
		}
	})
}

func TestFloat(t *testing.T) {
	t.Run("parses valid float", func(t *testing.T) {
		os.Setenv("TEST_VAR", "3.14")
		defer os.Unsetenv("TEST_VAR")

		result := Float("TEST_VAR")
		if result != 3.14 {
			t.Errorf("expected 3.14, got %f", result)
		}
	})

	t.Run("parses scientific notation", func(t *testing.T) {
		os.Setenv("TEST_VAR", "1.5e-10")
		defer os.Unsetenv("TEST_VAR")

		result := Float("TEST_VAR")
		if result != 1.5e-10 {
			t.Errorf("expected 1.5e-10, got %e", result)
		}
	})

	t.Run("panics on invalid float", func(t *testing.T) {
		os.Setenv("TEST_VAR", "not_a_float")
		defer os.Unsetenv("TEST_VAR")

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid float")
			}
		}()
		Float("TEST_VAR")
	})

	t.Run("returns default value", func(t *testing.T) {
		result := Float("MISSING_VAR", 2.5)
		if result != 2.5 {
			t.Errorf("expected 2.5, got %f", result)
		}
	})
}

func TestBool(t *testing.T) {
	trueValues := []string{"true", "True", "TRUE", "yes", "YES", "1", "on", "ON", "t", "y"}
	falseValues := []string{"false", "False", "FALSE", "no", "NO", "0", "off", "OFF", "f", "n"}

	t.Run("parses true values", func(t *testing.T) {
		for _, value := range trueValues {
			os.Setenv("TEST_VAR", value)
			result := Bool("TEST_VAR")
			if !result {
				t.Errorf("expected true for value '%s', got false", value)
			}
			os.Unsetenv("TEST_VAR")
		}
	})

	t.Run("parses false values", func(t *testing.T) {
		for _, value := range falseValues {
			os.Setenv("TEST_VAR", value)
			result := Bool("TEST_VAR")
			if result {
				t.Errorf("expected false for value '%s', got true", value)
			}
			os.Unsetenv("TEST_VAR")
		}
	})

	t.Run("panics on invalid boolean", func(t *testing.T) {
		os.Setenv("TEST_VAR", "maybe")
		defer os.Unsetenv("TEST_VAR")

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid boolean")
			}
		}()
		Bool("TEST_VAR")
	})

	t.Run("returns default value", func(t *testing.T) {
		result := Bool("MISSING_VAR", true)
		if !result {
			t.Error("expected true, got false")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("parses comma-separated list", func(t *testing.T) {
		os.Setenv("TEST_VAR", "a,b,c")
		defer os.Unsetenv("TEST_VAR")

		result := List("TEST_VAR")
		expected := []string{"a", "b", "c"}

		if len(result) != len(expected) {
			t.Errorf("expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("expected '%s' at index %d, got '%s'", v, i, result[i])
			}
		}
	})

	t.Run("trims whitespace", func(t *testing.T) {
		os.Setenv("TEST_VAR", "a, b , c")
		defer os.Unsetenv("TEST_VAR")

		result := List("TEST_VAR")
		expected := []string{"a", "b", "c"}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("expected '%s' at index %d, got '%s'", v, i, result[i])
			}
		}
	})

	t.Run("uses custom separator", func(t *testing.T) {
		os.Setenv("TEST_VAR", "a:b:c")
		defer os.Unsetenv("TEST_VAR")

		result := List("TEST_VAR", &ListOpts{Separator: ":"})
		expected := []string{"a", "b", "c"}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("expected '%s' at index %d, got '%s'", v, i, result[i])
			}
		}
	})

	t.Run("returns empty list for empty string", func(t *testing.T) {
		os.Setenv("TEST_VAR", "")
		defer os.Unsetenv("TEST_VAR")

		result := List("TEST_VAR")
		if len(result) != 0 {
			t.Errorf("expected empty list, got %v", result)
		}
	})

	t.Run("returns default value", func(t *testing.T) {
		result := List("MISSING_VAR", &ListOpts{Separator: ",", Default: []string{"x", "y"}})
		expected := []string{"x", "y"}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("expected '%s' at index %d, got '%s'", v, i, result[i])
			}
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("parses map", func(t *testing.T) {
		os.Setenv("TEST_VAR", "key1=value1,key2=value2")
		defer os.Unsetenv("TEST_VAR")

		result := Map("TEST_VAR")
		expected := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		if len(result) != len(expected) {
			t.Errorf("expected length %d, got %d", len(expected), len(result))
		}

		for k, v := range expected {
			if result[k] != v {
				t.Errorf("expected '%s' for key '%s', got '%s'", v, k, result[k])
			}
		}
	})

	t.Run("trims whitespace", func(t *testing.T) {
		os.Setenv("TEST_VAR", "key1 = value1 , key2 = value2")
		defer os.Unsetenv("TEST_VAR")

		result := Map("TEST_VAR")
		expected := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		for k, v := range expected {
			if result[k] != v {
				t.Errorf("expected '%s' for key '%s', got '%s'", v, k, result[k])
			}
		}
	})

	t.Run("uses custom separators", func(t *testing.T) {
		os.Setenv("TEST_VAR", "key1:value1;key2:value2")
		defer os.Unsetenv("TEST_VAR")

		result := Map("TEST_VAR", &MapOpts{ItemSeparator: ";", KeyValueSeparator: ":"})
		expected := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		for k, v := range expected {
			if result[k] != v {
				t.Errorf("expected '%s' for key '%s', got '%s'", v, k, result[k])
			}
		}
	})

	t.Run("returns empty map for empty string", func(t *testing.T) {
		os.Setenv("TEST_VAR", "")
		defer os.Unsetenv("TEST_VAR")

		result := Map("TEST_VAR")
		if len(result) != 0 {
			t.Errorf("expected empty map, got %v", result)
		}
	})

	t.Run("panics on invalid format", func(t *testing.T) {
		os.Setenv("TEST_VAR", "invalid_format")
		defer os.Unsetenv("TEST_VAR")

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid map format")
			}
		}()
		Map("TEST_VAR")
	})
}

func TestURL(t *testing.T) {
	t.Run("accepts http URL", func(t *testing.T) {
		os.Setenv("TEST_VAR", "http://example.com")
		defer os.Unsetenv("TEST_VAR")

		result := URL("TEST_VAR")
		if result != "http://example.com" {
			t.Errorf("expected 'http://example.com', got '%s'", result)
		}
	})

	t.Run("accepts https URL", func(t *testing.T) {
		os.Setenv("TEST_VAR", "https://example.com/path")
		defer os.Unsetenv("TEST_VAR")

		result := URL("TEST_VAR")
		if result != "https://example.com/path" {
			t.Errorf("expected 'https://example.com/path', got '%s'", result)
		}
	})

	t.Run("panics on invalid URL", func(t *testing.T) {
		os.Setenv("TEST_VAR", "not-a-url")
		defer os.Unsetenv("TEST_VAR")

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid URL")
			}
		}()
		URL("TEST_VAR")
	})

	t.Run("returns default value", func(t *testing.T) {
		result := URL("MISSING_VAR", "https://default.com")
		if result != "https://default.com" {
			t.Errorf("expected 'https://default.com', got '%s'", result)
		}
	})
}

func TestStringE(t *testing.T) {
	t.Run("returns value and nil error when present", func(t *testing.T) {
		os.Setenv("TEST_VAR", "hello")
		defer os.Unsetenv("TEST_VAR")

		result, err := StringE("TEST_VAR")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "hello" {
			t.Errorf("expected 'hello', got '%s'", result)
		}
	})

	t.Run("returns error when required variable is missing", func(t *testing.T) {
		_, err := StringE("MISSING_VAR")
		if _, ok := err.(*EnvVarNotFoundError); !ok {
			t.Errorf("expected EnvVarNotFoundError, got %T", err)
		}
	})
}

func TestEnum(t *testing.T) {
	type Environment string
	const (
		EnvDevelopment Environment = "development"
		EnvStaging     Environment = "staging"
	)

	values := map[string]Environment{
		"development": EnvDevelopment,
		"staging":     EnvStaging,
	}

	t.Run("parses valid enum value case-insensitively", func(t *testing.T) {
		os.Setenv("TEST_VAR", "STAGING")
		defer os.Unsetenv("TEST_VAR")

		result := Enum("TEST_VAR", values)
		if result != EnvStaging {
			t.Errorf("expected staging, got %v", result)
		}
	})

	t.Run("returns default when variable is missing", func(t *testing.T) {
		result := Enum("MISSING_VAR", values, EnvDevelopment)
		if result != EnvDevelopment {
			t.Errorf("expected development, got %v", result)
		}
	})
}

func TestAliases(t *testing.T) {
	t.Run("Str is an alias for String", func(t *testing.T) {
		os.Setenv("TEST_VAR", "hello")
		defer os.Unsetenv("TEST_VAR")

		if Str("TEST_VAR") != "hello" {
			t.Errorf("Str alias did not return expected value")
		}
	})

	t.Run("Dict is an alias for Map", func(t *testing.T) {
		os.Setenv("TEST_VAR", "a=1,b=2")
		defer os.Unsetenv("TEST_VAR")

		result := Dict("TEST_VAR")
		if result["a"] != "1" || result["b"] != "2" {
			t.Errorf("Dict alias did not return expected value")
		}
	})

	t.Run("list default is returned without serialization", func(t *testing.T) {
		result := List("MISSING_VAR", &ListOpts{Default: []string{"a,b"}})
		if len(result) != 1 || result[0] != "a,b" {
			t.Errorf("expected ['a,b'], got %v", result)
		}
	})
}

func TestRealWorldScenarios(t *testing.T) {
	t.Run("database configuration", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgresql://localhost:5432/mydb")
		os.Setenv("DB_POOL_SIZE", "20")
		os.Setenv("DB_TIMEOUT", "30.5")
		os.Setenv("DB_SSL", "true")
		defer func() {
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("DB_POOL_SIZE")
			os.Unsetenv("DB_TIMEOUT")
			os.Unsetenv("DB_SSL")
		}()

		dbURL := URL("DATABASE_URL")
		poolSize := Int("DB_POOL_SIZE")
		timeout := Float("DB_TIMEOUT")
		ssl := Bool("DB_SSL")

		if dbURL != "postgresql://localhost:5432/mydb" {
			t.Errorf("unexpected DATABASE_URL: %s", dbURL)
		}
		if poolSize != 20 {
			t.Errorf("unexpected DB_POOL_SIZE: %d", poolSize)
		}
		if timeout != 30.5 {
			t.Errorf("unexpected DB_TIMEOUT: %f", timeout)
		}
		if !ssl {
			t.Error("expected DB_SSL to be true")
		}
	})

	t.Run("application configuration", func(t *testing.T) {
		os.Setenv("APP_NAME", "MyApp")
		os.Setenv("DEBUG", "false")
		os.Setenv("ALLOWED_HOSTS", "localhost,127.0.0.1,example.com")
		defer func() {
			os.Unsetenv("APP_NAME")
			os.Unsetenv("DEBUG")
			os.Unsetenv("ALLOWED_HOSTS")
		}()

		appName := String("APP_NAME")
		debug := Bool("DEBUG")
		allowedHosts := List("ALLOWED_HOSTS")

		if appName != "MyApp" {
			t.Errorf("unexpected APP_NAME: %s", appName)
		}
		if debug {
			t.Error("expected DEBUG to be false")
		}
		if len(allowedHosts) != 3 {
			t.Errorf("expected 3 allowed hosts, got %d", len(allowedHosts))
		}
	})
}
