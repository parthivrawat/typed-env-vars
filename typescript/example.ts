/**
 * Example usage of typed-env-vars library
 * 
 * This example demonstrates how to use the typed-env-vars library
 * to load and validate environment variables in a real application.
 */

import { env, EnvVarNotFoundError, EnvVarTypeError } from './src/index';

enum Environment {
  DEVELOPMENT = 'DEVELOPMENT',
  STAGING = 'STAGING',
  PRODUCTION = 'PRODUCTION',
}

enum LogLevel {
  DEBUG = 'DEBUG',
  INFO = 'INFO',
  WARNING = 'WARNING',
  ERROR = 'ERROR',
}

/**
 * Application configuration loaded from environment variables.
 * 
 * This class demonstrates best practices for managing configuration:
 * - All config in one place
 * - Type-safe access
 * - Sensible defaults for non-critical settings
 * - Clear documentation
 */
class AppConfig {
  // Application Settings
  static readonly APP_NAME = env.str('APP_NAME', 'MyApp');
  static readonly ENVIRONMENT = env.enum('ENVIRONMENT', Environment, Environment.DEVELOPMENT);
  static readonly DEBUG = env.bool('DEBUG', false);
  static readonly VERSION = env.str('VERSION', '1.0.0');
  
  // Server Settings
  static readonly HOST = env.str('HOST', '0.0.0.0');
  static readonly PORT = env.int('PORT', 8000);
  static readonly WORKERS = env.int('WORKERS', 4);
  
  // Database Settings
  static readonly DATABASE_URL = env.url('DATABASE_URL', 'postgresql://localhost:5432/myapp');
  static readonly DB_POOL_SIZE = env.int('DB_POOL_SIZE', 10);
  static readonly DB_POOL_TIMEOUT = env.float('DB_POOL_TIMEOUT', 30.0);
  static readonly DB_SSL_ENABLED = env.bool('DB_SSL_ENABLED', true);
  
  // Redis Settings
  static readonly REDIS_URL = env.url('REDIS_URL', 'redis://localhost:6379/0');
  static readonly REDIS_MAX_CONNECTIONS = env.int('REDIS_MAX_CONNECTIONS', 50);
  
  // Security Settings
  static readonly SECRET_KEY = env.str('SECRET_KEY', 'dev-secret-key-change-in-production');
  static readonly ALLOWED_HOSTS = env.list('ALLOWED_HOSTS', { default: ['localhost', '127.0.0.1'] });
  static readonly CORS_ORIGINS = env.list('CORS_ORIGINS', { default: ['http://localhost:3000'] });
  
  // Logging Settings
  static readonly LOG_LEVEL = env.enum('LOG_LEVEL', LogLevel, LogLevel.INFO);
  static readonly LOG_FORMAT = env.str('LOG_FORMAT', 'json');
  
  // Feature Flags
  static readonly FEATURE_FLAGS = env.dict('FEATURE_FLAGS', {
    default: {
      new_ui: 'false',
      beta_features: 'false',
    },
  });
  
  // External Services
  static readonly SMTP_HOST = env.str('SMTP_HOST', 'localhost');
  static readonly SMTP_PORT = env.int('SMTP_PORT', 587);
  static readonly SMTP_USER = env.str('SMTP_USER', '');
  static readonly SMTP_PASSWORD = env.str('SMTP_PASSWORD', '');
  static readonly SMTP_USE_TLS = env.bool('SMTP_USE_TLS', true);
  
  // API Keys
  static readonly STRIPE_API_KEY = env.str('STRIPE_API_KEY', '');
  static readonly SENDGRID_API_KEY = env.str('SENDGRID_API_KEY', '');
  
  // Monitoring
  static readonly SENTRY_DSN = env.str('SENTRY_DSN', '');
  static readonly ENABLE_METRICS = env.bool('ENABLE_METRICS', false);
  
  /**
   * Validate critical configuration.
   * 
   * This method should be called at application startup to ensure
   * all required configuration is present and valid.
   */
  static validate(): void {
    const errors: string[] = [];
    
    // Check critical settings in production
    if (this.ENVIRONMENT === Environment.PRODUCTION) {
      if (this.SECRET_KEY === 'dev-secret-key-change-in-production') {
        errors.push('SECRET_KEY must be changed in production');
      }
      
      if (this.DEBUG) {
        errors.push('DEBUG must be false in production');
      }
      
      if (!this.SENTRY_DSN) {
        errors.push('SENTRY_DSN should be set in production');
      }
    }
    
    if (errors.length > 0) {
      throw new Error(`Configuration errors:\n${errors.map(e => `  - ${e}`).join('\n')}`);
    }
  }
  
  /**
   * Display current configuration (safe for logging).
   */
  static display(): void {
    console.log('='.repeat(60));
    console.log('Application Configuration');
    console.log('='.repeat(60));
    console.log(`App Name:        ${this.APP_NAME}`);
    console.log(`Environment:     ${this.ENVIRONMENT}`);
    console.log(`Debug Mode:      ${this.DEBUG}`);
    console.log(`Version:         ${this.VERSION}`);
    console.log(`Host:Port:       ${this.HOST}:${this.PORT}`);
    console.log(`Workers:         ${this.WORKERS}`);
    console.log(`Database:        ${this.DATABASE_URL.includes('@') ? this.DATABASE_URL.split('@')[1] : this.DATABASE_URL}`);
    console.log(`Redis:           ${this.REDIS_URL}`);
    console.log(`Log Level:       ${this.LOG_LEVEL}`);
    console.log(`Allowed Hosts:   ${this.ALLOWED_HOSTS.join(', ')}`);
    console.log(`Feature Flags:   ${Object.keys(this.FEATURE_FLAGS).length} flags`);
    console.log('='.repeat(60));
  }
}

function exampleBasicUsage(): void {
  console.log('\n=== Example 1: Basic Usage ===\n');
  
  // Set some example environment variables
  process.env.API_KEY = 'secret-key-123';
  process.env.MAX_RETRIES = '3';
  process.env.TIMEOUT = '30.5';
  process.env.ENABLE_CACHE = 'true';
  
  // Load them with type safety
  const apiKey = env.str('API_KEY');
  const maxRetries = env.int('MAX_RETRIES');
  const timeout = env.float('TIMEOUT');
  const enableCache = env.bool('ENABLE_CACHE');
  
  console.log(`API Key: ${apiKey}`);
  console.log(`Max Retries: ${maxRetries} (type: ${typeof maxRetries})`);
  console.log(`Timeout: ${timeout} (type: ${typeof timeout})`);
  console.log(`Enable Cache: ${enableCache} (type: ${typeof enableCache})`);
}

function exampleWithDefaults(): void {
  console.log('\n=== Example 2: Default Values ===\n');
  
  // These variables don't exist, so defaults are used
  const port = env.int('PORT_EXAMPLE', 8000);
  const debug = env.bool('DEBUG_EXAMPLE', false);
  const allowedHosts = env.list('ALLOWED_HOSTS_EXAMPLE', { default: ['localhost'] });
  
  console.log(`Port: ${port}`);
  console.log(`Debug: ${debug}`);
  console.log(`Allowed Hosts: ${allowedHosts}`);
}

function exampleListsAndDicts(): void {
  console.log('\n=== Example 3: Lists and Dictionaries ===\n');
  
  process.env.TAGS = 'typescript,javascript,node';
  process.env.SETTINGS = 'theme=dark,lang=en,timezone=UTC';
  
  const tags = env.list('TAGS');
  const settings = env.dict('SETTINGS');
  
  console.log(`Tags: ${tags}`);
  console.log(`Settings:`, settings);
  console.log(`  Theme: ${settings.theme}`);
  console.log(`  Language: ${settings.lang}`);
}

function exampleEnums(): void {
  console.log('\n=== Example 4: Enums ===\n');
  
  process.env.LOG_LEVEL = 'INFO';
  process.env.ENVIRONMENT = 'PRODUCTION';
  
  const logLevel = env.enum('LOG_LEVEL', LogLevel);
  const environment = env.enum('ENVIRONMENT', Environment);
  
  console.log(`Log Level: ${logLevel}`);
  console.log(`Environment: ${environment}`);
}

function exampleErrorHandling(): void {
  console.log('\n=== Example 5: Error Handling ===\n');
  
  // Missing required variable
  try {
    const missing = env.str('MISSING_REQUIRED_VAR');
  } catch (error) {
    if (error instanceof EnvVarNotFoundError) {
      console.log(`❌ Error: ${error.message}`);
    }
  }
  
  // Invalid type conversion
  process.env.INVALID_INT = 'not-a-number';
  try {
    const invalid = env.int('INVALID_INT');
  } catch (error) {
    if (error instanceof EnvVarTypeError) {
      console.log(`❌ Error: ${error.message}`);
    }
  }
  
  // Invalid boolean
  process.env.INVALID_BOOL = 'maybe';
  try {
    const invalid = env.bool('INVALID_BOOL');
  } catch (error) {
    if (error instanceof EnvVarTypeError) {
      console.log(`❌ Error: ${error.message}`);
    }
  }
}

function exampleCustomConverter(): void {
  console.log('\n=== Example 6: Custom Converter ===\n');
  
  process.env.RELEASE_DATE = '2024-01-15';
  
  const parseDate = (s: string): Date => {
    const [year, month, day] = s.split('-').map(Number);
    return new Date(year, month - 1, day);
  };
  
  const releaseDate = env.custom('RELEASE_DATE', parseDate);
  console.log(`Release Date: ${releaseDate.toISOString().split('T')[0]}`);
  console.log(`Type: ${releaseDate.constructor.name}`);
}

function exampleApplicationConfig(): void {
  console.log('\n=== Example 7: Application Configuration ===\n');
  
  // Set some environment variables
  process.env.APP_NAME = 'MyAwesomeApp';
  process.env.ENVIRONMENT = 'DEVELOPMENT';
  process.env.DEBUG = 'true';
  process.env.PORT = '3000';
  process.env.DATABASE_URL = 'postgresql://localhost:5432/myapp';
  process.env.REDIS_URL = 'redis://localhost:6379/0';
  process.env.LOG_LEVEL = 'DEBUG';
  process.env.ALLOWED_HOSTS = 'localhost,127.0.0.1,example.com';
  process.env.FEATURE_FLAGS = 'new_ui=true,beta_features=false';
  
  // Display configuration
  AppConfig.display();
  
  // Validate configuration
  try {
    AppConfig.validate();
    console.log('\n✅ Configuration is valid!');
  } catch (error) {
    if (error instanceof Error) {
      console.log(`\n❌ Configuration errors:\n${error.message}`);
    }
  }
}

function main(): void {
  console.log('='.repeat(60));
  console.log('Typed Environment Variables - Examples');
  console.log('='.repeat(60));
  
  exampleBasicUsage();
  exampleWithDefaults();
  exampleListsAndDicts();
  exampleEnums();
  exampleErrorHandling();
  exampleCustomConverter();
  exampleApplicationConfig();
  
  console.log('\n' + '='.repeat(60));
  console.log('Examples completed!');
  console.log('='.repeat(60));
}

// Run examples if this file is executed directly
if (require.main === module) {
  main();
}

export { AppConfig };
