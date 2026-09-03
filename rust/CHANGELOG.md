# Changelog

All notable changes to the Rust implementation of `typed-env-vars` are documented in this file.

## [1.1.0] - 2026-09-03

### Added

- `EnvVar::dict` alias for `EnvVar::map`.
- `EnvVar::enum_value` for type-safe enum lookup with optional defaults.

### Changed

- `EnvVar::list` and `EnvVar::map` defaults are returned as copies instead of being serialized and parsed back.

### Fixed

- URL validation no longer maintains a hard-coded scheme allow-list; it now parses `scheme://host` generically and accepts database and message-broker URLs.

## [1.0.0] - 2024-01-15

### Added

- Initial release.
- Support for `String`, `i32`, `i64`, `f64`, `bool`, `Vec`, `HashMap`, `URL`, and `custom` types.
