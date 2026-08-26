# Typed Environment Variables (Go)

A type-safe environment variable library for Go that provides automatic type conversion and validation.

## Features

- ✅ **Type Safety**: Automatic type conversion with validation
- ✅ **Rich Types**: Support for string, int, float64, bool, list, map, URL
- ✅ **Default Values**: Optional default values for missing variables
- ✅ **Clear Errors**: Descriptive error messages for debugging
- ✅ **Zero Dependencies**: No external dependencies required
- ✅ **Production Ready**: Comprehensive test coverage

## Installation

```bash
go get github.com/parthivrawat/typed-env-vars/go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/parthivrawat/typed-env-vars/go"
)

func main() {
    databaseURL := env.URL("DATABASE_URL")
    maxConnections := env.Int("MAX_CONNECTIONS", 10)
    debug := env.Bool("DEBUG", false)
    allowedHosts := env.List("ALLOWED_HOSTS", ",", []string{"localhost"})
    
    fmt.Printf("Database: %s\n", databaseURL)
    fmt.Printf("Max Connections: %d\n", maxConnections)
    fmt.Printf("Debug: %v\n", debug)
    fmt.Printf("Allowed Hosts: %v\n", allowedHosts)
}
```

## Usage Examples

### String Variables

```go
import "github.com/parthivrawat/typed-env-vars/go"

// Required string
apiKey := env.String("API_KEY")

// Optional string with default
appName := env.String("APP_NAME", "MyApp")
```

### Integer Variables

```go
// Required integer
port := env.Int("PORT")

// Optional integer with default
maxRetries := env.Int("MAX_RETRIES", 3)
```

### Float Variables

```go
// Required float
timeout := env.Float("TIMEOUT")

// Optional float with default
rateLimit := env.Float("RATE_LIMIT", 1.5)
```

### Boolean Variables

```go
// Accepts: true/false, yes/no, 1/0, on/off (case-insensitive)
debug := env.Bool("DEBUG", false)
enableCache := env.Bool("ENABLE_CACHE", true)
```

### List Variables

```go
// Comma-separated by default
allowedHosts := env.List("ALLOWED_HOSTS")

// Custom separator
tags := env.List("TAGS", ":")

// With default
features := env.List("FEATURES", ",", []string{"feature1", "feature2"})
```

### Map Variables

```go
// Format: KEY1=VALUE1,KEY2=VALUE2
featureFlags := env.Map("FEATURE_FLAGS")

// Custom separators
config := env.Map("CONFIG", ";", ":")

// With default
settings := env.Map("SETTINGS", ",", "=", map[string]string{"mode": "production"})
```

### URL Variables

```go
// Validates URL format
databaseURL := env.URL("DATABASE_URL")
apiEndpoint := env.URL("API_ENDPOINT", "https://api.example.com")
```

### Custom Type Conversion

```go
import "time"

releaseDate := env.Custom("RELEASE_DATE", func(s string) (interface{}, error) {
    return time.Parse("2006-01-02", s)
}).(time.Time)
```

## Real-World Example

```go
package main

import (
    "fmt"
    "github.com/parthivrawat/typed-env-vars/go"
)

type Config struct {
    // Application
    AppName     string
    Debug       bool
    
    // Server
    Host        string
    Port        int
    
    // Database
    DatabaseURL string
    DBPoolSize  int
    DBTimeout   float64
    DBSSL       bool
    
    // Security
    SecretKey     string
    AllowedHosts  []string
    CORSOrigins   []string
    
    // Features
    FeatureFlags map[string]string
}

func LoadConfig() *Config {
    return &Config{
        // Application
        AppName: env.String("APP_NAME", "MyApp"),
        Debug:   env.Bool("DEBUG", false),
        
        // Server
        Host: env.String("HOST", "0.0.0.0"),
        Port: env.Int("PORT", 8000),
        
        // Database
        DatabaseURL: env.URL("DATABASE_URL"),
        DBPoolSize:  env.Int("DB_POOL_SIZE", 10),
        DBTimeout:   env.Float("DB_TIMEOUT", 30.0),
        DBSSL:       env.Bool("DB_SSL", true),
        
        // Security
        SecretKey:    env.String("SECRET_KEY"),
        AllowedHosts: env.List("ALLOWED_HOSTS", ",", []string{"localhost"}),
        CORSOrigins:  env.List("CORS_ORIGINS", ",", []string{}),
        
        // Features
        FeatureFlags: env.Map("FEATURE_FLAGS", ",", "=", map[string]string{}),
    }
}

func main() {
    config := LoadConfig()
    
    if config.Debug {
        fmt.Printf("Running %s in debug mode\n", config.AppName)
    }
    
    fmt.Printf("Server listening on %s:%d\n", config.Host, config.Port)
}
```

## Error Handling

```go
import "github.com/parthivrawat/typed-env-vars/go"

func main() {
    defer func() {
        if r := recover(); r != nil {
            switch err := r.(type) {
            case *env.EnvVarNotFoundError:
                fmt.Printf("Missing required variable: %v\n", err)
            case *env.EnvVarTypeError:
                fmt.Printf("Invalid type: %v\n", err)
            default:
                panic(r)
            }
        }
    }()
    
    port := env.Int("PORT")
    fmt.Printf("Port: %d\n", port)
}
```

## API Reference

### Functions

- `env.String(key, default...)` - Get string value
- `env.Int(key, default...)` - Get integer value
- `env.Float(key, default...)` - Get float64 value
- `env.Bool(key, default...)` - Get boolean value
- `env.List(key, separator, default...)` - Get list value
- `env.Map(key, itemSep, kvSep, default...)` - Get map value
- `env.URL(key, default...)` - Get URL value with validation
- `env.Custom(key, converter, default...)` - Get value with custom converter

### Must Functions

- `env.MustString(key)` - Get required string (panics if missing)
- `env.MustInt(key)` - Get required int (panics if missing)
- `env.MustFloat(key)` - Get required float64 (panics if missing)
- `env.MustBool(key)` - Get required bool (panics if missing)
- `env.MustURL(key)` - Get required URL (panics if missing)

### Error Types

- `EnvVarError` - Base error type
- `EnvVarNotFoundError` - Raised when required variable is missing
- `EnvVarTypeError` - Raised when type conversion fails

## Testing

```bash
# Run tests
go test -v

# Run tests with coverage
go test -v -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Best Practices

1. **Load configuration at startup**: Initialize all config in one place
2. **Use struct for config**: Keep all environment variables in a struct
3. **Provide sensible defaults**: For non-critical settings
4. **Validate early**: Load configuration at application startup
5. **Document required variables**: In README or .env.example file

## Comparison with Other Libraries

| Feature | typed-env-vars | godotenv | envconfig | viper |
|---------|---------------|----------|-----------|-------|
| Type safety | ✅ | ❌ | ✅ | ✅ |
| Zero dependencies | ✅ | ✅ | ✅ | ❌ |
| Map support | ✅ | ❌ | ❌ | ✅ |
| URL validation | ✅ | ❌ | ❌ | ❌ |
| Custom converters | ✅ | ❌ | ✅ | ✅ |
| Simple API | ✅ | ✅ | ⚠️ | ❌ |

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
