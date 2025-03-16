package config

import (
	"github.com/Mark-Grigorev/chat_analyzer/internal/model"
	"github.com/kelseyhightower/envconfig"
)

func Read() *model.Config {
	var config model.Config
	envconfig.Process("", &config)
	return &config
}
