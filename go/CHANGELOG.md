# Changelog

All notable changes to the Go implementation of `typed-env-vars` are documented in this file.

## [1.1.0] - 2026-09-03

### Added

- Error-returning `*E` variants (`StringE`, `IntE`, `FloatE`, `BoolE`, `ListE`, `MapE`, `URLE`, `CustomE`, `EnumE`) for idiomatic `(T, error)` handling.
- `Must*` panic-on-failure helpers for every getter.
- Generic enum support via `env.Enum` and `env.EnumE`.
- Cross-language vocabulary aliases `env.Str`/`env.StrE` and `env.Dict`/`env.DictE`.

### Changed

- `env.List` and `env.Map` defaults are returned as copies instead of being serialized to a string and parsed back, preserving values that contain separators.

### Fixed

- URL validation now uses `url.Parse` and accepts any valid scheme with a host (for example, `postgresql://`, `redis://`, or `amqp://`).

## [1.0.0] - 2024-01-15

### Added

- Initial release.
- Support for `string`, `int`, `float64`, `bool`, `list`, `map`, `URL`, and `custom` types.
