# miggo

A simple, flexible SQL migration CLI tool for Go projects.

miggo manages SQL migrations using a directory-based structure, with support for applying migrations, rolling back changes, locking migrations and managing multiple databases through a single configuration file.

## Features

- 📁 **Directory-based migrations** - Each migration is stored in its own folder with `.up.sql` and `.down.sql` files
- 🔢 **Sequential numbering** - Automatic migration indexing (`001_`, `002_`, etc.)
- ⚡ **Transaction safety** - Each migration runs inside a database transaction
- 🔄 **Rollback support** - Safely rollback migrations using down files
- 🔒 **Rollback boundaries** - Lock migrations to prevent accidental rollbacks
- 🗄️ **Multiple databases** - Configure and manage multiple databases from one file
- ⚙️ **YAML configuration** - Simple `miggo.yaml` project configuration
- 🎨 **Colored output** - Clear CLI feedback with status messages
- 💻 **Interactive shell** - Run miggo commands through an interactive console

## Installation

Install miggo using Go:

```bash
go install github.com/matheusbastani/miggo@latest
```

Verify installation:

```bash
miggo
```

---

## Getting Started

Initialize miggo in your project:

```bash
miggo init
```

This creates a `miggo.yaml` file in the current directory:

```yaml
databases:
  development:
    driver: postgres
    url: postgres://user:password@localhost/database?sslmode=disable
    path: ./migrations
    environment: development
```

SQLite is also supported:

```yaml
databases:
  development:
    driver: sqlite
    url: ./dev.db
    path: ./migrations
    environment: development
```

### Configuration

Each database entry contains:

| Field         | Description                                                                            |
| ------------- | -------------------------------------------------------------------------------------- |
| `driver`      | Database driver name (`postgres` or `sqlite`)                                          |
| `url`         | Database connection URL                                                                |
| `path`        | Migration directory                                                                    |
| `environment` | Defines the execution environment and enables safety rules for production environments |

### Environments

miggo recognizes the following environments:

| Environment           | Behavior                                                                     |
| --------------------- | ---------------------------------------------------------------------------- |
| `dev` / `development` | Development mode. Allows normal migration operations                         |
| `prod` / `production` | Production mode. Enables additional safety checks for destructive operations |

Example production configuration:

```yaml
databases:
  production:
    driver: postgres
    url: postgres://user:password@production/database
    path: ./migrations
    environment: production
```

When `environment` is set to `prod` or `production`, miggo automatically enables secure behavior for destructive commands such as `reset` and `reset-drop`.

# Commands

Every command that operates on a database takes a `--db` / `-d` flag specifying which database from `miggo.yaml` to use.

## Create a migration

Creates a new migration folder with `.up.sql` and `.down.sql` files.

```bash
miggo create create_users --db development
```

Example output:

```
migrations/
└── 001_create_users/
    ├── 20260719120000_create_users.up.sql
    └── 20260719120000_create_users.down.sql
```

---

## Apply migrations

Apply all pending migrations:

```bash
miggo up --db development
```

miggo automatically creates the migration tracking table when needed.

---

## Show current version

Display the latest applied migration:

```bash
miggo version --db development
```

Example:

```
latest migration:
003_add_indexes
```

---

## Rollback latest migration

Rollback the most recently applied migration:

```bash
miggo down --db development
```

---

## Lock a migration

Create a rollback boundary.

Locked migrations cannot be rolled back.

```bash
miggo lock 005 --db development
```

Example:

```
001
002
003
004
005 🔒
006
007
```

A rollback will stop at migration `005`.

---

## Unlock a migration

Remove a rollback boundary.

The migration index must always be provided.

```bash
miggo unlock 005 --db development
```

---

## Reset migrations

Rollback all migrations:

```bash
miggo reset --db development
```

For destructive environments:

```bash
miggo reset --db development --force
```

---

## Reset and drop migration table

Rollback all migrations and remove miggo's tracking table:

```bash
miggo reset-drop --db development
```

Force mode:

```bash
miggo reset-drop --db development --force
```

---

## Insert migration

Create a migration at a specific index.

```bash
miggo insert add_email_verification 3 --db development
```

Existing migrations are automatically renumbered.

Before:

```
001_create_users
002_create_posts
003_create_comments
```

After:

```
001_create_users
002_create_posts
003_add_email_verification
004_create_comments
```

---

# Interactive Shell

Running miggo without arguments starts the interactive shell:

```bash
miggo
```

Example:

```
  __  __ _
 |  \/  (_)__ _ __ _ ___
 | |\/| | / _` / _` / _ \
 |_|  |_|_\__, \__, \___/
          |___/|___/

miggo>
```

Commands can be executed directly, using the same flags as the CLI:

```
miggo> up --db development
miggo> version --db development
miggo> down --db development
```

Exit:

```
miggo> exit
```

---

# Migration Structure

A typical project:

```
project/
├── miggo.yaml
└── migrations/
    ├── 001_create_users/
    │   ├── 20260719120000_create_users.up.sql
    │   └── 20260719120000_create_users.down.sql
    │
    ├── 002_create_posts/
    │   ├── 20260719130000_create_posts.up.sql
    │   └── 20260719130000_create_posts.down.sql
    │
    └── 003_add_indexes/
        ├── 20260719140000_add_indexes.up.sql
        └── 20260719140000_add_indexes.down.sql
```

Each migration contains two files.

## Up migration

Executed when applying migrations:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL
);
```

## Down migration

Executed during rollback:

```sql
DROP TABLE users;
```

---

# Database Support

miggo currently supports the following drivers:

- PostgreSQL
- SQLite

Support for additional `database/sql` compatible drivers may be added in the future.

---

# Best Practices

1. Always create `.down.sql` files
2. Keep migrations small and focused
3. Never modify migrations already applied in production
4. Create new migrations for changes
5. Test both up and down migrations
6. Use rollback boundaries for important releases
7. Backup databases before destructive operations

---

# Development

Install development tools:

```bash
go install github.com/evilmartians/lefthook@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/tools/gopls@latest

lefthook install
```

---

# Contributing

Contributions are welcome.

Feel free to open issues or submit pull requests.

---

# License

MIT License
