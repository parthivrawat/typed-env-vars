# Changelog

All notable changes to the Python implementation of `typed-env-vars` are documented in this file.

## [1.1.0] - 2026-09-03

### Changed

- `pyproject.toml` now contains full project metadata and dev dependencies; `setup.py` defers to it.
- `env.list` and `env.dict` defaults are returned as shallow copies instead of being shared with the caller.

### Fixed

- URL validation no longer rejects database URLs; it now uses `urllib.parse.urlparse` and accepts any valid scheme such as `postgresql://`, `redis://`, or `amqp://`.

## [1.0.0] - 2024-01-15

### Added

- Initial release.
- Support for `str`, `int`, `float`, `bool`, `list`, `dict`, `enum`, `URL`, and `custom` types.
