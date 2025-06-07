# Database Migrations and Seeds

This project uses a custom migration system built on top of [golang-migrate](https://github.com/golang-migrate/migrate) to manage database schema changes and seed data.

## Overview

The migration system replaces the manual SQL setup from `docs/sql/blog.sql` with automated, versioned database migrations. This provides better control over database schema changes and makes it easier to manage different environments.

## Features

- **Automated migrations**: Runs automatically when the application starts
- **Version control**: Each migration is versioned and tracked
- **Rollback support**: Ability to rollback migrations if needed
- **Seed data**: Automated insertion of initial data
- **CLI tool**: Command-line interface for manual migration management

## Migration Files

Migration files are stored in the `migrations/` directory with the following naming convention:
- `{version}_{description}.up.sql` - Forward migration
- `{version}_{description}.down.sql` - Rollback migration

### Current Migrations

1. **001_create_tables**: Creates the initial database schema
   - `blog_tag` - Article tags management
   - `blog_article` - Article management
   - `blog_auth` - Authentication

2. **002_seed_data**: Inserts initial seed data
   - Default auth user (test/test123)
   - Sample tags (Technology, Programming, etc.)
   - Sample articles

## Usage

### Automatic Migrations

Migrations run automatically when the application starts. This ensures the database is always up-to-date.

### Manual Migration Management

Use the CLI tool for manual migration management:

```bash
# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Check current migration version
make migrate-version

# Migrate to a specific version
make migrate-to VERSION=1

# Reset database (WARNING: drops all data)
make migrate-reset
```

### Creating New Migrations

To create a new migration:

```bash
# Create migration files
make new-migration NAME=add_user_email_column

# This creates:
# migrations/20231205120000_add_user_email_column.up.sql
# migrations/20231205120000_add_user_email_column.down.sql
```

Then edit the generated files:

**up.sql** (forward migration):
```sql
ALTER TABLE blog_auth ADD COLUMN email VARCHAR(255) DEFAULT '';
```

**down.sql** (rollback migration):
```sql
ALTER TABLE blog_auth DROP COLUMN email;
```

## Configuration

The migration system uses the same database configuration as the main application from `conf/app.ini`:

```ini
[database]
Type = mysql
User = root
Password = password
Host = 127.0.0.1:3306
Name = blog
TablePrefix = blog_
```

## Migration Package

The `pkg/migration` package provides programmatic access to migrations:

```go
import "github.com/EDDYCJY/go-gin-example/pkg/migration"

// Run all pending migrations
err := migration.RunMigrations()

// Rollback last migration
err := migration.RollbackMigrations()

// Get current version
version, dirty, err := migration.GetMigrationVersion()

// Migrate to specific version
err := migration.MigrateToVersion(2)
```

## Best Practices

1. **Always create both up and down migrations** - This ensures you can rollback changes if needed
2. **Test migrations thoroughly** - Test both forward and rollback migrations
3. **Keep migrations small** - Break large changes into smaller, manageable migrations
4. **Never modify existing migrations** - Once a migration is deployed, create a new migration for changes
5. **Use descriptive names** - Migration names should clearly describe what they do
6. **Backup before major changes** - Always backup your database before running major migrations

## Troubleshooting

### Migration Failed

If a migration fails, the database might be in a "dirty" state. Check the migration version:

```bash
make migrate-version
```

If the migration is dirty, you may need to manually fix the database and then force the version:

```bash
# After fixing the database manually
make migrate-to VERSION=<correct_version>
```

### Database Connection Issues

Ensure your database configuration in `conf/app.ini` is correct and the database server is running.

### Permission Issues

Make sure the database user has the necessary permissions to create/drop tables and modify schema.

## Migration vs Direct GORM AutoMigrate

This migration system is preferred over GORM's `AutoMigrate()` because:

- **Version control**: Explicit versioning of schema changes
- **Rollback capability**: Ability to undo changes
- **Production safety**: Controlled deployment of schema changes
- **Team collaboration**: Clear history of database changes
- **Environment consistency**: Same schema across all environments

## Files Overview

- `migrations/` - Migration SQL files
- `pkg/migration/migration.go` - Migration package implementation
- `cmd/migrate/main.go` - CLI tool for manual migration management
- `Makefile` - Convenient commands for migration management
- `MIGRATIONS.md` - This documentation file

## Dependencies

- [golang-migrate/migrate](https://github.com/golang-migrate/migrate) - Core migration library
- MySQL driver for migrations
- File source driver for reading migration files