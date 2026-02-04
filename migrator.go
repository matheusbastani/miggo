package miggo

import "database/sql"

// Migrator encapsulates the migration configuration and provides
// convenient methods for running migrations without repeatedly
// passing the same parameters.
type Migrator struct {
	db      *sql.DB
	baseDir string
}

// New creates a new Migrator instance with the given database connection
// and base directory for migrations.
//
// Parameters:
//   - db: database connection
//   - baseDir: base directory containing migration folders
//
// Example:
//
//	m := miggo.New(db, "./migrations")
//	m.Up()
//	m.Create("add_users_table")
func New(db *sql.DB, baseDir string) *Migrator {
	return &Migrator{
		db:      db,
		baseDir: baseDir,
	}
}

// Up applies all pending migrations.
func (m *Migrator) Up() {
	Up(m.db, m.baseDir)
}

// Down rolls back the most recently applied migration.
func (m *Migrator) Down() {
	Down(m.db, m.baseDir)
}

// Reset rolls back all applied migrations.
func (m *Migrator) Reset() {
	Reset(m.db, m.baseDir)
}

// ResetAndDrop rolls back all migrations and drops the migrations table.
func (m *Migrator) ResetAndDrop() {
	ResetAndDrop(m.db, m.baseDir)
}

// Version displays the latest applied migration.
func (m *Migrator) Version() {
	Version(m.db)
}

// Create creates a new migration with the given name.
// Optionally accepts an index number.
func (m *Migrator) Create(name string, index ...int) {
	Create(m.baseDir, name, index...)
}

// Insert creates a new migration at a specific index,
// renumbering existing migrations as needed.
func (m *Migrator) Insert(name string, insertIndex int) {
	Insert(m.baseDir, name, insertIndex)
}
