## 1.3.0

### Added

- `default` option in `miggo.yaml` to mark a database as the default
- `--db` is now optional: miggo uses the only configured database, or the one marked as `default`

### Changed

- When multiple databases are configured and none is marked as default, the error now lists the available databases

## 1.2.0

### Breaking Changes

- database name is now passed via --db / -d flag instead of a positional argument.
  - Before: miggo create users_table table_test
  - After: miggo create users_table --db table_test

### Added

- Added SQLite support alongside Postgres.

### Fixed

- Fixed a bug where reset-drop could fail after a reset due to the connection being closed prematurely.

## 1.1.1

### Added

- Added support for expanding environment variables across all YAML configuration fields
- Added support for loading .env.local by default before falling back to .env

### Changed

- Improved configuration flexibility by allowing any setting to use ${ENV_VARIABLE} placeholders

## 1.1.0

### Breaking Changes

- Migrated from a Go library to a standalone CLI tool
- Added miggo.yaml configuration with multi-database support
- Replaced secure with environment-based safety (development / production)

### Added

- Added migration rollback boundaries with lock and unlock
- Added interactive shell mode
- Added new CLI commands: init, create, up, down, reset, reset-drop, insert

## 1.0.1

### Changed

- Replaced migrations table name to `schema_migrations`

### 1.0.0

### Added

- Initial release
