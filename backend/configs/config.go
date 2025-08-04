package configs

import (
	"fmt"
	"os"
	"strconv"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	K8s      K8sConfig      `mapstructure:"kubernetes"`
	Docker   DockerConfig   `mapstructure:"docker"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string `mapstructure:"secret_key"`
	ExpiresIn int    `mapstructure:"expires_in"`
}

// K8sConfig holds Kubernetes configuration
type K8sConfig struct {
	ConfigPath string `mapstructure:"config_path"`
	Namespace  string `mapstructure:"namespace"`
}

// DockerConfig holds Docker configuration
type DockerConfig struct {
	RegistryURL string `mapstructure:"registry_url"`
	RegistryUser string `mapstructure:"registry_user"`
	RegistryPassword string `mapstructure:"registry_password"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Read environment variables
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// If config file exists, read it
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.idle_timeout", 120)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "minipaas")
	viper.SetDefault("database.sslmode", "disable")

	// JWT defaults
	viper.SetDefault("jwt.secret_key", "your-secret-key")
	viper.SetDefault("jwt.expires_in", 24) // hours

	// Kubernetes defaults
	viper.SetDefault("kubernetes.config_path", "")
	viper.SetDefault("kubernetes.namespace", "default")

	// Docker defaults
	viper.SetDefault("docker.registry_url", "docker.io")
	viper.SetDefault("docker.registry_user", "")
	viper.SetDefault("docker.registry_password", "")
}

// GetDSN returns database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetServerAddr returns server address
func (c *ServerConfig) GetServerAddr() string {
	return c.Host + ":" + c.Port
}

// IsDevelopment returns true if running in development mode
func IsDevelopment() bool {
	return os.Getenv("GIN_MODE") != "release"
} 