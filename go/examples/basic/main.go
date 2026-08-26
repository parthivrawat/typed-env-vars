package main

import (
	"fmt"
	"os"
	"strings"
	env "github.com/parthivrawat/typed-env-vars"
)

// Config holds all application configuration
type Config struct {
	// Application Settings
	AppName     string
	Debug       bool
	Version     string
	
	// Server Settings
	Host    string
	Port    int
	Workers int
	
	// Database Settings
	DatabaseURL   string
	DBPoolSize    int
	DBPoolTimeout float64
	DBSSLEnabled  bool
	
	// Redis Settings
	RedisURL           string
	RedisMaxConnections int
	
	// Security Settings
	SecretKey    string
	AllowedHosts []string
	CORSOrigins  []string
	
	// Feature Flags
	FeatureFlags map[string]string
	
	// External Services
	SMTPHost    string
	SMTPPort    int
	SMTPUseTLS  bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		// Application Settings
		AppName: env.String("APP_NAME", "MyApp"),
		Debug:   env.Bool("DEBUG", false),
		Version: env.String("VERSION", "1.0.0"),
		
		// Server Settings
		Host:    env.String("HOST", "0.0.0.0"),
		Port:    env.Int("PORT", 8000),
		Workers: env.Int("WORKERS", 4),
		
		// Database Settings
		DatabaseURL:   env.URL("DATABASE_URL", "postgresql://localhost:5432/myapp"),
		DBPoolSize:    env.Int("DB_POOL_SIZE", 10),
		DBPoolTimeout: env.Float("DB_POOL_TIMEOUT", 30.0),
		DBSSLEnabled:  env.Bool("DB_SSL_ENABLED", true),
		
		// Redis Settings
		RedisURL:           env.URL("REDIS_URL", "redis://localhost:6379/0"),
		RedisMaxConnections: env.Int("REDIS_MAX_CONNECTIONS", 50),
		
		// Security Settings
		SecretKey:    env.String("SECRET_KEY", "dev-secret-key-change-in-production"),
		AllowedHosts: env.List("ALLOWED_HOSTS", ",", []string{"localhost", "127.0.0.1"}),
		CORSOrigins:  env.List("CORS_ORIGINS", ",", []string{"http://localhost:3000"}),
		
		// Feature Flags
		FeatureFlags: env.Map("FEATURE_FLAGS", ",", "=", map[string]string{
			"new_ui":        "false",
			"beta_features": "false",
		}),
		
		// External Services
		SMTPHost:   env.String("SMTP_HOST", "localhost"),
		SMTPPort:   env.Int("SMTP_PORT", 587),
		SMTPUseTLS: env.Bool("SMTP_USE_TLS", true),
	}
}

// Display prints the current configuration (safe for logging)
func (c *Config) Display() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Application Configuration")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("App Name:        %s\n", c.AppName)
	fmt.Printf("Debug Mode:      %v\n", c.Debug)
	fmt.Printf("Version:         %s\n", c.Version)
	fmt.Printf("Host:Port:       %s:%d\n", c.Host, c.Port)
	fmt.Printf("Workers:         %d\n", c.Workers)
	fmt.Printf("Database:        %s\n", c.DatabaseURL)
	fmt.Printf("Redis:           %s\n", c.RedisURL)
	fmt.Printf("Allowed Hosts:   %v\n", c.AllowedHosts)
	fmt.Printf("Feature Flags:   %d flags\n", len(c.FeatureFlags))
	fmt.Println(strings.Repeat("=", 60))
}

func exampleBasicUsage() {
	fmt.Println("\n=== Example 1: Basic Usage ===")
	
	// Set some example environment variables
	os.Setenv("API_KEY", "secret-key-123")
	os.Setenv("MAX_RETRIES", "3")
	os.Setenv("TIMEOUT", "30.5")
	os.Setenv("ENABLE_CACHE", "true")
	
	// Load them with type safety
	apiKey := env.String("API_KEY")
	maxRetries := env.Int("MAX_RETRIES")
	timeout := env.Float("TIMEOUT")
	enableCache := env.Bool("ENABLE_CACHE")
	
	fmt.Printf("API Key: %s\n", apiKey)
	fmt.Printf("Max Retries: %d\n", maxRetries)
	fmt.Printf("Timeout: %.1f\n", timeout)
	fmt.Printf("Enable Cache: %v\n", enableCache)
}

func exampleWithDefaults() {
	fmt.Println("\n=== Example 2: Default Values ===")
	
	// These variables don't exist, so defaults are used
	port := env.Int("PORT_EXAMPLE", 8000)
	debug := env.Bool("DEBUG_EXAMPLE", false)
	allowedHosts := env.List("ALLOWED_HOSTS_EXAMPLE", ",", []string{"localhost"})
	
	fmt.Printf("Port: %d\n", port)
	fmt.Printf("Debug: %v\n", debug)
	fmt.Printf("Allowed Hosts: %v\n", allowedHosts)
}

func exampleListsAndMaps() {
	fmt.Println("\n=== Example 3: Lists and Maps ===")
	
	os.Setenv("TAGS", "go,rust,python")
	os.Setenv("SETTINGS", "theme=dark,lang=en,timezone=UTC")
	
	tags := env.List("TAGS")
	settings := env.Map("SETTINGS")
	
	fmt.Printf("Tags: %v\n", tags)
	fmt.Printf("Settings: %v\n", settings)
	fmt.Printf("  Theme: %s\n", settings["theme"])
	fmt.Printf("  Language: %s\n", settings["lang"])
}

func exampleURLValidation() {
	fmt.Println("\n=== Example 4: URL Validation ===")
	
	os.Setenv("DATABASE_URL", "postgresql://localhost:5432/mydb")
	os.Setenv("API_ENDPOINT", "https://api.example.com")
	
	dbURL := env.URL("DATABASE_URL")
	apiEndpoint := env.URL("API_ENDPOINT")
	
	fmt.Printf("Database URL: %s\n", dbURL)
	fmt.Printf("API Endpoint: %s\n", apiEndpoint)
}

func exampleErrorHandling() {
	fmt.Println("\n=== Example 5: Error Handling ===")
	
	// Missing required variable
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*env.EnvVarNotFoundError); ok {
				fmt.Printf("❌ Error: %v\n", err)
			}
		}
	}()
	
	_ = env.String("MISSING_REQUIRED_VAR")
}

func exampleApplicationConfig() {
	fmt.Println("\n=== Example 6: Application Configuration ===")
	
	// Set some environment variables
	os.Setenv("APP_NAME", "MyAwesomeApp")
	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgresql://localhost:5432/myapp")
	os.Setenv("REDIS_URL", "redis://localhost:6379/0")
	os.Setenv("ALLOWED_HOSTS", "localhost,127.0.0.1,example.com")
	os.Setenv("FEATURE_FLAGS", "new_ui=true,beta_features=false")
	
	// Load configuration
	config := LoadConfig()
	
	// Display configuration
	config.Display()
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Typed Environment Variables - Examples")
	fmt.Println(strings.Repeat("=", 60))
	
	exampleBasicUsage()
	exampleWithDefaults()
	exampleListsAndMaps()
	exampleURLValidation()
	exampleErrorHandling()
	exampleApplicationConfig()
	
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Examples completed!")
	fmt.Println(strings.Repeat("=", 60))
}
