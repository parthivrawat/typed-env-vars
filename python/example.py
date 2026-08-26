"""
Example usage of typed-env-vars library

This example demonstrates how to use the typed-env-vars library
to load and validate environment variables in a real application.
"""

import os
from enum import Enum
from typed_env import env, EnvVarNotFoundError, EnvVarTypeError


class Environment(Enum):
    """Application environment types."""
    DEVELOPMENT = 1
    STAGING = 2
    PRODUCTION = 3


class LogLevel(Enum):
    """Logging levels."""
    DEBUG = 1
    INFO = 2
    WARNING = 3
    ERROR = 4


class AppConfig:
    """
    Application configuration loaded from environment variables.
    
    This class demonstrates best practices for managing configuration:
    - All config in one place
    - Type-safe access
    - Sensible defaults for non-critical settings
    - Clear documentation
    """
    
    # Application Settings
    APP_NAME = env.str('APP_NAME', default='MyApp')
    ENVIRONMENT = env.enum('ENVIRONMENT', Environment, default=Environment.DEVELOPMENT)
    DEBUG = env.bool('DEBUG', default=False)
    VERSION = env.str('VERSION', default='1.0.0')
    
    # Server Settings
    HOST = env.str('HOST', default='0.0.0.0')
    PORT = env.int('PORT', default=8000)
    WORKERS = env.int('WORKERS', default=4)
    
    # Database Settings
    DATABASE_URL = env.url('DATABASE_URL', default='postgresql://localhost:5432/myapp')
    DB_POOL_SIZE = env.int('DB_POOL_SIZE', default=10)
    DB_POOL_TIMEOUT = env.float('DB_POOL_TIMEOUT', default=30.0)
    DB_SSL_ENABLED = env.bool('DB_SSL_ENABLED', default=True)
    
    # Redis Settings
    REDIS_URL = env.url('REDIS_URL', default='redis://localhost:6379/0')
    REDIS_MAX_CONNECTIONS = env.int('REDIS_MAX_CONNECTIONS', default=50)
    
    # Security Settings
    SECRET_KEY = env.str('SECRET_KEY', default='dev-secret-key-change-in-production')
    ALLOWED_HOSTS = env.list('ALLOWED_HOSTS', default=['localhost', '127.0.0.1'])
    CORS_ORIGINS = env.list('CORS_ORIGINS', default=['http://localhost:3000'])
    
    # Logging Settings
    LOG_LEVEL = env.enum('LOG_LEVEL', LogLevel, default=LogLevel.INFO)
    LOG_FORMAT = env.str('LOG_FORMAT', default='json')
    
    # Feature Flags
    FEATURE_FLAGS = env.dict('FEATURE_FLAGS', default={
        'new_ui': 'false',
        'beta_features': 'false'
    })
    
    # External Services
    SMTP_HOST = env.str('SMTP_HOST', default='localhost')
    SMTP_PORT = env.int('SMTP_PORT', default=587)
    SMTP_USER = env.str('SMTP_USER', default='')
    SMTP_PASSWORD = env.str('SMTP_PASSWORD', default='')
    SMTP_USE_TLS = env.bool('SMTP_USE_TLS', default=True)
    
    # API Keys (should be set in production)
    STRIPE_API_KEY = env.str('STRIPE_API_KEY', default='')
    SENDGRID_API_KEY = env.str('SENDGRID_API_KEY', default='')
    
    # Monitoring
    SENTRY_DSN = env.str('SENTRY_DSN', default='')
    ENABLE_METRICS = env.bool('ENABLE_METRICS', default=False)
    
    @classmethod
    def validate(cls):
        """
        Validate critical configuration.
        
        This method should be called at application startup to ensure
        all required configuration is present and valid.
        """
        errors = []
        
        # Check critical settings in production
        if cls.ENVIRONMENT == Environment.PRODUCTION:
            if cls.SECRET_KEY == 'dev-secret-key-change-in-production':
                errors.append("SECRET_KEY must be changed in production")
            
            if cls.DEBUG:
                errors.append("DEBUG must be False in production")
            
            if not cls.SENTRY_DSN:
                errors.append("SENTRY_DSN should be set in production")
        
        if errors:
            raise ValueError(f"Configuration errors:\n" + "\n".join(f"  - {e}" for e in errors))
    
    @classmethod
    def display(cls):
        """Display current configuration (safe for logging)."""
        print("=" * 60)
        print("Application Configuration")
        print("=" * 60)
        print(f"App Name:        {cls.APP_NAME}")
        print(f"Environment:     {cls.ENVIRONMENT.name}")
        print(f"Debug Mode:      {cls.DEBUG}")
        print(f"Version:         {cls.VERSION}")
        print(f"Host:Port:       {cls.HOST}:{cls.PORT}")
        print(f"Workers:         {cls.WORKERS}")
        print(f"Database:        {cls.DATABASE_URL.split('@')[-1] if '@' in cls.DATABASE_URL else cls.DATABASE_URL}")
        print(f"Redis:           {cls.REDIS_URL}")
        print(f"Log Level:       {cls.LOG_LEVEL.name}")
        print(f"Allowed Hosts:   {', '.join(cls.ALLOWED_HOSTS)}")
        print(f"Feature Flags:   {len(cls.FEATURE_FLAGS)} flags")
        print("=" * 60)


def example_basic_usage():
    """Example: Basic usage of environment variables."""
    print("\n=== Example 1: Basic Usage ===\n")
    
    # Set some example environment variables
    os.environ['API_KEY'] = 'secret-key-123'
    os.environ['MAX_RETRIES'] = '3'
    os.environ['TIMEOUT'] = '30.5'
    os.environ['ENABLE_CACHE'] = 'true'
    
    # Load them with type safety
    api_key = env.str('API_KEY')
    max_retries = env.int('MAX_RETRIES')
    timeout = env.float('TIMEOUT')
    enable_cache = env.bool('ENABLE_CACHE')
    
    print(f"API Key: {api_key}")
    print(f"Max Retries: {max_retries} (type: {type(max_retries).__name__})")
    print(f"Timeout: {timeout} (type: {type(timeout).__name__})")
    print(f"Enable Cache: {enable_cache} (type: {type(enable_cache).__name__})")


def example_with_defaults():
    """Example: Using default values."""
    print("\n=== Example 2: Default Values ===\n")
    
    # These variables don't exist, so defaults are used
    port = env.int('PORT', default=8000)
    debug = env.bool('DEBUG', default=False)
    allowed_hosts = env.list('ALLOWED_HOSTS', default=['localhost'])
    
    print(f"Port: {port}")
    print(f"Debug: {debug}")
    print(f"Allowed Hosts: {allowed_hosts}")


def example_lists_and_dicts():
    """Example: Working with lists and dictionaries."""
    print("\n=== Example 3: Lists and Dictionaries ===\n")
    
    os.environ['TAGS'] = 'python,typescript,go'
    os.environ['SETTINGS'] = 'theme=dark,lang=en,timezone=UTC'
    
    tags = env.list('TAGS')
    settings = env.dict('SETTINGS')
    
    print(f"Tags: {tags}")
    print(f"Settings: {settings}")
    print(f"  Theme: {settings['theme']}")
    print(f"  Language: {settings['lang']}")


def example_enums():
    """Example: Using enums for type-safe values."""
    print("\n=== Example 4: Enums ===\n")
    
    os.environ['LOG_LEVEL'] = 'INFO'
    os.environ['ENVIRONMENT'] = 'production'
    
    log_level = env.enum('LOG_LEVEL', LogLevel)
    environment = env.enum('ENVIRONMENT', Environment)
    
    print(f"Log Level: {log_level.name} (value: {log_level.value})")
    print(f"Environment: {environment.name} (value: {environment.value})")


def example_error_handling():
    """Example: Error handling."""
    print("\n=== Example 5: Error Handling ===\n")
    
    # Missing required variable
    try:
        missing = env.str('MISSING_REQUIRED_VAR')
    except EnvVarNotFoundError as e:
        print(f"❌ Error: {e}")
    
    # Invalid type conversion
    os.environ['INVALID_INT'] = 'not-a-number'
    try:
        invalid = env.int('INVALID_INT')
    except EnvVarTypeError as e:
        print(f"❌ Error: {e}")
    
    # Invalid boolean
    os.environ['INVALID_BOOL'] = 'maybe'
    try:
        invalid = env.bool('INVALID_BOOL')
    except EnvVarTypeError as e:
        print(f"❌ Error: {e}")


def example_custom_converter():
    """Example: Custom type converter."""
    print("\n=== Example 6: Custom Converter ===\n")
    
    from datetime import datetime
    
    os.environ['RELEASE_DATE'] = '2024-01-15'
    
    def parse_date(s):
        return datetime.strptime(s, '%Y-%m-%d').date()
    
    release_date = env.custom('RELEASE_DATE', parse_date)
    print(f"Release Date: {release_date}")
    print(f"Type: {type(release_date).__name__}")


def example_application_config():
    """Example: Complete application configuration."""
    print("\n=== Example 7: Application Configuration ===\n")
    
    # Set some environment variables
    os.environ['APP_NAME'] = 'MyAwesomeApp'
    os.environ['ENVIRONMENT'] = 'DEVELOPMENT'
    os.environ['DEBUG'] = 'true'
    os.environ['PORT'] = '3000'
    os.environ['DATABASE_URL'] = 'postgresql://localhost:5432/myapp'
    os.environ['REDIS_URL'] = 'redis://localhost:6379/0'
    os.environ['LOG_LEVEL'] = 'DEBUG'
    os.environ['ALLOWED_HOSTS'] = 'localhost,127.0.0.1,example.com'
    os.environ['FEATURE_FLAGS'] = 'new_ui=true,beta_features=false'
    
    # Display configuration
    AppConfig.display()
    
    # Validate configuration
    try:
        AppConfig.validate()
        print("\n✅ Configuration is valid!")
    except ValueError as e:
        print(f"\n❌ Configuration errors:\n{e}")


def main():
    """Run all examples."""
    print("=" * 60)
    print("Typed Environment Variables - Examples")
    print("=" * 60)
    
    example_basic_usage()
    example_with_defaults()
    example_lists_and_dicts()
    example_enums()
    example_error_handling()
    example_custom_converter()
    example_application_config()
    
    print("\n" + "=" * 60)
    print("Examples completed!")
    print("=" * 60)


if __name__ == '__main__':
    main()
