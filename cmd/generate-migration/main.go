package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

type ColumnInfo struct {
	Name    string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

type ModelField struct {
	Name     string
	DBName   string
	Type     string
	Size     int
	Default  string
	Comment  string
	Nullable bool
}

type ModelInfo struct {
	Name      string
	TableName string
	Instance  interface{}
}

func main() {
	var (
		name = flag.String("name", "", "Migration name")
		auto = flag.Bool("auto", false, "Auto-generate migration based on model changes")
	)
	flag.Parse()

	if *name == "" {
		log.Fatal("Migration name is required. Use -name flag.")
	}

	// Load configuration
	setting.Setup()

	// Get next migration version
	nextVersion, err := getNextMigrationVersion()
	if err != nil {
		log.Fatalf("Failed to get next migration version: %v", err)
	}

	upFile := fmt.Sprintf("migrations/%d_%s.up.sql", nextVersion, *name)
	downFile := fmt.Sprintf("migrations/%d_%s.down.sql", nextVersion, *name)

	var upSQL, downSQL string

	if *auto {
		// Auto-generate migration based on model analysis
		upSQL, downSQL, err = generateMigrationFromModels(*name)
		if err != nil {
			log.Fatalf("Failed to generate migration: %v", err)
		}
	} else {
		// Create template migration
		upSQL = fmt.Sprintf("-- Migration: %s\n-- Add your forward migration SQL here\n", *name)
		downSQL = fmt.Sprintf("-- Rollback migration: %s\n-- Add your rollback migration SQL here\n", *name)
	}

	// Write migration files
	err = os.WriteFile(upFile, []byte(upSQL), 0644)
	if err != nil {
		log.Fatalf("Failed to create up migration file: %v", err)
	}

	err = os.WriteFile(downFile, []byte(downSQL), 0644)
	if err != nil {
		log.Fatalf("Failed to create down migration file: %v", err)
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  %s\n", upFile)
	fmt.Printf("  %s\n", downFile)

	if *auto {
		fmt.Printf("\nGenerated SQL:\n%s\n", upSQL)
	}

	fmt.Printf("\nTo apply migration, run:\n")
	fmt.Printf("  go run cmd/migrate/main.go -action=up\n")
}

func getNextMigrationVersion() (int, error) {
	maxFileVersion := 0
	err := filepath.Walk("migrations", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".up.sql") {
			// Extract version number from filename
			re := regexp.MustCompile(`^(\d+)_`)
			matches := re.FindStringSubmatch(info.Name())
			if len(matches) > 1 {
				version, err := strconv.Atoi(matches[1])
				if err == nil && version > maxFileVersion {
					maxFileVersion = version
				}
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return maxFileVersion + 1, nil
}

func generateMigrationFromModels(migrationName string) (string, string, error) {
	var upSQL, downSQL strings.Builder

	upSQL.WriteString(fmt.Sprintf("-- Migration: %s\n-- Auto-generated based on model changes\n\n", migrationName))
	downSQL.WriteString(fmt.Sprintf("-- Rollback migration: %s\n-- Auto-generated rollback\n\n", migrationName))

	// Analyze models and generate SQL
	changes, err := detectModelChanges()
	if err != nil {
		return "", "", err
	}

	if len(changes) == 0 {
		upSQL.WriteString("-- No model changes detected\n")
		downSQL.WriteString("-- No rollback needed\n")
	} else {
		for _, change := range changes {
			upSQL.WriteString(change.UpSQL + "\n")
			downSQL.WriteString(change.DownSQL + "\n")
		}
	}

	return upSQL.String(), downSQL.String(), nil
}

type ModelChange struct {
	Type    string // "add_column", "drop_column", "modify_column"
	Table   string
	Column  string
	UpSQL   string
	DownSQL string
}

func detectModelChanges() ([]ModelChange, error) {
	var changes []ModelChange

	// Connect to database to get current schema
	db, err := connectToDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get all registered models dynamically
	modelRegistry := getModelRegistry()
	
	for _, modelInfo := range modelRegistry {
		fmt.Printf("Analyzing %s model...\n", modelInfo.Name)
		
		// Get actual table name using GORM naming conventions
		tableName := getTableName(modelInfo.Instance)
		
		modelChanges, err := analyzeModelChanges(db, tableName, modelInfo.Instance)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze %s model: %v", modelInfo.Name, err)
		}
		changes = append(changes, modelChanges...)
	}

	fmt.Printf("Total changes detected: %d\n", len(changes))
	return changes, nil
}

// getModelRegistry returns all models that should be analyzed for schema changes
func getModelRegistry() []ModelInfo {
	// For now, return known models, but this could be made fully dynamic
	// by using reflection to scan the models package
	modelRegistry := []ModelInfo{
		{Name: "Tag", Instance: models.Tag{}},
		{Name: "Article", Instance: models.Article{}},
		{Name: "Auth", Instance: models.Auth{}},
	}
	
	// Optionally, add auto-discovery logic here
	// discoveredModels := discoverModelsInPackage("github.com/EDDYCJY/go-gin-example/models")
	// modelRegistry = append(modelRegistry, discoveredModels...)
	
	return modelRegistry
}

// discoverModelsInPackage would dynamically discover all model structs in a package
// This is commented out as it requires advanced reflection techniques
/*
func discoverModelsInPackage(packagePath string) []ModelInfo {
	// This would use reflection to scan the models package and find all structs
	// that embed the Model struct or have GORM tags
	// Implementation would require runtime package scanning
	return []ModelInfo{}
}
*/

// getTableName gets the actual table name using GORM conventions
func getTableName(model interface{}) string {
	// Load settings to get table prefix
	setting.Setup()
	
	// Create a temporary GORM DB to get table name
	tempDB, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))
	
	if err != nil {
		// Fallback to simple naming if GORM connection fails
		modelType := reflect.TypeOf(model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		return setting.DatabaseSetting.TablePrefix + strings.ToLower(modelType.Name())
	}
	
	defer tempDB.Close()
	
	// Configure GORM with same settings as models package
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}
	tempDB.SingularTable(true)
	
	// Get the actual table name GORM would use
	return tempDB.NewScope(model).TableName()
}

func connectToDatabase() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)
	return sql.Open("mysql", dsn)
}

func analyzeModelChanges(db *sql.DB, tableName string, model interface{}) ([]ModelChange, error) {
	var changes []ModelChange

	// Get current database columns
	dbColumns, err := getTableColumns(db, tableName)
	if err != nil {
		return nil, err
	}

	// Get model fields
	modelFields := getModelFields(model)

	// Create maps for easier lookup
	dbColumnMap := make(map[string]ColumnInfo)
	for _, col := range dbColumns {
		dbColumnMap[col.Name] = col
	}
	
	modelFieldMap := make(map[string]ModelField)
	for _, field := range modelFields {
		modelFieldMap[field.DBName] = field
	}

	// Find new columns (in model but not in DB)
	for _, field := range modelFields {
		if _, exists := dbColumnMap[field.DBName]; !exists {
			change := ModelChange{
				Type:    "add_column",
				Table:   tableName,
				Column:  field.DBName,
				UpSQL:   generateAddColumnSQL(tableName, field),
				DownSQL: generateDropColumnSQL(tableName, field.DBName),
			}
			changes = append(changes, change)
		}
	}

	// Find removed columns (in DB but not in model)
	for _, col := range dbColumns {
		if _, exists := modelFieldMap[col.Name]; !exists {
			// Skip system columns that might be added by GORM or MySQL
			if isSystemColumn(col.Name) {
				continue
			}
			
			change := ModelChange{
				Type:    "drop_column",
				Table:   tableName,
				Column:  col.Name,
				UpSQL:   generateDropColumnSQL(tableName, col.Name),
				DownSQL: generateAddColumnSQLFromDB(tableName, col),
			}
			changes = append(changes, change)
		}
	}

	// Find modified columns (type changes, etc.)
	for _, field := range modelFields {
		if dbCol, exists := dbColumnMap[field.DBName]; exists {
			if needsColumnUpdate(field, dbCol) {
				change := ModelChange{
					Type:    "modify_column",
					Table:   tableName,
					Column:  field.DBName,
					UpSQL:   generateModifyColumnSQL(tableName, field),
					DownSQL: generateModifyColumnSQLFromDB(tableName, dbCol),
				}
				changes = append(changes, change)
			}
		}
	}

	return changes, nil
}

// isSystemColumn checks if a column is a system column that should be ignored
func isSystemColumn(columnName string) bool {
	systemColumns := []string{
		"created_at", "updated_at", "deleted_at", // Common GORM timestamps
	}
	
	for _, sysCol := range systemColumns {
		if columnName == sysCol {
			return true
		}
	}
	return false
}

// needsColumnUpdate checks if a model field differs from the database column
func needsColumnUpdate(field ModelField, dbCol ColumnInfo) bool {
	expectedType := mapGoTypeToSQL(field.Type)
	actualType := strings.ToUpper(dbCol.Type)
	
	// Simple type comparison - can be made more sophisticated
	return !strings.Contains(actualType, expectedType)
}

// generateModifyColumnSQL generates SQL to modify an existing column
func generateModifyColumnSQL(tableName string, field ModelField) string {
	sqlType := mapGoTypeToSQLWithSize(field.Type, field.Size, field.DBName)
	nullable := "NOT NULL"
	if field.Nullable {
		nullable = "NULL"
	}
	
	return fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` %s %s;",
		tableName, field.DBName, sqlType, nullable)
}

// generateModifyColumnSQLFromDB generates SQL to restore a column to its original state
func generateModifyColumnSQLFromDB(tableName string, col ColumnInfo) string {
	nullable := "NOT NULL"
	if col.Null == "YES" {
		nullable = "NULL"
	}
	
	defaultClause := ""
	if col.Default.Valid {
		defaultClause = fmt.Sprintf("DEFAULT '%s'", col.Default.String)
	}
	
	return fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` %s %s %s;",
		tableName, col.Name, col.Type, nullable, defaultClause)
}

// generateAddColumnSQLFromDB generates SQL to add a column based on DB info
func generateAddColumnSQLFromDB(tableName string, col ColumnInfo) string {
	nullable := "NOT NULL"
	if col.Null == "YES" {
		nullable = "NULL"
	}
	
	defaultClause := ""
	if col.Default.Valid {
		defaultClause = fmt.Sprintf("DEFAULT '%s'", col.Default.String)
	}
	
	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s %s %s;",
		tableName, col.Name, col.Type, nullable, defaultClause)
}

func getTableColumns(db *sql.DB, tableName string) ([]ColumnInfo, error) {
	query := "DESCRIBE " + tableName
	rows, err := db.Query(query)
	if err != nil {
		// Table might not exist yet, return empty
		return []ColumnInfo{}, nil
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		err := rows.Scan(&col.Name, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

func getModelFields(model interface{}) []ModelField {
	var fields []ModelField
	modelType := reflect.TypeOf(model)
	
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Process embedded Model fields recursively
	fields = append(fields, processStructFields(modelType, "")...)

	return fields
}

// processStructFields recursively processes struct fields including embedded structs
func processStructFields(structType reflect.Type, prefix string) []ModelField {
	var fields []ModelField
	
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		
		// Handle embedded structs (like Model)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			embeddedFields := processStructFields(field.Type, prefix)
			fields = append(fields, embeddedFields...)
			continue
		}
		
		// Skip non-anonymous struct fields that are relationships
		if field.Type.Kind() == reflect.Struct && !field.Anonymous {
			continue
		}
		
		// Skip slice/array fields (typically relationships)
		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			continue
		}
		
		// Skip pointer to struct fields (relationships)
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			continue
		}

		dbName := getDBColumnName(field)
		if dbName == "" || dbName == "-" {
			continue
		}

		// Parse GORM tags for additional field information
		gormInfo := parseGORMTags(field.Tag.Get("gorm"))
		
		modelField := ModelField{
			Name:     field.Name,
			DBName:   prefix + dbName,
			Type:     field.Type.String(),
			Nullable: !gormInfo.NotNull,
			Size:     gormInfo.Size,
			Default:  gormInfo.Default,
			Comment:  gormInfo.Comment,
		}

		fields = append(fields, modelField)
	}

	return fields
}

// GORMTagInfo holds parsed GORM tag information
type GORMTagInfo struct {
	NotNull   bool
	Size      int
	Default   string
	Comment   string
	Index     bool
	Unique    bool
	PrimaryKey bool
}

// parseGORMTags parses GORM struct tags
func parseGORMTags(gormTag string) GORMTagInfo {
	info := GORMTagInfo{}
	
	if gormTag == "" {
		return info
	}
	
	parts := strings.Split(gormTag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		
		switch {
		case part == "not null":
			info.NotNull = true
		case part == "primary_key":
			info.PrimaryKey = true
			info.NotNull = true
		case part == "index":
			info.Index = true
		case part == "unique":
			info.Unique = true
		case strings.HasPrefix(part, "size:"):
			if size, err := strconv.Atoi(strings.TrimPrefix(part, "size:")); err == nil {
				info.Size = size
			}
		case strings.HasPrefix(part, "default:"):
			info.Default = strings.TrimPrefix(part, "default:")
		case strings.HasPrefix(part, "comment:"):
			info.Comment = strings.TrimPrefix(part, "comment:")
		}
	}
	
	return info
}

func getDBColumnName(field reflect.StructField) string {
	// First check for gorm tag to get actual database column name
	gormTag := field.Tag.Get("gorm")
	if gormTag != "" {
		// Look for column name in gorm tag
		for _, part := range strings.Split(gormTag, ";") {
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:")
			}
		}
	}

	// Fall back to JSON tag as DB column name
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		// If no json tag, convert field name to snake_case
		return toSnakeCase(field.Name)
	}
	
	// Remove options like ,omitempty
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}

func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func generateAddColumnSQL(tableName string, field ModelField) string {
	sqlType := mapGoTypeToSQLWithSize(field.Type, field.Size, field.DBName)
	
	// Handle nullable
	nullable := "NULL"
	if !field.Nullable {
		nullable = "NOT NULL"
	}
	
	// Handle default value
	defaultClause := ""
	if field.Default != "" {
		if isNumericType(field.Type) {
			defaultClause = fmt.Sprintf("DEFAULT %s", field.Default)
		} else {
			defaultClause = fmt.Sprintf("DEFAULT '%s'", field.Default)
		}
	} else {
		// Set sensible defaults based on type
		if !field.Nullable {
			switch {
			case strings.Contains(field.Type, "string"):
				defaultClause = "DEFAULT ''"
			case strings.Contains(field.Type, "int"):
				defaultClause = "DEFAULT 0"
			case strings.Contains(field.Type, "bool"):
				defaultClause = "DEFAULT 0"
			}
		}
	}
	
	// Handle comment
	commentClause := ""
	if field.Comment != "" {
		commentClause = fmt.Sprintf("COMMENT '%s'", field.Comment)
	}

	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s %s %s %s;",
		tableName, field.DBName, sqlType, nullable, defaultClause, commentClause)
}

// mapGoTypeToSQLWithSize maps Go types to MySQL types with intelligent sizing
func mapGoTypeToSQLWithSize(goType string, gormSize int, fieldName string) string {
	switch {
	case goType == "string":
		return getStringColumnType(gormSize, fieldName)
	case goType == "int" || goType == "int32":
		return "INT"
	case goType == "int64":
		return "BIGINT"
	case goType == "int16":
		return "SMALLINT"
	case goType == "int8":
		return "TINYINT"
	case goType == "uint" || goType == "uint32":
		return "INT UNSIGNED"
	case goType == "uint64":
		return "BIGINT UNSIGNED"
	case goType == "uint16":
		return "SMALLINT UNSIGNED"
	case goType == "uint8":
		return "TINYINT UNSIGNED"
	case goType == "bool":
		return "TINYINT(1)"
	case goType == "float32":
		return "FLOAT"
	case goType == "float64":
		return "DOUBLE"
	case strings.Contains(goType, "time.Time"):
		return "DATETIME"
	case strings.Contains(goType, "sql.NullString"):
		return getStringColumnType(gormSize, fieldName)
	case strings.Contains(goType, "sql.NullInt"):
		return "INT"
	case strings.Contains(goType, "sql.NullBool"):
		return "TINYINT(1)"
	case strings.Contains(goType, "sql.NullFloat"):
		return "DOUBLE"
	case strings.Contains(goType, "sql.NullTime"):
		return "DATETIME"
	default:
		return "VARCHAR(255)" // Default fallback
	}
}

// getStringColumnType determines the appropriate SQL type for string fields
func getStringColumnType(gormSize int, fieldName string) string {
	// Priority 1: Use explicit GORM size if defined
	if gormSize > 0 {
		if gormSize > 65535 {
			return "LONGTEXT"
		} else if gormSize > 255 {
			return "TEXT"
		}
		return fmt.Sprintf("VARCHAR(%d)", gormSize)
	}
	
	// Priority 2: Use conservative defaults based on field name patterns
	fieldLower := strings.ToLower(fieldName)
	switch {
	case strings.Contains(fieldLower, "content") || strings.Contains(fieldLower, "description") || strings.Contains(fieldLower, "body"):
		return "TEXT" // For long text content
	case strings.Contains(fieldLower, "password"):
		return "VARCHAR(100)" // Conservative for passwords (bcrypt = 60, but allowing buffer)
	case strings.Contains(fieldLower, "username") || strings.Contains(fieldLower, "name"):
		return "VARCHAR(100)" // Standard name length
	case strings.Contains(fieldLower, "email"):
		return "VARCHAR(255)" // Email can be up to 320 chars, but 255 is common limit
	case strings.Contains(fieldLower, "title"):
		return "VARCHAR(200)" // Article/post titles
	case strings.Contains(fieldLower, "url") || strings.Contains(fieldLower, "link"):
		return "VARCHAR(500)" // URLs can be long
	case strings.Contains(fieldLower, "code") || strings.Contains(fieldLower, "token"):
		return "VARCHAR(255)" // API keys, tokens, codes
	default:
		return "VARCHAR(255)" // Conservative default
	}
}

// Keep the original function for backward compatibility in modify operations
func mapGoTypeToSQL(goType string) string {
	return mapGoTypeToSQLWithSize(goType, 0, "")
}

// isNumericType checks if a Go type is numeric
func isNumericType(goType string) bool {
	numericTypes := []string{
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "bool",
	}
	
	for _, numType := range numericTypes {
		if goType == numType {
			return true
		}
	}
	return false
}

func generateDropColumnSQL(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", tableName, columnName)
}