# Migration from Mixed Seeds to Environment-Specific Seeding

This document explains how we've implemented environment-specific seeding and how to migrate from the old system.

## What Changed

### Before (Old System)
- Seeds were mixed with migrations in the `migrations/` directory
- Same seed data was used for all environments
- Files like `4_seed_data.up.sql` and `5_seed_tags.up.sql` contained test data

### After (New System)
- Seeds are separated into environment-specific directories under `seeds/`
- Different data for development, production, and staging
- Clean separation between schema migrations and data seeding

## Directory Structure

```
seeds/
├── development/          # Rich test data for developers
│   ├── 001_seed_auth.up.sql
│   ├── 001_seed_auth.down.sql
│   ├── 002_seed_tags.up.sql
│   └── 002_seed_tags.down.sql
├── production/           # Minimal essential data
│   ├── 001_seed_auth.up.sql
│   ├── 001_seed_auth.down.sql
│   ├── 002_seed_tags.up.sql
│   └── 002_seed_tags.down.sql
└── staging/              # Production-like but safe for testing
    ├── 001_seed_auth.up.sql
    └── 001_seed_auth.down.sql
```

## Migration Steps

### 1. Remove Old Seed Migrations (Optional)
The old seed files in `migrations/` directory can be kept for reference but should not be used:
- `migrations/4_seed_data.up.sql` → Replaced by environment-specific auth seeds
- `migrations/5_seed_tags.up.sql` → Replaced by environment-specific tag seeds

### 2. Use New Seeding Commands
Replace old migration commands with new seeding commands:

**Old way:**
```bash
make migrate-up  # This would run both schema and seeds
```

**New way:**
```bash
make migrate-up  # Only runs schema migrations
make seed-dev    # Runs development seeds
# OR
make setup-dev   # Runs both migrations and dev seeds
```

## Environment Detection

The system automatically detects your environment:

1. **APP_ENV environment variable** (highest priority)
2. **RunMode from config** (fallback)
   - `debug` → `development`
   - `release` → `production`

```bash
# Set environment explicitly
export APP_ENV=development
export APP_ENV=production
export APP_ENV=staging

# Or use config file setting
[server]
RunMode = debug    # Maps to development
RunMode = release  # Maps to production
```

## Available Commands

### Make Commands (Recommended)
```bash
# Seeding
make seed-dev          # Run development seeds
make seed-prod         # Run production seeds
make seed-status       # Check status for all environments
make seed-status ENV=development  # Check specific environment
make seed-rollback ENV=development
make seed-list         # List available environments

# Full setup
make setup-dev         # Migrations + development seeds
make setup-prod        # Migrations + production seeds

# Creating new seeds
make new-seed ENV=development NAME=create_sample_articles
```

### Direct Commands
```bash
# Using seed command
go run cmd/seed/main.go -action=run -env=development
go run cmd/seed/main.go -action=run -env=production
go run cmd/seed/main.go -action=status -env=development
go run cmd/seed/main.go -action=list

# Using migrate command
go run cmd/migrate/main.go -action=seed -env=development
go run cmd/migrate/main.go -action=seed-status -env=development
```

## Data Differences by Environment

### Development Seeds
- Multiple test users (admin, testuser, developer)
- Rich tag data (Technology, Programming, Go, Web Development, API, Testing)
- Sample articles and test data
- Edge cases and boundary data

### Production Seeds
- Single admin user
- Minimal essential tags (General)
- No test data
- Only data required for application to function

### Staging Seeds
- Production-like data but safe for testing
- Staging-specific credentials
- Reduced dataset for performance

## Benefits

1. **Environment Safety**: No test data in production
2. **Development Efficiency**: Rich test data for developers
3. **Deployment Flexibility**: Different data per environment
4. **Data Consistency**: Predictable seed data per environment
5. **Easy Management**: Clear separation and commands

## Best Practices

### 1. Keep Production Seeds Minimal
```sql
-- Good: Only essential data
INSERT INTO `blog_auth` (`username`, `password`) VALUES ('admin', 'secure_hash');

-- Avoid: Test data in production
INSERT INTO `blog_auth` (`username`, `password`) VALUES 
('admin', 'secure_hash'),
('testuser', 'test123'),        -- Don't do this in production
('developer', 'dev123');        -- Don't do this in production
```

### 2. Make Seeds Idempotent
```sql
-- Good: Safe to run multiple times
INSERT IGNORE INTO `blog_auth` (`id`, `username`, `password`) 
VALUES (1, 'admin', 'admin123');

-- Or use conditional insert
INSERT INTO `blog_auth` (`username`, `password`) 
SELECT 'admin', 'admin123'
WHERE NOT EXISTS (SELECT 1 FROM `blog_auth` WHERE `username` = 'admin');
```

### 3. Use Secure Passwords in Production
```sql
-- Development: Simple passwords for testing
INSERT INTO `blog_auth` (`username`, `password`) VALUES ('admin', 'admin123');

-- Production: Secure hashed passwords
INSERT INTO `blog_auth` (`username`, `password`) VALUES 
('admin', '$2a$10$encrypted_password_hash_here');
```

## Troubleshooting

### Common Issues

1. **"No seed files found"**
   - Ensure environment directory exists: `seeds/development/`
   - Check file naming: `001_seed_name.up.sql`

2. **Wrong environment detected**
   - Set `APP_ENV` environment variable
   - Check `RunMode` in `conf/app.ini`

3. **Seeds not applying**
   - Check database connection
   - Verify SQL syntax in seed files
   - Check seed status: `make seed-status`

### Debugging
```bash
# Check current environment
go run cmd/seed/main.go -action=list

# Check what seeds are available
ls -la seeds/development/

# Check database connectivity
make migrate-version
```

## Migration Checklist

- [ ] Review old seed files in `migrations/` directory
- [ ] Create environment-specific seed files in `seeds/` directories
- [ ] Test development seeds: `make seed-dev`
- [ ] Test production seeds: `make seed-prod`
- [ ] Update deployment scripts to use new commands
- [ ] Update documentation for team members
- [ ] Set `APP_ENV` in production deployment
- [ ] Verify environment detection works correctly

This new system provides better separation of concerns, environment safety, and easier management of database seeding across different deployment environments.