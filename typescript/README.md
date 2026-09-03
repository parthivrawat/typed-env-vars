# Typed Environment Variables (TypeScript)

A type-safe environment variable library for TypeScript/Node.js that provides automatic type conversion and validation.

## Features

- ✅ **Type Safety**: Full TypeScript support with type inference
- ✅ **Rich Types**: Support for string, number, boolean, list, dict, enum, URL, and custom types
- ✅ **Default Values**: Optional default values for missing variables
- ✅ **Clear Errors**: Descriptive error messages for debugging
- ✅ **Zero Dependencies**: No external dependencies required
- ✅ **Production Ready**: Comprehensive test coverage

## Installation

```bash
npm install typed-env-vars
# or
yarn add typed-env-vars
# or
pnpm add typed-env-vars
```

## Quick Start

```typescript
import { env } from 'typed-env-vars';

const DATABASE_URL = env.str('DATABASE_URL');
const MAX_CONNECTIONS = env.int('MAX_CONNECTIONS', 10);
const DEBUG = env.bool('DEBUG', false);
const ALLOWED_HOSTS = env.list('ALLOWED_HOSTS', { default: ['localhost'] });
```

## Usage Examples

### String Variables

```typescript
import { env } from 'typed-env-vars';

const API_KEY = env.str('API_KEY');
const APP_NAME = env.str('APP_NAME', 'MyApp');
```

### Integer Variables

```typescript
const PORT = env.int('PORT');
const MAX_RETRIES = env.int('MAX_RETRIES', 3);
```

### Float Variables

```typescript
const TIMEOUT = env.float('TIMEOUT');
const RATE_LIMIT = env.float('RATE_LIMIT', 1.5);
```

### Boolean Variables

```typescript
const DEBUG = env.bool('DEBUG', false);
const ENABLE_CACHE = env.bool('ENABLE_CACHE', true);
```

### List Variables

```typescript
const ALLOWED_HOSTS = env.list('ALLOWED_HOSTS');

const TAGS = env.list('TAGS', { separator: ';' });

const FEATURES = env.list('FEATURES', { default: ['feature1', 'feature2'] });
```

### Dictionary Variables

```typescript
const FEATURE_FLAGS = env.dict('FEATURE_FLAGS');

const CONFIG = env.dict('CONFIG', {
  itemSeparator: ';',
  keyValueSeparator: ':',
});

const SETTINGS = env.dict('SETTINGS', {
  default: { mode: 'production' },
});
```

### Enum Variables

```typescript
enum LogLevel {
  DEBUG = 'DEBUG',
  INFO = 'INFO',
  WARNING = 'WARNING',
  ERROR = 'ERROR',
}

const LOG_LEVEL = env.enum('LOG_LEVEL', LogLevel, LogLevel.INFO);
```

### URL Variables

```typescript
// Validates URL format and accepts any valid scheme (e.g. postgresql://, redis://, amqp://)
const DATABASE_URL = env.url('DATABASE_URL');
const API_ENDPOINT = env.url('API_ENDPOINT', 'https://api.example.com');
```

### Custom Type Conversion

```typescript
const RELEASE_DATE = env.custom('RELEASE_DATE', (s) => {
  const [year, month, day] = s.split('-').map(Number);
  return new Date(year, month - 1, day);
});
```

## Real-World Example

```typescript
import { env } from 'typed-env-vars';

enum Environment {
  DEVELOPMENT = 'DEVELOPMENT',
  STAGING = 'STAGING',
  PRODUCTION = 'PRODUCTION',
}

class Config {
  static readonly APP_NAME = env.str('APP_NAME', 'MyApp');
  static readonly ENVIRONMENT = env.enum('ENVIRONMENT', Environment, Environment.DEVELOPMENT);
  static readonly DEBUG = env.bool('DEBUG', false);
  
  static readonly HOST = env.str('HOST', '0.0.0.0');
  static readonly PORT = env.int('PORT', 8000);
  
  static readonly DATABASE_URL = env.url('DATABASE_URL');
  static readonly DB_POOL_SIZE = env.int('DB_POOL_SIZE', 10);
  static readonly DB_TIMEOUT = env.float('DB_TIMEOUT', 30.0);
  static readonly DB_SSL = env.bool('DB_SSL', true);
  
  static readonly SECRET_KEY = env.str('SECRET_KEY');
  static readonly ALLOWED_HOSTS = env.list('ALLOWED_HOSTS', { default: ['localhost'] });
  static readonly CORS_ORIGINS = env.list('CORS_ORIGINS', { default: [] });
  
  static readonly FEATURE_FLAGS = env.dict('FEATURE_FLAGS', { default: {} });
  
  static readonly REDIS_URL = env.url('REDIS_URL', 'redis://localhost:6379/0');
  static readonly SMTP_HOST = env.str('SMTP_HOST', 'localhost');
  static readonly SMTP_PORT = env.int('SMTP_PORT', 587);
}

if (Config.DEBUG) {
  console.log(`Running ${Config.APP_NAME} in debug mode`);
}
```

## Error Handling

```typescript
import { env, EnvVarNotFoundError, EnvVarTypeError } from 'typed-env-vars';

try {
  const port = env.int('PORT');
} catch (error) {
  if (error instanceof EnvVarNotFoundError) {
    console.error('Missing required variable:', error.message);
  } else if (error instanceof EnvVarTypeError) {
    console.error('Invalid type:', error.message);
  }
}
```

## API Reference

`Env` is the class exported from the package, and `env` is a ready-to-use singleton instance. You can import either to create your own instance or to use the shared one.

### Methods

- `env.str(key, default?)` - Get string value
- `env.int(key, default?)` - Get integer value
- `env.float(key, default?)` - Get float value
- `env.bool(key, default?)` - Get boolean value
- `env.list(key, options?)` - Get list value
- `env.dict(key, options?)` - Get dict value
- `env.enum(key, enumObj, default?)` - Get enum value
- `env.url(key, default?)` - Get URL value with validation
- `env.custom(key, converter, default?)` - Get value with custom converter

### Exceptions

- `EnvVarError` - Base exception class
- `EnvVarNotFoundError` - Raised when required variable is missing
- `EnvVarTypeError` - Raised when type conversion fails

## Testing

```bash
npm test

npm run test:coverage

npm run test:watch
```

## Building

```bash
npm run build
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history.
