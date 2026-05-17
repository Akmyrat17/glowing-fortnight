package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Middleware MiddlewareConfig `mapstructure:"middleware"`
	Firebase   FirebaseConfig   `mapstructure:"firebase"`
	JWT        JWTConfig        `mapstructure:"jwt"`
}

type ServerConfig struct {
	Host                    string        `mapstructure:"host"`
	Port                    int           `mapstructure:"port"`
	Env                     string        `mapstructure:"env"`
	ReadTimeout             time.Duration `mapstructure:"read_timeout"`
	WriteTimeout            time.Duration `mapstructure:"write_timeout"`
	IdleTimeout             time.Duration `mapstructure:"idle_timeout"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
	HideBanner              bool          `mapstructure:"hide_banner"`
	HidePort                bool          `mapstructure:"hide_port"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type MiddlewareConfig struct {
	CORS        CORSConfig        `mapstructure:"cors"`
	RateLimit   RateLimitConfig   `mapstructure:"rate_limit"`
	RequestSize RequestSizeConfig `mapstructure:"request_size"`
}

type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

type RateLimitConfig struct {
	Enabled bool    `mapstructure:"enabled"`
	Rate    float64 `mapstructure:"rate"`
	Burst   int     `mapstructure:"burst"`
}

type RequestSizeConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	MaxSize string `mapstructure:"max_size"`
}

type FirebaseConfig struct {
	CredentialsFile string `mapstructure:"credentials_file"`
}

type JWTConfig struct {
	Secret                 string        `mapstructure:"secret"`
	AccessTokenExpiration  time.Duration `mapstructure:"access_token_expiration"`
	RefreshTokenExpiration time.Duration `mapstructure:"refresh_token_expiration"`
}

type Manager struct {
	v         *viper.Viper
	config    *Config
	mu        sync.RWMutex
	callbacks []func(*Config)
}

var (
	manager *Manager
	once    sync.Once
)

func Initialize(configPath string) error {
	var initError error

	once.Do(func() {
		manager = &Manager{}
		initError = manager.load(configPath)
	})

	return initError
}

func Load(configPath string) (*Config, error) {
	if err := Initialize(configPath); err != nil {
		return nil, err
	}

	return Get(), nil
}

func Get() *Config {
	if manager == nil || manager.config == nil {
		return nil
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	cfg := *manager.config
	return &cfg
}

func OnReload(callback func(*Config)) {
	if manager == nil {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	manager.callbacks = append(manager.callbacks, callback)
}

func (m *Manager) load(configPath string) error {
	v := viper.New()
	m.v = v

	setDefaults(v)
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config: %w", err)
		}
		log.Printf("No config file found at %q, using defaults and environment values", configPath)
	} else {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	}

	if err := mergeEnvFile(v, ".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to merge root .env overrides: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	m.config = &cfg

	if v.ConfigFileUsed() != "" {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("config file changed: %s", e.Name)
			m.reload()
		})
	}

	return nil
}

func (m *Manager) reload() {
	var cfg Config
	if err := m.v.Unmarshal(&cfg); err != nil {
		log.Printf("failed to reload config: %v", err)
		return
	}

	if err := cfg.Validate(); err != nil {
		log.Printf("reloaded config validation failed: %v", err)
		return
	}

	m.mu.Lock()
	m.config = &cfg
	callbacks := append([]func(*Config){}, m.callbacks...)
	m.mu.Unlock()

	for _, callback := range callbacks {
		callback(&cfg)
	}
	log.Printf("configuration reloaded successfully")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.env", "development")
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.graceful_shutdown_timeout", "10s")
	v.SetDefault("server.hide_banner", true)
	v.SetDefault("server.hide_port", false)

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "boilerplate")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "text")
	v.SetDefault("logger.output", "stdout")

	v.SetDefault("middleware.cors.enabled", true)
	v.SetDefault("middleware.cors.allow_origins", []string{"*"})
	v.SetDefault("middleware.cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	v.SetDefault("middleware.cors.allow_headers", []string{"*"})
	v.SetDefault("middleware.cors.allow_credentials", false)
	v.SetDefault("middleware.cors.max_age", 3600)
	v.SetDefault("middleware.rate_limit.enabled", false)
	v.SetDefault("middleware.rate_limit.rate", 20.0)
	v.SetDefault("middleware.rate_limit.burst", 5)
	v.SetDefault("middleware.request_size.enabled", true)
	v.SetDefault("middleware.request_size.max_size", "2M")

	v.SetDefault("firebase.credentials_file", "./configs/firebase-credentials.json")

	v.SetDefault("jwt.secret", "your-secret-key-change-this-in-production")
	v.SetDefault("jwt.access_token_expiration", "30m")
	v.SetDefault("jwt.refresh_token_expiration", "168h")
}

func mergeEnvFile(v *viper.Viper, path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}

	v.SetConfigFile(path)
	v.SetConfigType("env")
	return v.MergeInConfig()
}

func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be a valid port: %d", c.Server.Port)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("database.user is required")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("database.name is required")
	}

	if c.Redis.Port <= 0 || c.Redis.Port > 65535 {
		return fmt.Errorf("redis.port must be a valid port: %d", c.Redis.Port)
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}

	return nil
}

func (dc DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dc.User, dc.Password, dc.Host, dc.Port, dc.Name, dc.SSLMode)
}
