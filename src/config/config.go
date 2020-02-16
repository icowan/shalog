package config

import (
	"github.com/Unknwon/goconfig"
	"strings"
)

const (
	SectionServer = "server"
	ImageFilePath = "image_file_path"
	ImageDomain   = "image_domain"
)

type Config struct {
	*goconfig.ConfigFile
}

func NewConfig(path string) (*Config, error) {
	// 处理配置文件

	cfg, err := goconfig.LoadConfigFile(path)
	if err != nil {
		return nil, err
	}
	return &Config{cfg}, nil
}

func (c *Config) GetString(section, key string) string {
	var val string
	val, _ = c.GetValue(section, key)
	return val
}

func (c *Config) GetStrings(section, key string) []string {
	val := c.GetString(section, key)
	return strings.Split(val, ";")
}

func (c *Config) GetInt(section, key string) int {
	val, _ := c.Int(section, key)
	return val
}

func (c *Config) GetBool(section, key string) bool {
	val, _ := c.Bool(section, key)
	return val
}
