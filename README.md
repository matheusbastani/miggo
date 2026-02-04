# miggo

A simple, flexible SQL migration library for Go with zero dependencies on ORMs or frameworks.

## Features

- 📁 **Directory-based migrations** - Each migration is a folder with `.up.sql` and `.down.sql` files
- 🔢 **Sequential numbering** - Automatic migration indexing with `001_`, `002_`, etc.
- ⚡ **Transaction safety** - Each migration runs in a transaction
- 🔄 **Full rollback support** - Down migrations and reset functionality
- 🎯 **Insert at any position** - Add migrations between existing ones with automatic renumbering
- 🗄️ **PostgreSQL ready** - Works with any `database/sql` compatible driver
- 🎨 **Colored output** - Clear visual feedback with color-coded messages

## Installation

```bash
go get github.com/matheusbastani/miggo
```

## Quick Start

```go
package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/matheusbastani/miggo"
)

func main() {
    db, _ := sql.Open("postgres", "postgres://user:pass@localhost/dbname?sslmode=disable")
    defer db.Close()

    // Create a migrator instance
    m := miggo.New(db, "./migrations")

    // Create a new migration
    m.Create("create_users_table")

    // Apply all pending migrations
    m.Up()

    // Check current version
    m.Version()
}
```

## Migration Structure

Migrations are organized in numbered directories:

```
migrations/
├── 001_create_users/
│   ├── 20240204120000_create_users.up.sql
│   └── 20240204120000_create_users.down.sql
├── 002_add_posts/
│   ├── 20240204130000_add_posts.up.sql
│   └── 20240204130000_add_posts.down.sql
└── 003_add_comments/
    ├── 20240204140000_add_comments.up.sql
    └── 20240204140000_add_comments.down.sql
```

### Methods

#### Up()
```go
m.Up()
```
Applies all pending migrations in sequential order. Creates the `migrations` tracking table if it doesn't exist.

#### Down()
```go
m.Down()
```
Rolls back the most recently applied migration.

#### Reset()
```go
m.Reset()
```
Rolls back all applied migrations in reverse order.

#### ResetAndDrop()
```go
m.ResetAndDrop()
```
Rolls back all migrations and drops the `migrations` tracking table.

#### Version()
```go
m.Version()
```
Displays the latest applied migration.

#### Create(name string, index ...int)
```go
// Create next migration automatically
m.Create("add_email_verification")

// Create migration with specific index
m.Create("add_email_verification", 5)
```
Creates a new migration directory with `.up.sql` and `.down.sql` files.

#### Insert(name string, insertIndex int)
```go
m.Insert("add_missing_index", 3)
```
Creates a new migration at a specific position, automatically renumbering existing migrations.

## Usage Examples

### Basic Workflow

```go
m := miggo.New(db, "./migrations")

// 1. Create migrations
m.Create("create_users")
m.Create("create_posts")
m.Create("create_comments")

// 2. Write your SQL in the generated files
// Edit: migrations/001_create_users/*.up.sql and *.down.sql

// 3. Apply migrations
m.Up()

// 4. Check status
m.Version()
// Output: latest migration: 003_create_comments/20240204140000_create_comments.up.sql
```

### Using Standalone Functions

If you prefer not to use the `Migrator` type:

```go
import "github.com/matheusbastani/miggo"

// Create migration
miggo.Create("./migrations", "add_users", 1)

// Apply migrations
miggo.Up(db, "./migrations")

// Rollback
miggo.Down(db, "./migrations")

// Check version
miggo.Version(db)

// Reset all
miggo.Reset(db, "./migrations")
miggo.ResetAndDrop(db, "./migrations")

// Insert migration
miggo.Insert("./migrations", "add_index", 5)
```

## Database Support

Currently optimized for PostgreSQL, but should work with any `database/sql` compatible driver that supports:
- `CREATE TABLE IF NOT EXISTS`
- `UUID` type (or you can modify to use `TEXT`)
- Transactions

### Using with Other Databases

For MySQL/MariaDB, you may need to adjust the migrations table schema:

```sql
CREATE TABLE IF NOT EXISTS migrations (
    id VARCHAR(36) PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Best Practices

1. **Always include DOWN migrations** - Make your migrations reversible
2. **One logical change per migration** - Keep migrations focused and atomic
3. **Test migrations on a copy** - Validate both up and down migrations
4. **Don't modify applied migrations** - Create a new migration to fix issues
5. **Use transactions** - miggo does this automatically, but be aware
6. **Backup before production** - Always backup before running migrations in production

## Error Handling

miggo uses colored output for clear error reporting:

- 🟢 **Green**: Successful operations
- 🟡 **Yellow**: Warnings or empty states
- 🔴 **Red**: Errors (also exits with status code 1)

All errors are printed to stdout with context about what failed.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Uses [fatih/color](https://github.com/fatih/color) for colored output
- Uses [google/uuid](https://github.com/google/uuid) for UUID generation