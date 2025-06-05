package migration

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

const (
	createSeedTrackingTable = `
CREATE TABLE IF NOT EXISTS seed_migrations (
    version VARCHAR(255) NOT NULL,
    environment VARCHAR(50) NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (version, environment)
)`
)

// InitSeedTracking creates the seed tracking table if it doesn't exist
func InitSeedTracking(db *sql.DB) error {
	_, err := db.Exec(createSeedTrackingTable)
	if err != nil {
		return fmt.Errorf("failed to create seed tracking table: %v", err)
	}
	return nil
}

// GetAppliedSeeds returns a list of seed versions that have been applied for an environment
func GetAppliedSeeds(environment string) ([]string, error) {
	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize seed tracking
	if err := InitSeedTracking(db); err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT version FROM seed_migrations WHERE environment = ? ORDER BY version", environment)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied seeds: %v", err)
	}
	defer rows.Close()

	var appliedSeeds []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan seed version: %v", err)
		}
		appliedSeeds = append(appliedSeeds, version)
	}

	return appliedSeeds, nil
}

// MarkSeedApplied marks a seed as applied for an environment
func MarkSeedApplied(environment, version string) error {
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

	_, err = db.Exec("INSERT INTO seed_migrations (version, environment) VALUES (?, ?)", version, environment)
	if err != nil {
		return fmt.Errorf("failed to mark seed as applied: %v", err)
	}

	return nil
}

// MarkSeedUnapplied removes a seed from the applied list (for rollback)
func MarkSeedUnapplied(environment, version string) error {
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

	_, err = db.Exec("DELETE FROM seed_migrations WHERE version = ? AND environment = ?", version, environment)
	if err != nil {
		return fmt.Errorf("failed to mark seed as unapplied: %v", err)
	}

	return nil
}

// GetPendingSeeds returns seeds that haven't been applied yet
func GetPendingSeeds(environment string) ([]string, error) {
	// Get all seed files for the environment
	seedsDir := filepath.Join("seeds", environment)
	files, err := filepath.Glob(filepath.Join(seedsDir, "*.up.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to read seed files: %v", err)
	}

	// Extract versions from filenames
	var allVersions []string
	for _, file := range files {
		filename := filepath.Base(file)
		// Extract version (everything before the first underscore)
		parts := strings.Split(filename, "_")
		if len(parts) > 0 {
			allVersions = append(allVersions, parts[0])
		}
	}

	// Get applied seeds
	appliedSeeds, err := GetAppliedSeeds(environment)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	appliedMap := make(map[string]bool)
	for _, applied := range appliedSeeds {
		appliedMap[applied] = true
	}

	// Find pending seeds
	var pendingSeeds []string
	for _, version := range allVersions {
		if !appliedMap[version] {
			pendingSeeds = append(pendingSeeds, version)
		}
	}

	// Sort by version number
	sort.Slice(pendingSeeds, func(i, j int) bool {
		vi, _ := strconv.Atoi(pendingSeeds[i])
		vj, _ := strconv.Atoi(pendingSeeds[j])
		return vi < vj
	})

	return pendingSeeds, nil
}

// GetSeedStatus returns detailed information about seed status
func GetDetailedSeedStatus(environment string) (map[string]interface{}, error) {
	appliedSeeds, err := GetAppliedSeeds(environment)
	if err != nil {
		return nil, err
	}

	pendingSeeds, err := GetPendingSeeds(environment)
	if err != nil {
		return nil, err
	}

	// Get all available seeds
	seedsDir := filepath.Join("seeds", environment)
	files, err := filepath.Glob(filepath.Join(seedsDir, "*.up.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to read seed files: %v", err)
	}

	var allSeeds []string
	for _, file := range files {
		filename := filepath.Base(file)
		parts := strings.Split(filename, "_")
		if len(parts) > 0 {
			allSeeds = append(allSeeds, parts[0])
		}
	}

	// Sort all seeds
	sort.Slice(allSeeds, func(i, j int) bool {
		vi, _ := strconv.Atoi(allSeeds[i])
		vj, _ := strconv.Atoi(allSeeds[j])
		return vi < vj
	})

	var latestApplied string
	if len(appliedSeeds) > 0 {
		latestApplied = appliedSeeds[len(appliedSeeds)-1]
	}

	return map[string]interface{}{
		"environment":     environment,
		"applied_seeds":   appliedSeeds,
		"pending_seeds":   pendingSeeds,
		"all_seeds":       allSeeds,
		"latest_applied":  latestApplied,
		"total_available": len(allSeeds),
		"total_applied":   len(appliedSeeds),
		"total_pending":   len(pendingSeeds),
	}, nil
}