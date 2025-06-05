#!/bin/bash

# Script to create new migration files
# Usage: ./scripts/new-migration.sh migration_name

if [ -z "$1" ]; then
    echo "Usage: ./scripts/new-migration.sh <migration_name>"
    echo "Example: ./scripts/new-migration.sh add_user_email_column"
    exit 1
fi

MIGRATION_NAME=$1
TIMESTAMP=$(date +%Y%m%d%H%M%S)
UP_FILE="migrations/${TIMESTAMP}_${MIGRATION_NAME}.up.sql"
DOWN_FILE="migrations/${TIMESTAMP}_${MIGRATION_NAME}.down.sql"

# Create migration files
touch "$UP_FILE"
touch "$DOWN_FILE"

# Add basic templates
cat > "$UP_FILE" << EOF
-- Migration: $MIGRATION_NAME
-- Created at: $(date)

-- Add your forward migration SQL here
-- Example:
-- ALTER TABLE blog_users ADD COLUMN email VARCHAR(255) DEFAULT '';
EOF

cat > "$DOWN_FILE" << EOF
-- Rollback migration: $MIGRATION_NAME
-- Created at: $(date)

-- Add your rollback migration SQL here
-- Example:
-- ALTER TABLE blog_users DROP COLUMN email;
EOF

echo "Created migration files:"
echo "  $UP_FILE"
echo "  $DOWN_FILE"
echo ""
echo "Edit these files with your migration SQL, then run:"
echo "  make migrate-up"