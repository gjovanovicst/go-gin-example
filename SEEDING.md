# Environment-Specific Database Seeding

This document explains how to use the environment-specific database seeding system in the Go Gin Example project.

## Overview

The seeding system allows you to have different seed data for different environments (development, production, staging, etc.). This separation ensures that:

- **Development**: Has rich test data for developers to work with
- **Production**: Has minimal, essential data only
- **Testing**: Can have specific test fixtures
- **Staging**: Can mirror production or have staging-specific data

## Directory Structure

```
seeds/
├── development/
│   ├── 001_seed_auth.up.sql
│   ├── 001_seed_auth.down.sql
│   ├── 002_seed_tags.up.sql
│   └── 002_seed_tags.down.sql
├── production/
│   ├── 001_seed_auth.up.sql
│   ├── 001_seed_auth.down.sql
│   ├── 002_seed_tags.up.sql
│   └── 002_seed_tags.down.sql
└── staging/           # Optional
    ├── 001_seed_auth.up.sql
    └── ...
```

## Environment Detection

The system automatically detects the environment using:

1. **APP_ENV environment variable** (highest priority)
2. **RunMode from config file** (fallback)
   - `debug` → `development`
   - `release` → `production`

### Setting Environment

```bash
# Method 1: Environment variable
export APP_ENV=development
export APP_ENV=production

# Method 2: Config file (conf/app.ini)
[server]
RunMode = debug    # Maps to development
RunMode = release  # Maps to production
```

## Usage

### Using Make Commands (Recommended)

```bash
# Run development seeds
make seed-dev

# Run production seeds
make seed-prod

# Check seed status for specific environment
make seed-status ENV=development

# Check seed status for all environments
make seed-status

# Rollback seeds for specific environment
make seed-rollback ENV=development

# List available environments
make seed-list

# Full environment setup (migrations + seeds)
make setup-dev    # Runs migrations + development seeds
make setup-prod   # Runs migrations + production seeds
```

### Using Direct Commands

```bash
# Using the seed command
go run cmd/seed/main.go -action=run -env=development
go run cmd/seed/main.go -action=run -env=production
go run cmd/seed/main.go -action=status -env=development
go run cmd/seed/main.go -action=rollback -env=development
go run cmd/seed/main.go -action=list

# Using the migrate command
go run cmd/migrate/main.go -action=seed -env=development
go run cmd/migrate/main.go -action=seed-status -env=development
```

### Manual Seeding

For direct SQL execution without migration tracking:

```bash
go run cmd/seed/main.go -action=run -env=development -manual
```

## Creating Seed Files

### Naming Convention

Seed files follow the same naming pattern as migrations:
- `{sequence}_{description}.up.sql` - For applying seeds
- `{sequence}_{description}.down.sql` - For rolling back seeds

### Example Development Seeds

**seeds/development/001_seed_auth.up.sql:**
```sql
-- Development seed data for auth table
INSERT INTO `blog_auth` (`id`, `username`, `password`) VALUES 
(1, 'admin', 'admin123'),
(2, 'testuser', 'test123'),
(3, 'developer', 'dev123');
```

**seeds/development/001_seed_auth.down.sql:**
```sql
-- Rollback development seed data for auth table
DELETE FROM `blog_auth` WHERE `id` IN (1, 2, 3);
```

### Example Production Seeds

**seeds/production/001_seed_auth.up.sql:**
```sql
-- Production seed data for auth table - minimal essential data
INSERT INTO `blog_auth` (`id`, `username`, `password`) VALUES 
(1, 'admin', 'admin123');
```

**seeds/production/001_seed_auth.down.sql:**
```sql
-- Rollback production seed data for auth table
DELETE FROM `blog_auth` WHERE `id` = 1;
```

## Best Practices

### 1. Keep Production Seeds Minimal
- Only include essential data required for the application to function
- Avoid test users, sample content, or development-specific data

### 2. Rich Development Seeds
- Include comprehensive test data
- Multiple users with different roles
- Sample content for testing UI/UX
- Edge cases and boundary data

### 3. Idempotent Seeds
Make seeds safe to run multiple times:

```sql
-- Good: Use INSERT IGNORE or ON DUPLICATE KEY UPDATE
INSERT IGNORE INTO `blog_auth` (`id`, `username`, `password`) VALUES (1, 'admin', 'admin123');

-- Or check existence first
INSERT INTO `blog_auth` (`id`, `username`, `password`) 
SELECT 1, 'admin', 'admin123'
WHERE NOT EXISTS (SELECT 1 FROM `blog_auth` WHERE `id` = 1);
```

### 4. Environment-Specific Data
- Development: Rich test data, multiple users, sample content
- Production: Admin user, essential categories, minimal config
- Staging: Production-like data but safe for testing

### 5. Secure Passwords
For production:
```sql
-- Use secure, hashed passwords
INSERT INTO `blog_auth` (`username`, `password`) VALUES 
('admin', '$2a$10$encrypted_password_hash');
```

## Migration vs Seeding

| Aspect | Migrations | Seeds |
|--------|------------|-------|
| **Purpose** | Schema changes | Initial data |
| **Environment** | Same across all | Different per environment |
| **Versioning** | Sequential, required | Environment-specific |
| **Rollback** | Critical for schema | Optional for data |
| **Frequency** | Every schema change | Once per environment setup |

## Troubleshooting

### Common Issues

1. **"No seed files found"**
   - Check if the environment directory exists in `seeds/`
   - Verify file naming convention

2. **"Failed to connect to database"**
   - Check database configuration in `conf/app.ini`
   - Ensure database server is running

3. **"Environment not detected"**
   - Set `APP_ENV` environment variable
   - Check `RunMode` in config file

### Debugging

```bash
# Check current environment
go run cmd/seed/main.go -action=list

# Check seed status
make seed-status

# Check migration status
make migrate-version
```

## Docker Integration

For Docker deployments:

```dockerfile
# Set environment in Dockerfile
ENV APP_ENV=production

# Or pass via docker-compose
environment:
  - APP_ENV=production
```

## CI/CD Integration

```yaml
# Example GitHub Actions
- name: Setup Development Environment
  run: |
    make migrate-up
    make seed-dev
  if: github.ref == 'refs/heads/develop'

- name: Setup Production Environment
  run: |
    make migrate-up
    make seed-prod
  if: github.ref == 'refs/heads/main'
```

## Advanced Usage

### Creating New Environments

1. Create directory: `seeds/staging/`
2. Add seed files following naming convention
3. Use: `go run cmd/seed/main.go -action=run -env=staging`

### Conditional Seeding

You can add logic to seeds:

```sql
-- Only insert if in development
INSERT INTO `blog_auth` (`username`, `password`) 
SELECT 'testuser', 'test123'
WHERE (SELECT COUNT(*) FROM `blog_auth`) < 5;
```

This environment-specific seeding system provides flexibility while maintaining data integrity across different deployment environments.