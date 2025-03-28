# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-15

### Added

- Core CLI application with Cobra command framework
- Tag management: add, remove, and list tags for files
- File search with boolean tag queries
- SQLite-backed persistence for tags and file associations
- Import and export functionality for tag data
- Configuration management with support for log levels and database paths
- Structured logging with severity levels (debug, info, warn, error)
- Database initialization and migration utilities
- Input validation and error handling utilities
- Repository layer for data access abstraction
- Service layer for business logic

### Fixed

- Initial release with stable API and CLI interface
