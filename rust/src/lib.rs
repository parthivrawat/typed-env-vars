//! Typed Environment Variables
//!
//! A type-safe environment variable library for Rust that provides automatic type conversion
//! and validation.
//!
//! # Features
//!
//! - **Type Safety**: Automatic type conversion with validation
//! - **Rich Types**: Support for String, i32, i64, f64, bool, Vec, HashMap, URL
//! - **Default Values**: Optional default values for missing variables
//! - **Clear Errors**: Descriptive error messages for debugging
//! - **Zero Dependencies**: No external dependencies required
//!
//! # Examples
//!
//! ```
//! use typed_env_vars::EnvVar;
//! use std::env;
//!
//! // Set environment variables for the example
//! env::set_var("API_KEY", "secret-key-123");
//!
//! // Required string
//! let api_key = EnvVar::string("API_KEY").unwrap();
//!
//! // Optional with default
//! let port = EnvVar::int("PORT").unwrap_or(8000);
//! let debug = EnvVar::bool("DEBUG").unwrap_or(false);
//! ```

use std::collections::HashMap;
use std::env;
use std::fmt;
use std::num::{ParseFloatError, ParseIntError};
use std::str::ParseBoolError;

/// Error types for environment variable operations
#[derive(Debug)]
pub enum EnvVarError {
    /// Environment variable not found
    NotFound(String),
    /// Type conversion error
    TypeError {
        key: String,
        value: String,
        expected_type: String,
        details: Option<String>,
    },
}

impl fmt::Display for EnvVarError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            EnvVarError::NotFound(key) => {
                write!(f, "Required environment variable '{}' is not set", key)
            }
            EnvVarError::TypeError {
                key,
                value,
                expected_type,
                details,
            } => {
                write!(
                    f,
                    "Environment variable '{}' has value '{}' which cannot be converted to {}",
                    key, value, expected_type
                )?;
                if let Some(d) = details {
                    write!(f, ": {}", d)?;
                }
                Ok(())
            }
        }
    }
}

impl std::error::Error for EnvVarError {}

impl From<ParseIntError> for EnvVarError {
    fn from(err: ParseIntError) -> Self {
        EnvVarError::TypeError {
            key: String::new(),
            value: String::new(),
            expected_type: "integer".to_string(),
            details: Some(err.to_string()),
        }
    }
}

impl From<ParseFloatError> for EnvVarError {
    fn from(err: ParseFloatError) -> Self {
        EnvVarError::TypeError {
            key: String::new(),
            value: String::new(),
            expected_type: "float".to_string(),
            details: Some(err.to_string()),
        }
    }
}

impl From<ParseBoolError> for EnvVarError {
    fn from(err: ParseBoolError) -> Self {
        EnvVarError::TypeError {
            key: String::new(),
            value: String::new(),
            expected_type: "bool".to_string(),
            details: Some(err.to_string()),
        }
    }
}

/// Type-safe environment variable accessor
pub struct EnvVar;

impl EnvVar {
    /// Get string environment variable
    ///
    /// # Examples
    ///
    /// ```
    /// use typed_env_vars::EnvVar;
    ///
    /// // Required string
    /// let api_key = EnvVar::string("API_KEY");
    /// ```
    pub fn string(key: &str) -> Result<String, EnvVarError> {
        env::var(key).map_err(|_| EnvVarError::NotFound(key.to_string()))
    }

    /// Get integer environment variable (i32)
    ///
    /// # Examples
    ///
    /// ```
    /// use typed_env_vars::EnvVar;
    ///
    /// let port = EnvVar::int("PORT").unwrap_or(8000);
    /// ```
    pub fn int(key: &str) -> Result<i32, EnvVarError> {
        let raw = Self::string(key)?;
        raw.parse::<i32>().map_err(|_| EnvVarError::TypeError {
            key: key.to_string(),
            value: raw.clone(),
            expected_type: "i32".to_string(),
            details: None,
        })
    }

    /// Get 64-bit integer environment variable
    ///
    /// # Examples
    ///
    /// ```
    /// use typed_env_vars::EnvVar;
    ///
    /// let max_size = EnvVar::int64("MAX_SIZE").unwrap_or(1000000);
    /// ```
    pub fn int64(key: &str) -> Result<i64, EnvVarError> {
        let raw = Self::string(key)?;
        raw.parse::<i64>().map_err(|_| EnvVarError::TypeError {
            key: key.to_string(),
            value: raw.clone(),
            expected_type: "i64".to_string(),
            details: None,
        })
    }

    /// Get float environment variable (f64)
    ///
    /// # Examples
    ///
    /// ```
    /// use typed_env_vars::EnvVar;
    ///
    /// let timeout = EnvVar::float("TIMEOUT").unwrap_or(30.0);
    /// ```
    pub fn float(key: &str) -> Result<f64, EnvVarError> {
        let raw = Self::string(key)?;
        raw.parse::<f64>().map_err(|_| EnvVarError::TypeError {
            key: key.to_string(),
            value: raw.clone(),
            expected_type: "f64".to_string(),
            details: None,
        })
    }

    /// Get boolean environment variable
    ///
    /// Accepts: true/false, yes/no, 1/0, on/off (case-insensitive)
    ///
    /// # Examples
    ///
    /// ```
    /// use typed_env_vars::EnvVar;
    ///
    /// let debug = EnvVar::bool("DEBUG").unwrap_or(false);
    /// ```
    pub fn bool(key: &str) -> Result<bool, EnvVarError> {
        let raw = Self::string(key)?;
        let normalized = raw.to_lowercase().trim().to_string();

        match normalized.as_str() {
            "true" | "yes" | "1" | "on" | "t" | "y" => Ok(true),
            "false" | "no" | "0" | "off" | "f" | "n" => Ok(false),
            _ => Err(EnvVarError::TypeError {
                key: key.to_string(),
                value: raw,
                expected_type: "bool".to_string(),
                details: Some("Valid values: true, false, yes, no, 1, 0, on, off".to_string()),
            }),
        }
    }

    /// Get list environment variable
    ///
    /// # Examples
    ///
    /// ```
    /// use typed_env_vars::EnvVar;
    ///
    /// // Comma-separated by default
    /// let hosts = EnvVar::list("ALLOWED_HOSTS", ",").unwrap_or_default();
    /// ```
    pub fn list(key: &str, separator: &str) -> Result<Vec<String>, EnvVarError> {
        let raw = Self::string(key)?;
        if raw.trim().is_empty() {
            return Ok(Vec::new());
        }

        Ok(raw
            .split(separator)
            .map(|s| s.trim().to_string())
            .filter(|s| !s.is_empty())
            .collect())
    }

    /// Get map environment variable
    ///
    /// Format: KEY1=VALUE1,KEY2=VALUE2
    ///
    /// # Examples
    ///
    /// ```
    /// use typed_env_vars::EnvVar;
    ///
    /// let flags = EnvVar::map("FEATURE_FLAGS", ",", "=").unwrap_or_default();
    /// ```
    pub fn map(
        key: &str,
        item_separator: &str,
        kv_separator: &str,
    ) -> Result<HashMap<String, String>, EnvVarError> {
        let raw = Self::string(key)?;
        if raw.trim().is_empty() {
            return Ok(HashMap::new());
        }

        let mut result = HashMap::new();
        for item in raw.split(item_separator) {
            let trimmed = item.trim();
            if trimmed.is_empty() {
                continue;
            }

            let parts: Vec<&str> = trimmed.splitn(2, kv_separator).collect();
            if parts.len() != 2 {
                return Err(EnvVarError::TypeError {
                    key: key.to_string(),
                    value: raw,
                    expected_type: "map".to_string(),
                    details: Some(format!(
                        "Expected format: KEY1{}VALUE1{}KEY2{}VALUE2",
                        kv_separator, item_separator, kv_separator
                    )),
                });
            }

            result.insert(parts[0].trim().to_string(), parts[1].trim().to_string());
        }

        Ok(result)
    }

    /// Get URL environment variable with validation
    ///
    /// # Examples
    ///
    /// ```
    /// use typed_env_vars::EnvVar;
    ///
    /// let db_url = EnvVar::url("DATABASE_URL");
    /// ```
    pub fn url(key: &str) -> Result<String, EnvVarError> {
        let value = Self::string(key)?;

        if value.is_empty() {
            return Ok(value);
        }

        let valid_schemes = ["http://", "https://", "ftp://", "ws://", "wss://"];
        let has_valid_scheme = valid_schemes.iter().any(|s| value.starts_with(s));

        if !has_valid_scheme {
            return Err(EnvVarError::TypeError {
                key: key.to_string(),
                value,
                expected_type: "URL".to_string(),
                details: Some("URL must start with http://, https://, ftp://, ws://, or wss://".to_string()),
            });
        }

        Ok(value)
    }

    /// Get environment variable with custom converter
    ///
    /// # Examples
    ///
    /// ```
    /// use typed_env_vars::EnvVar;
    ///
    /// let date = EnvVar::custom("RELEASE_DATE", |s| {
    ///     // Custom parsing logic
    ///     Ok(s.to_string())
    /// });
    /// ```
    pub fn custom<T, F>(key: &str, converter: F) -> Result<T, EnvVarError>
    where
        F: FnOnce(&str) -> Result<T, Box<dyn std::error::Error>>,
    {
        let raw = Self::string(key)?;
        converter(&raw).map_err(|e| EnvVarError::TypeError {
            key: key.to_string(),
            value: raw,
            expected_type: "custom".to_string(),
            details: Some(e.to_string()),
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_string() {
        env::set_var("TEST_STRING", "hello");
        assert_eq!(EnvVar::string("TEST_STRING").unwrap(), "hello");
        env::remove_var("TEST_STRING");
    }

    #[test]
    fn test_string_not_found() {
        env::remove_var("MISSING_VAR");
        assert!(matches!(
            EnvVar::string("MISSING_VAR"),
            Err(EnvVarError::NotFound(_))
        ));
    }

    #[test]
    fn test_int() {
        env::set_var("TEST_INT", "42");
        assert_eq!(EnvVar::int("TEST_INT").unwrap(), 42);
        env::remove_var("TEST_INT");
    }

    #[test]
    fn test_int_invalid() {
        env::set_var("TEST_INT", "not_a_number");
        assert!(matches!(
            EnvVar::int("TEST_INT"),
            Err(EnvVarError::TypeError { .. })
        ));
        env::remove_var("TEST_INT");
    }

    #[test]
    fn test_float() {
        env::set_var("TEST_FLOAT", "3.14");
        assert_eq!(EnvVar::float("TEST_FLOAT").unwrap(), 3.14);
        env::remove_var("TEST_FLOAT");
    }

    #[test]
    fn test_bool() {
        let true_values = vec!["true", "True", "TRUE", "yes", "YES", "1", "on", "ON"];
        for val in true_values {
            env::set_var("TEST_BOOL", val);
            assert_eq!(EnvVar::bool("TEST_BOOL").unwrap(), true);
        }

        let false_values = vec!["false", "False", "FALSE", "no", "NO", "0", "off", "OFF"];
        for val in false_values {
            env::set_var("TEST_BOOL", val);
            assert_eq!(EnvVar::bool("TEST_BOOL").unwrap(), false);
        }

        env::remove_var("TEST_BOOL");
    }

    #[test]
    fn test_bool_invalid() {
        env::set_var("TEST_BOOL_INVALID", "maybe");
        assert!(matches!(
            EnvVar::bool("TEST_BOOL_INVALID"),
            Err(EnvVarError::TypeError { .. })
        ));
        env::remove_var("TEST_BOOL_INVALID");
    }

    #[test]
    fn test_list() {
        env::set_var("TEST_LIST_1", "a,b,c");
        let result = EnvVar::list("TEST_LIST_1", ",").unwrap();
        assert_eq!(result, vec!["a", "b", "c"]);
        env::remove_var("TEST_LIST_1");
    }

    #[test]
    fn test_list_with_whitespace() {
        env::set_var("TEST_LIST_2", "a, b , c");
        let result = EnvVar::list("TEST_LIST_2", ",").unwrap();
        assert_eq!(result, vec!["a", "b", "c"]);
        env::remove_var("TEST_LIST_2");
    }

    #[test]
    fn test_map() {
        env::set_var("TEST_MAP", "key1=value1,key2=value2");
        let result = EnvVar::map("TEST_MAP", ",", "=").unwrap();
        assert_eq!(result.get("key1"), Some(&"value1".to_string()));
        assert_eq!(result.get("key2"), Some(&"value2".to_string()));
        env::remove_var("TEST_MAP");
    }

    #[test]
    fn test_url() {
        env::set_var("TEST_URL_VALID", "https://example.com");
        assert_eq!(
            EnvVar::url("TEST_URL_VALID").unwrap(),
            "https://example.com"
        );
        env::remove_var("TEST_URL_VALID");
    }

    #[test]
    fn test_url_invalid() {
        env::set_var("TEST_URL_INVALID", "not-a-url");
        assert!(matches!(
            EnvVar::url("TEST_URL_INVALID"),
            Err(EnvVarError::TypeError { .. })
        ));
        env::remove_var("TEST_URL_INVALID");
    }
}
