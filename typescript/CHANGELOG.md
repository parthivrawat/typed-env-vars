# Changelog

All notable changes to the TypeScript implementation of `typed-env-vars` are documented in this file.

## [1.1.0] - 2026-09-03

### Added

- `env` is exported as a ready-to-use singleton instance of the `Env` class, while `Env` remains available for custom instances.

### Changed

- `env.list` and `env.dict` defaults are cloned (`[...default]` / `{ ...default }`) to avoid accidental mutation of shared state.

### Fixed

- URL validation accepts any valid scheme, including database URLs like `postgresql://` and `redis://`.

## [1.0.0] - 2024-01-15

### Added

- Initial release.
- Support for `string`, `number`, `boolean`, `list`, `dict`, `enum`, `URL`, and `custom` types.
