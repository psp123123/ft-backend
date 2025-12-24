package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	File     FileConfig     `yaml:"file"`
	Redis    RedisConfig    `yaml:"redis"`
}

type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type JWTConfig struct {
	SecretKey       string `yaml:"secret_key"`
	AccessTokenExp  int    `yaml:"access_token_exp"`
	RefreshTokenExp int    `yaml:"refresh_token_exp"`
}

type FileConfig struct {
	UploadDir      string   `yaml:"upload_dir"`
	MaxFileSize    int64    `yaml:"max_file_size"`
	ChunkSize      int      `yaml:"chunk_size"`
	AllowedFormats []string `yaml:"allowed_formats"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func LoadConfig() (*Config, error) {
	configFile := "config.yaml"

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// 如果不存在，创建默认配置文件
		defaultConfig := &Config{
			Server: ServerConfig{
				Host:         "0.0.0.0",
				Port:         "8080",
				ReadTimeout:  30,
				WriteTimeout: 30,
			},
			Database: DatabaseConfig{
				Host:     "192.168.56.11",
				Port:     "3306",
				User:     "root",
				Password: "123456",
				DBName:   "filetransfer",
				Charset:  "utf8mb4",
			},
			JWT: JWTConfig{
				SecretKey:       "your-secret-key-here",
				AccessTokenExp:  15,
				RefreshTokenExp: 1440,
			},
			File: FileConfig{
				UploadDir:      "uploads",
				MaxFileSize:    1073741824,
				ChunkSize:      1048576,
				AllowedFormats: []string{"jpg", "png", "pdf", "txt", "zip", "rar"},
			},
			Redis: RedisConfig{
				Host:     "localhost",
				Port:     "6379",
				Password: "",
				DB:       0,
			},
		}

		// 保存默认配置
		if err := SaveConfig(defaultConfig, configFile); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}

		return defaultConfig, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
