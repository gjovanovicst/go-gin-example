package setting

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

// Setup initialize the configuration instance
func Setup() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load App settings
	AppSetting.JwtSecret = getEnvString("JWT_SECRET", "233")
	AppSetting.PageSize = getEnvInt("APP_PAGE_SIZE", 10)
	AppSetting.PrefixUrl = getEnvString("PREFIX_URL", "http://127.0.0.1:8000")
	AppSetting.RuntimeRootPath = getEnvString("RUNTIME_ROOT_PATH", "runtime/")
	AppSetting.ImageSavePath = getEnvString("IMAGE_SAVE_PATH", "upload/images/")
	AppSetting.ImageMaxSize = getEnvInt("IMAGE_MAX_SIZE", 5)
	AppSetting.ImageAllowExts = getEnvStringSlice("IMAGE_ALLOW_EXTS", []string{".jpg", ".jpeg", ".png"})
	AppSetting.ExportSavePath = getEnvString("EXPORT_SAVE_PATH", "export/")
	AppSetting.QrCodeSavePath = getEnvString("QR_CODE_SAVE_PATH", "qrcode/")
	AppSetting.FontSavePath = getEnvString("FONT_SAVE_PATH", "fonts/")
	AppSetting.LogSavePath = getEnvString("LOG_SAVE_PATH", "logs/")
	AppSetting.LogSaveName = getEnvString("LOG_SAVE_NAME", "log")
	AppSetting.LogFileExt = getEnvString("LOG_FILE_EXT", "log")
	AppSetting.TimeFormat = getEnvString("TIME_FORMAT", "20060102")

	// Load Server settings
	ServerSetting.RunMode = getEnvString("RUN_MODE", "debug")
	ServerSetting.HttpPort = getEnvInt("HTTP_PORT", 8000)
	ServerSetting.ReadTimeout = time.Duration(getEnvInt("READ_TIMEOUT", 60))
	ServerSetting.WriteTimeout = time.Duration(getEnvInt("WRITE_TIMEOUT", 60))

	// Load Database settings
	DatabaseSetting.Type = getEnvString("DB_TYPE", "mysql")
	DatabaseSetting.User = getEnvString("DB_USER", "root")
	DatabaseSetting.Password = getEnvString("DB_PASSWORD", "rootpassword")
	DatabaseSetting.Host = getEnvString("DB_HOST", "127.0.0.1:3306")
	DatabaseSetting.Name = getEnvString("DB_NAME", "blog")
	DatabaseSetting.TablePrefix = getEnvString("DB_TABLE_PREFIX", "blog_")

	// Load Redis settings
	RedisSetting.Host = getEnvString("REDIS_HOST", "127.0.0.1:6379")
	RedisSetting.Password = getEnvString("REDIS_PASSWORD", "")
	RedisSetting.MaxIdle = getEnvInt("REDIS_MAX_IDLE", 30)
	RedisSetting.MaxActive = getEnvInt("REDIS_MAX_ACTIVE", 30)
	RedisSetting.IdleTimeout = time.Duration(getEnvInt("REDIS_IDLE_TIMEOUT", 200))

	// Apply transformations
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// getEnvString gets string value from environment variable with default
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets integer value from environment variable with default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvStringSlice gets string slice from environment variable with default
func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
