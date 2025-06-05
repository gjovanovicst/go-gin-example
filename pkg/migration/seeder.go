package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

// GetEnvironment returns the current environment (development, production, etc.)
func GetEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		// Default to development if not set
		if setting.ServerSetting.RunMode == "debug" {
			return "development"
		}
		return "production"
	}
	return env
}

// RunSeeds executes all seed files for the specified environment
func RunSeeds(environment string) error {
	if environment == "" {
		environment = GetEnvironment()
	}

	log.Printf("Running seeds for environment: %s", environment)

	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migration driver
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %v", err)
	}

	// Construct the seeds path for the environment
	seedsPath := fmt.Sprintf("file://seeds/%s", environment)

	// Create migrate instance for seeds
	m, err := migrate.NewWithDatabaseInstance(
		seedsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance for seeds: %v", err)
	}

	// Run seed migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run seeds: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Printf("No new seeds to run for environment: %s", environment)
	} else {
		log.Printf("Seeds completed successfully for environment: %s", environment)
	}

	return nil
}

// RollbackSeeds rolls back the last seed for the specified environment
func RollbackSeeds(environment string) error {
	if environment == "" {
		environment = GetEnvironment()
	}

	log.Printf("Rolling back seeds for environment: %s", environment)

	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migration driver
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %v", err)
	}

	// Construct the seeds path for the environment
	seedsPath := fmt.Sprintf("file://seeds/%s", environment)

	// Create migrate instance for seeds
	m, err := migrate.NewWithDatabaseInstance(
		seedsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance for seeds: %v", err)
	}

	// Roll back one step
	err = m.Steps(-1)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback seeds: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Printf("No seeds to rollback for environment: %s", environment)
	} else {
		log.Printf("Seed rollback completed successfully for environment: %s", environment)
	}

	return nil
}

// ListAvailableEnvironments returns a list of available seed environments
func ListAvailableEnvironments() ([]string, error) {
	seedsDir := "seeds"
	
	entries, err := os.ReadDir(seedsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read seeds directory: %v", err)
	}

	var environments []string
	for _, entry := range entries {
		if entry.IsDir() {
			environments = append(environments, entry.Name())
		}
	}

	sort.Strings(environments)
	return environments, nil
}

// GetSeedStatus returns the current seed status for an environment
func GetSeedStatus(environment string) (uint, bool, error) {
	if environment == "" {
		environment = GetEnvironment()
	}

	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return 0, false, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migration driver
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration driver: %v", err)
	}

	// Construct the seeds path for the environment
	seedsPath := fmt.Sprintf("file://seeds/%s", environment)

	// Create migrate instance for seeds
	m, err := migrate.NewWithDatabaseInstance(
		seedsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance for seeds: %v", err)
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get seed version: %v", err)
	}

	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}

	return version, dirty, nil
}

// RunSeedsWithTracking executes seed files with proper tracking (recommended)
func RunSeedsWithTracking(environment string) error {
	if environment == "" {
		environment = GetEnvironment()
	}

	log.Printf("Running tracked seeds for environment: %s", environment)

	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize seed tracking
	if err := InitSeedTracking(db); err != nil {
		return err
	}

	// Get pending seeds
	pendingSeeds, err := GetPendingSeeds(environment)
	if err != nil {
		return err
	}

	if len(pendingSeeds) == 0 {
		log.Printf("No pending seeds for environment: %s", environment)
		return nil
	}

	log.Printf("Found %d pending seeds: %v", len(pendingSeeds), pendingSeeds)

	// Execute each pending seed
	for _, version := range pendingSeeds {
		seedFile := filepath.Join("seeds", environment, fmt.Sprintf("%s_*.up.sql", version))
		files, err := filepath.Glob(seedFile)
		if err != nil || len(files) == 0 {
			log.Printf("Warning: Could not find seed file for version %s", version)
			continue
		}

		file := files[0] // Take the first match
		log.Printf("Executing seed file: %s", file)

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %v", file, err)
		}

		// Execute the seed
		if err := executeSeedContent(db, file, content); err != nil {
			return err
		}

		// Mark as applied
		if err := MarkSeedApplied(environment, version); err != nil {
			return fmt.Errorf("failed to mark seed %s as applied: %v", version, err)
		}

		log.Printf("Successfully applied seed version %s", version)
	}

	log.Printf("All pending seeds applied successfully for environment: %s", environment)
	return nil
}

// RunSeedsManually executes seed files manually by reading and executing SQL files
// This is useful when you need more control over the seeding process
func RunSeedsManually(environment string) error {
	if environment == "" {
		environment = GetEnvironment()
	}

	log.Printf("Running seeds manually for environment: %s", environment)

	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get seed files for the environment
	seedsDir := filepath.Join("seeds", environment)
	files, err := filepath.Glob(filepath.Join(seedsDir, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("failed to read seed files: %v", err)
	}

	if len(files) == 0 {
		log.Printf("No seed files found for environment: %s", environment)
		return nil
	}

	// Sort files to ensure they run in order
	sort.Strings(files)

	// Execute each seed file
	for _, file := range files {
		log.Printf("Executing seed file: %s", file)
		
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %v", file, err)
		}

		// Process the content to handle multi-line SQL statements properly
		contentStr := strings.TrimSpace(string(content))
		
		// Remove comments and empty lines, but preserve the SQL structure
		lines := strings.Split(contentStr, "\n")
		var cleanLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			// Skip empty lines and comment-only lines
			if line != "" && !strings.HasPrefix(line, "--") {
				cleanLines = append(cleanLines, line)
			}
		}
		
		if len(cleanLines) == 0 {
			log.Printf("No SQL statements found in %s, skipping", file)
			continue
		}
		
		// Join all lines and split by semicolon to get complete statements
		fullContent := strings.Join(cleanLines, "\n")
		statements := strings.Split(fullContent, ";")
		
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			
			// Show first 100 characters of statement for logging
			preview := stmt
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			log.Printf("Executing SQL from %s: %s", filepath.Base(file), preview)
			_, err = db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("failed to execute statement in %s: %v\nStatement: %s", file, err, stmt)
			}
		}
		
		log.Printf("Successfully executed seed file: %s", file)
	}

	log.Printf("All seeds executed successfully for environment: %s", environment)
	return nil
}

// executeSeedContent executes the SQL content from a seed file
func executeSeedContent(db *sql.DB, file string, content []byte) error {
	// Process the content to handle multi-line SQL statements properly
	contentStr := strings.TrimSpace(string(content))
	
	// Remove comments and empty lines, but preserve the SQL structure
	lines := strings.Split(contentStr, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comment-only lines
		if line != "" && !strings.HasPrefix(line, "--") {
			cleanLines = append(cleanLines, line)
		}
	}
	
	if len(cleanLines) == 0 {
		log.Printf("No SQL statements found in %s, skipping", file)
		return nil
	}
	
	// Join all lines and split by semicolon to get complete statements
	fullContent := strings.Join(cleanLines, "\n")
	statements := strings.Split(fullContent, ";")
	
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		
		// Show first 100 characters of statement for logging
		preview := stmt
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		log.Printf("Executing SQL from %s: %s", filepath.Base(file), preview)
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement in %s: %v\nStatement: %s", file, err, stmt)
		}
	}
	return nil
}

// RollbackLastSeed rolls back the last applied seed for an environment
func RollbackLastSeed(environment string) error {
	if environment == "" {
		environment = GetEnvironment()
	}

	log.Printf("Rolling back last seed for environment: %s", environment)

	// Get applied seeds
	appliedSeeds, err := GetAppliedSeeds(environment)
	if err != nil {
		return err
	}

	if len(appliedSeeds) == 0 {
		log.Printf("No seeds to rollback for environment: %s", environment)
		return nil
	}

	// Get the last applied seed
	lastSeed := appliedSeeds[len(appliedSeeds)-1]
	log.Printf("Rolling back seed version: %s", lastSeed)

	// Find the down file
	downFile := filepath.Join("seeds", environment, fmt.Sprintf("%s_*.down.sql", lastSeed))
	files, err := filepath.Glob(downFile)
	if err != nil || len(files) == 0 {
		return fmt.Errorf("could not find rollback file for seed version %s", lastSeed)
	}

	file := files[0] // Take the first match
	log.Printf("Executing rollback file: %s", file)

	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read rollback file %s: %v", file, err)
	}

	// Execute the rollback
	if err := executeSeedContent(db, file, content); err != nil {
		return err
	}

	// Mark as unapplied
	if err := MarkSeedUnapplied(environment, lastSeed); err != nil {
		return fmt.Errorf("failed to mark seed %s as unapplied: %v", lastSeed, err)
	}

	log.Printf("Successfully rolled back seed version %s", lastSeed)
	return nil
}