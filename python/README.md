# Typed Environment Variables

A type-safe environment variable library for Python that provides automatic type conversion and validation.

## Features

- ✅ **Type Safety**: Automatic type conversion with validation
- ✅ **Rich Types**: Support for str, int, float, bool, list, dict, enum, URL, and custom types
- ✅ **Default Values**: Optional default values for missing variables
- ✅ **Clear Errors**: Descriptive error messages for debugging
- ✅ **Zero Dependencies**: No external dependencies required
- ✅ **Production Ready**: Comprehensive test coverage

## Installation

```bash
pip install typed-env-vars
```

## Quick Start

```python
from typed_env import env

DATABASE_URL = env.str('DATABASE_URL')
MAX_CONNECTIONS = env.int('MAX_CONNECTIONS', default=10)
DEBUG = env.bool('DEBUG', default=False)
ALLOWED_HOSTS = env.list('ALLOWED_HOSTS', default=['localhost'])
```

## Usage Examples

### String Variables

```python
from typed_env import env

# Required string
API_KEY = env.str('API_KEY')

# Optional string with default
APP_NAME = env.str('APP_NAME', default='MyApp')
```

### Integer Variables

```python
# Required integer
PORT = env.int('PORT')

# Optional integer with default
MAX_RETRIES = env.int('MAX_RETRIES', default=3)
```

### Float Variables

```python
# Required float
TIMEOUT = env.float('TIMEOUT')

# Optional float with default
RATE_LIMIT = env.float('RATE_LIMIT', default=1.5)
```

### Boolean Variables

```python
# Accepts: true/false, yes/no, 1/0, on/off (case-insensitive)
DEBUG = env.bool('DEBUG', default=False)
ENABLE_CACHE = env.bool('ENABLE_CACHE', default=True)
```

### List Variables

```python
# Comma-separated by default
ALLOWED_HOSTS = env.list('ALLOWED_HOSTS')
# Example: ALLOWED_HOSTS=localhost,127.0.0.1,example.com

# Custom separator
TAGS = env.list('TAGS', separator=';')
# Example: TAGS=tag1;tag2;tag3

# With default
FEATURES = env.list('FEATURES', default=['feature1', 'feature2'])
```

### Dictionary Variables

```python
# Format: KEY1=VALUE1,KEY2=VALUE2
FEATURE_FLAGS = env.dict('FEATURE_FLAGS')
# Example: FEATURE_FLAGS=feature1=true,feature2=false

# Custom separators
CONFIG = env.dict('CONFIG', item_separator=';', key_value_separator=':')
# Example: CONFIG=key1:value1;key2:value2

# With default
SETTINGS = env.dict('SETTINGS', default={'mode': 'production'})
```

### Enum Variables

```python
from enum import Enum
from typed_env import env

class LogLevel(Enum):
    DEBUG = 1
    INFO = 2
    WARNING = 3
    ERROR = 4

# Case-insensitive enum matching
LOG_LEVEL = env.enum('LOG_LEVEL', LogLevel, default=LogLevel.INFO)
# Example: LOG_LEVEL=DEBUG or LOG_LEVEL=debug
```

### URL Variables

```python
# Validates URL format
DATABASE_URL = env.url('DATABASE_URL')
# Example: DATABASE_URL=postgresql://localhost:5432/mydb

API_ENDPOINT = env.url('API_ENDPOINT', default='https://api.example.com')
```

### Custom Type Conversion

```python
from datetime import datetime

def parse_date(s):
    return datetime.strptime(s, '%Y-%m-%d').date()

RELEASE_DATE = env.custom('RELEASE_DATE', parse_date)
# Example: RELEASE_DATE=2024-01-15
```

## Real-World Example

```python
from enum import Enum
from typed_env import env

class Environment(Enum):
    DEVELOPMENT = 1
    STAGING = 2
    PRODUCTION = 3

class Config:
    # Application
    APP_NAME = env.str('APP_NAME', default='MyApp')
    ENVIRONMENT = env.enum('ENVIRONMENT', Environment, default=Environment.DEVELOPMENT)
    DEBUG = env.bool('DEBUG', default=False)
    
    # Server
    HOST = env.str('HOST', default='0.0.0.0')
    PORT = env.int('PORT', default=8000)
    
    # Database
    DATABASE_URL = env.url('DATABASE_URL')
    DB_POOL_SIZE = env.int('DB_POOL_SIZE', default=10)
    DB_TIMEOUT = env.float('DB_TIMEOUT', default=30.0)
    DB_SSL = env.bool('DB_SSL', default=True)
    
    # Security
    SECRET_KEY = env.str('SECRET_KEY')
    ALLOWED_HOSTS = env.list('ALLOWED_HOSTS', default=['localhost'])
    CORS_ORIGINS = env.list('CORS_ORIGINS', default=[])
    
    # Features
    FEATURE_FLAGS = env.dict('FEATURE_FLAGS', default={})
    
    # External Services
    REDIS_URL = env.url('REDIS_URL', default='redis://localhost:6379/0')
    SMTP_HOST = env.str('SMTP_HOST', default='localhost')
    SMTP_PORT = env.int('SMTP_PORT', default=587)

# Usage
if Config.DEBUG:
    print(f"Running {Config.APP_NAME} in debug mode")
```

## Error Handling

```python
from typed_env import env, EnvVarNotFoundError, EnvVarTypeError

try:
    port = env.int('PORT')
except EnvVarNotFoundError as e:
    print(f"Missing required variable: {e}")
except EnvVarTypeError as e:
    print(f"Invalid type: {e}")
```

## API Reference

### Methods

- `env.str(key, default=None)` - Get string value
- `env.int(key, default=None)` - Get integer value
- `env.float(key, default=None)` - Get float value
- `env.bool(key, default=None)` - Get boolean value
- `env.list(key, default=None, separator=',')` - Get list value
- `env.dict(key, default=None, item_separator=',', key_value_separator='=')` - Get dict value
- `env.enum(key, enum_class, default=None)` - Get enum value
- `env.url(key, default=None)` - Get URL value with validation
- `env.custom(key, converter, default=None)` - Get value with custom converter

### Exceptions

- `EnvVarError` - Base exception class
- `EnvVarNotFoundError` - Raised when required variable is missing
- `EnvVarTypeError` - Raised when type conversion fails

## Testing

```bash
# Install dev dependencies
pip install -e ".[dev]"

# Run tests
pytest test_typed_env.py -v

# Run with coverage
pytest test_typed_env.py --cov=typed_env --cov-report=html
```

## Best Practices

1. **Define configuration in a class**: Keep all environment variables in one place
2. **Use type hints**: Make your configuration self-documenting
3. **Provide sensible defaults**: For non-critical settings
4. **Validate early**: Load configuration at application startup
5. **Document required variables**: In README or .env.example file

## Comparison with Other Libraries

| Feature | typed-env-vars | python-decouple | environs |
|---------|---------------|-----------------|----------|
| Type safety | ✅ | ✅ | ✅ |
| Zero dependencies | ✅ | ✅ | ❌ |
| Dict support | ✅ | ❌ | ✅ |
| Enum support | ✅ | ❌ | ✅ |
| URL validation | ✅ | ❌ | ✅ |
| Custom converters | ✅ | ✅ | ✅ |

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Changelog

### 1.0.0 (2024-01-15)
- Initial release
- Support for str, int, float, bool, list, dict, enum, URL types
- Custom type converters
- Comprehensive test coverage
