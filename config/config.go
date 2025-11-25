package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Log      LogConfig      `yaml:"log"`
}

type AppConfig struct {
	Name  string `yaml:"name"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	Dsn          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type LogConfig struct {
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`    // MB
	MaxAge     int    `yaml:"max_age"`     // 天
	MaxBackups int    `yaml:"max_backups"` // 个
	Compress   bool   `yaml:"compress"`
}

var Cfg Config // 全局配置变量，方便其他地方调用

// Init 初始化配置（建议在程序启动时调用）
func Init(path string) error {
	// 读取文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 解析YAML到结构体
	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		return err
	}

	return nil
}
