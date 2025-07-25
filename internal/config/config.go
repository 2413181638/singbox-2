package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DatabasePath string      `yaml:"database_path" mapstructure:"database_path"`
	LogLevel     string      `yaml:"log_level" mapstructure:"log_level"`
	XBoard       XBoardConfig `yaml:"xboard" mapstructure:"xboard"`
	SingBox      SingBoxConfig `yaml:"singbox" mapstructure:"singbox"`
}

type XBoardConfig struct {
	URL      string `yaml:"url" mapstructure:"url"`
	Token    string `yaml:"token" mapstructure:"token"`
	NodeID   int    `yaml:"node_id" mapstructure:"node_id"`
	Interval int    `yaml:"interval" mapstructure:"interval"` // 订阅更新间隔(秒)
}

type SingBoxConfig struct {
	ConfigPath string `yaml:"config_path" mapstructure:"config_path"`
	LogPath    string `yaml:"log_path" mapstructure:"log_path"`
	APIPort    int    `yaml:"api_port" mapstructure:"api_port"`
}

func Load() (*Config, error) {
	// 获取配置文件路径
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "config.yaml")

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := getDefaultConfig(configDir)
		if err := saveConfig(cfg, configPath); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	// 读取配置文件
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func getDefaultConfig(configDir string) *Config {
	return &Config{
		DatabasePath: filepath.Join(configDir, "data.db"),
		LogLevel:     "info",
		XBoard: XBoardConfig{
			URL:      "",
			Token:    "",
			NodeID:   1,
			Interval: 300, // 5分钟
		},
		SingBox: SingBoxConfig{
			ConfigPath: filepath.Join(configDir, "singbox.json"),
			LogPath:    filepath.Join(configDir, "singbox.log"),
			APIPort:    9090,
		},
	}
}

func saveConfig(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".singbox-xboard")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

func (c *Config) Save() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	return saveConfig(c, configPath)
}