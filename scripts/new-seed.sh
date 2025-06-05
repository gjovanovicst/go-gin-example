#!/bin/bash

# Script to create new seed files for environment-specific seeding
# Usage: ./scripts/new-seed.sh <environment> <name>
# Example: ./scripts/new-seed.sh development create_sample_articles

set -e

# Check if required arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <environment> <name>"
    echo "Example: $0 development create_sample_articles"
    echo "Example: $0 production create_admin_user"
    echo ""
    echo "Available environments:"
    if [ -d "seeds" ]; then
        ls -1 seeds/ 2>/dev/null || echo "  No environments found in seeds/ directory"
    else
        echo "  No seeds/ directory found"
    fi
    exit 1
fi

ENVIRONMENT="$1"
NAME="$2"

# Validate environment name
if [[ ! "$ENVIRONMENT" =~ ^[a-zA-Z0-9_-]+$ ]]; then
    echo "Error: Environment name can only contain letters, numbers, underscores, and hyphens"
    exit 1
fi

# Validate seed name
if [[ ! "$NAME" =~ ^[a-zA-Z0-9_-]+$ ]]; then
    echo "Error: Seed name can only contain letters, numbers, underscores, and hyphens"
    exit 1
fi

# Create seeds directory structure if it doesn't exist
SEEDS_DIR="seeds/$ENVIRONMENT"
mkdir -p "$SEEDS_DIR"

# Find the next sequence number
NEXT_SEQ=1
if [ -d "$SEEDS_DIR" ]; then
    # Find the highest existing sequence number
    LAST_SEQ=$(find "$SEEDS_DIR" -name "*.up.sql" | sed 's/.*\/\([0-9]\+\)_.*/\1/' | sort -n | tail -1 2>/dev/null || echo "0")
    if [ -n "$LAST_SEQ" ] && [ "$LAST_SEQ" -gt 0 ]; then
        NEXT_SEQ=$((LAST_SEQ + 1))
    fi
fi

# Format sequence number with leading zeros
FORMATTED_SEQ=$(printf "%03d" "$NEXT_SEQ")

# Create file paths
UP_FILE="$SEEDS_DIR/${FORMATTED_SEQ}_${NAME}.up.sql"
DOWN_FILE="$SEEDS_DIR/${FORMATTED_SEQ}_${NAME}.down.sql"

# Check if files already exist
if [ -f "$UP_FILE" ]; then
    echo "Error: File $UP_FILE already exists"
    exit 1
fi

if [ -f "$DOWN_FILE" ]; then
    echo "Error: File $DOWN_FILE already exists"
    exit 1
fi

# Create UP file with template
cat > "$UP_FILE" << EOF
-- Seed: $NAME for $ENVIRONMENT environment
-- Created: $(date '+%Y-%m-%d %H:%M:%S')

-- TODO: Add your seed data here
-- Example:
-- INSERT INTO \`blog_table\` (\`column1\`, \`column2\`) VALUES
-- ('value1', 'value2'),
-- ('value3', 'value4');

EOF

# Create DOWN file with template
cat > "$DOWN_FILE" << EOF
-- Rollback seed: $NAME for $ENVIRONMENT environment
-- Created: $(date '+%Y-%m-%d %H:%M:%S')

-- TODO: Add rollback statements here
-- Example:
-- DELETE FROM \`blog_table\` WHERE \`column1\` IN ('value1', 'value3');

EOF

echo "Created seed files for environment '$ENVIRONMENT':"
echo "  UP:   $UP_FILE"
echo "  DOWN: $DOWN_FILE"
echo ""
echo "Next steps:"
echo "1. Edit the files to add your seed data"
echo "2. Run the seeds:"
echo "   make seed-status ENV=$ENVIRONMENT"
echo "   go run cmd/seed/main.go -action=run -env=$ENVIRONMENT"
echo ""
echo "Available commands:"
echo "  go run cmd/seed/main.go -action=run -env=$ENVIRONMENT      # Run seeds"
echo "  go run cmd/seed/main.go -action=rollback -env=$ENVIRONMENT # Rollback seeds"
echo "  go run cmd/seed/main.go -action=status -env=$ENVIRONMENT   # Check status"

# Make the files visible in common editors
if command -v code > /dev/null 2>&1; then
    echo ""
    echo "Opening files in VS Code..."
    code "$UP_FILE" "$DOWN_FILE"
fi