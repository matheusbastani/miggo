package errs

import "errors"

var (
	ErrSecureModeEnabled      = errors.New("operation disabled: secure mode is enabled")
	ErrRollbackBoundaryExists = errors.New("operation blocked: a migration is locked, use --force")
	ErrMigrationNotFound      = errors.New("migration not found")
)
