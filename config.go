package main

type AppConfig struct {
	TemplateFolder string
	QueueSize      int
}

var Config AppConfig

func NewAppConfig(projectFolder string) AppConfig {
	return AppConfig{
		TemplateFolder: projectFolder + "/templates",
		QueueSize:      50,
	}
}

func init() {
	Config = NewAppConfig("")
}
