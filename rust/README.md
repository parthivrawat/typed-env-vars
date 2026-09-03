# Typed Environment Variables (Rust)

A type-safe environment variable library for Rust that provides automatic type conversion and validation.

## Features

- ✅ **Type Safety**: Automatic type conversion with validation
- ✅ **Rich Types**: Support for String, i32, i64, f64, bool, Vec, HashMap, enum, URL, and custom types
- ✅ **Cross-Language Vocabulary**: `dict` alias for `map` and a shared API shape
- ✅ **Default Values**: Optional default values using `unwrap_or`
- ✅ **Clear Errors**: Descriptive error messages for debugging
- ✅ **Zero Dependencies**: No external dependencies required
- ✅ **Production Ready**: Comprehensive test coverage

## Installation

Add this to your `Cargo.toml`:

```toml
[dependencies]
typed-env-vars = "1.1"
```

## Quick Start

```rust
use typed_env_vars::EnvVar;

fn main() {
    let database_url = EnvVar::url("DATABASE_URL").expect("DATABASE_URL must be set");
    let max_connections = EnvVar::int("MAX_CONNECTIONS").unwrap_or(10);
    let debug = EnvVar::bool("DEBUG").unwrap_or(false);
    let allowed_hosts = EnvVar::list("ALLOWED_HOSTS", ",").unwrap_or_default();
    
    println!("Database: {}", database_url);
    println!("Max Connections: {}", max_connections);
    println!("Debug: {}", debug);
    println!("Allowed Hosts: {:?}", allowed_hosts);
}
```

## Usage Examples

### String Variables

```rust
use typed_env_vars::EnvVar;

// Required string
let api_key = EnvVar::string("API_KEY").expect("API_KEY is required");

// Optional string with default
let app_name = EnvVar::string("APP_NAME").unwrap_or_else(|_| "MyApp".to_string());
```

### Integer Variables

```rust
// Required integer
let port = EnvVar::int("PORT").expect("PORT is required");

// Optional integer with default
let max_retries = EnvVar::int("MAX_RETRIES").unwrap_or(3);

// 64-bit integer
let max_size = EnvVar::int64("MAX_SIZE").unwrap_or(1000000);
```

### Float Variables

```rust
// Required float
let timeout = EnvVar::float("TIMEOUT").expect("TIMEOUT is required");

// Optional float with default
let rate_limit = EnvVar::float("RATE_LIMIT").unwrap_or(1.5);
```

### Boolean Variables

```rust
// Accepts: true/false, yes/no, 1/0, on/off (case-insensitive)
let debug = EnvVar::bool("DEBUG").unwrap_or(false);
let enable_cache = EnvVar::bool("ENABLE_CACHE").unwrap_or(true);
```

### List Variables

```rust
// Comma-separated by default
let allowed_hosts = EnvVar::list("ALLOWED_HOSTS", ",").unwrap_or_default();
// Example: ALLOWED_HOSTS=localhost,127.0.0.1,example.com

// Custom separator
let tags = EnvVar::list("TAGS", ";").unwrap_or_default();
// Example: TAGS=tag1;tag2;tag3
```

### Map Variables

```rust
use std::collections::HashMap;

// Format: KEY1=VALUE1,KEY2=VALUE2
let feature_flags = EnvVar::map("FEATURE_FLAGS", ",", "=").unwrap_or_default();
// Example: FEATURE_FLAGS=feature1=true,feature2=false

// Custom separators
let config = EnvVar::map("CONFIG", ";", ":").unwrap_or_default();
// Example: CONFIG=key1:value1;key2:value2
```

### URL Variables

```rust
// Validates URL format and accepts any valid scheme (e.g. postgresql://, redis://, amqp://)
let database_url = EnvVar::url("DATABASE_URL").expect("DATABASE_URL is required");
// Example: DATABASE_URL=postgresql://localhost:5432/mydb

let api_endpoint = EnvVar::url("API_ENDPOINT")
    .unwrap_or_else(|_| "https://api.example.com".to_string());
```

### Dict Alias

```rust
use std::collections::HashMap;

// `dict` is an alias for `map` that matches the cross-language vocabulary
let feature_flags = EnvVar::dict("FEATURE_FLAGS", ",", "=").unwrap_or_default();
```

### Enum Variables

```rust
use std::collections::HashMap;

#[derive(Clone, Debug, PartialEq)]
enum Environment {
    Development,
    Staging,
    Production,
}

let mut values = HashMap::new();
values.insert("development".to_string(), Environment::Development);
values.insert("staging".to_string(), Environment::Staging);
values.insert("production".to_string(), Environment::Production);

let environment = EnvVar::enum_value("ENVIRONMENT", &values, Some(Environment::Development))
    .unwrap();
```

### Custom Type Conversion

```rust
let timeout = EnvVar::custom("TIMEOUT_MS", |s| {
    let timeout_ms: u64 = s.parse().map_err(|e| Box::new(e) as Box<dyn std::error::Error>)?;
    Ok(std::time::Duration::from_millis(timeout_ms))
}).unwrap_or_else(|_| std::time::Duration::from_secs(30));
```

## Real-World Example

```rust
use typed_env_vars::EnvVar;
use std::collections::HashMap;

struct Config {
    // Application
    app_name: String,
    debug: bool,
    
    // Server
    host: String,
    port: i32,
    
    // Database
    database_url: String,
    db_pool_size: i32,
    db_timeout: f64,
    db_ssl: bool,
    
    // Security
    secret_key: String,
    allowed_hosts: Vec<String>,
    cors_origins: Vec<String>,
    
    // Features
    feature_flags: HashMap<String, String>,
}

impl Config {
    fn from_env() -> Result<Self, Box<dyn std::error::Error>> {
        Ok(Config {
            // Application
            app_name: EnvVar::string("APP_NAME").unwrap_or_else(|_| "MyApp".to_string()),
            debug: EnvVar::bool("DEBUG").unwrap_or(false),
            
            // Server
            host: EnvVar::string("HOST").unwrap_or_else(|_| "0.0.0.0".to_string()),
            port: EnvVar::int("PORT").unwrap_or(8000),
            
            // Database
            database_url: EnvVar::url("DATABASE_URL")?,
            db_pool_size: EnvVar::int("DB_POOL_SIZE").unwrap_or(10),
            db_timeout: EnvVar::float("DB_TIMEOUT").unwrap_or(30.0),
            db_ssl: EnvVar::bool("DB_SSL").unwrap_or(true),
            
            // Security
            secret_key: EnvVar::string("SECRET_KEY")?,
            allowed_hosts: EnvVar::list("ALLOWED_HOSTS", ",")
                .unwrap_or_else(|_| vec!["localhost".to_string()]),
            cors_origins: EnvVar::list("CORS_ORIGINS", ",").unwrap_or_default(),
            
            // Features
            feature_flags: EnvVar::map("FEATURE_FLAGS", ",", "=").unwrap_or_default(),
        })
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let config = Config::from_env()?;
    
    if config.debug {
        println!("Running {} in debug mode", config.app_name);
    }
    
    println!("Server listening on {}:{}", config.host, config.port);
    Ok(())
}
```

## Error Handling

```rust
use typed_env_vars::{EnvVar, EnvVarError};

fn main() {
    match EnvVar::int("PORT") {
        Ok(port) => println!("Port: {}", port),
        Err(EnvVarError::NotFound(key)) => {
            eprintln!("Missing required variable: {}", key);
        }
        Err(EnvVarError::TypeError { key, value, expected_type, details }) => {
            eprintln!("Invalid type for {}: {} (expected {})", key, value, expected_type);
            if let Some(d) = details {
                eprintln!("Details: {}", d);
            }
        }
    }
}
```

## API Reference

### Methods

- `EnvVar::string(key)` - Get string value
- `EnvVar::int(key)` - Get i32 value
- `EnvVar::int64(key)` - Get i64 value
- `EnvVar::float(key)` - Get f64 value
- `EnvVar::bool(key)` - Get boolean value
- `EnvVar::list(key, separator)` - Get Vec<String> value
- `EnvVar::map(key, item_sep, kv_sep)` - Get HashMap<String, String> value
- `EnvVar::dict(key, item_sep, kv_sep)` - Alias for `map`
- `EnvVar::enum_value(key, values, default)` - Get enum value with a `HashMap<String, T>` lookup
- `EnvVar::url(key)` - Get URL value with validation
- `EnvVar::custom(key, converter)` - Get value with custom converter

### Error Types

- `EnvVarError::NotFound` - Environment variable not found
- `EnvVarError::TypeError` - Type conversion failed

## Testing

```bash
# Run tests
cargo test

# Run tests with output
cargo test -- --nocapture

# Run tests with coverage
cargo tarpaulin --out Html
```

## Best Practices

1. **Load configuration at startup**: Initialize all config in one place
2. **Use struct for config**: Keep all environment variables in a struct
3. **Provide sensible defaults**: For non-critical settings using `unwrap_or`
4. **Validate early**: Load configuration at application startup
5. **Document required variables**: In README or .env.example file

## Comparison with Other Libraries

| Feature | typed-env-vars | dotenv | envy | config |
|---------|---------------|--------|------|--------|
| Type safety | ✅ | ❌ | ✅ | ✅ |
| Zero dependencies | ✅ | ❌ | ❌ | ❌ |
| Map support | ✅ | ❌ | ❌ | ✅ |
| URL validation | ✅ | ❌ | ❌ | ❌ |
| Custom converters | ✅ | ❌ | ✅ | ✅ |
| Simple API | ✅ | ✅ | ⚠️ | ❌ |

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history.
