# Typed Environment Variables

A cross-language, type-safe environment variable library with automatic type conversion and validation.

## Overview

This library provides a consistent, cross-language API for loading and validating environment variables across Python, TypeScript, Go, and Rust. It supports strings, integers, floats, booleans, lists, maps/dicts, enums, URLs, and custom conversions, with optional default values and clear error messages.

## Languages

| Language | Package | README |
|---|---|---|
| **Go** | `github.com/parthivrawat/typed-env-vars/go` | [Go README](go/README.md) |
| **Python** | `pip install typed-env-vars` | [Python README](python/README.md) |
| **Rust** | `cargo add typed-env-vars` | [Rust README](rust/README.md) |
| **TypeScript** | `npm install typed-env-vars` | [TypeScript README](typescript/README.md) |

## Repository Layout

```
typed-env-vars/
├── go/            # Go module
├── python/        # Python package
├── rust/          # Rust crate
├── typescript/    # TypeScript/npm package
└── README.md      # This file
```

## License

MIT License. See the root [LICENSE](LICENSE) file for details. Each language directory also contains a copy of the license.
