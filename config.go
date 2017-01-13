package main

import (
	"path/filepath"
)

type AppConfig struct {
	MediaPlayer    string
	MediaFolder    string
	TemplateFolder string
	QueueSize      int
}

var Config AppConfig

func NewAppConfig() AppConfig {
	rootPath, _ := filepath.Abs(".")
	return AppConfig{
		MediaFolder:    rootPath + "/media",
		TemplateFolder: rootPath + "/templates",
		QueueSize:      50,
	}
}

func init() {
	Config = NewAppConfig()
}
