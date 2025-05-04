package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type FolderConfig struct {
	Folder  string   `yaml:"folder"`
	Depth   int      `yaml:"depth"`
	Ignores []string `yaml:"ignores"`
}

type Config struct {
	Folders []FolderConfig `yaml:"folders"`
	Ignores []string       `yaml:"ignores"`
}

func GetConfig() (conf Config, err error) {
	configPath, err := GetCurDirFile("./config.yaml")
	if err != nil {
		return
	}
	yamlFile, err := os.ReadFile(configPath)
	// 本地调试使用的文件
	if yamlFile == nil {
		yamlFile, err = os.ReadFile("config.yaml")
	}
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return
	}

	return
}
