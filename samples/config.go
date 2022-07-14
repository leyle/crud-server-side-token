package main

import (
	"github.com/leyle/go-api-starter/confighelper"
)

const (
	LogFormatJson = "json"
	LogFormatLine = "line"
)

type Config struct {
	Debug  bool                           `yaml:"debug"`
	Server *confighelper.ConnectionOption `yaml:"server"`
	SST    *SSTOption                     `yaml:"sst"`
	Log    *LogOption                     `yaml:"log"`
}

type SSTOption struct {
	AesKey      string `yaml:"aesKey"`
	ServiceName string `yaml:"serviceName"`
}

type LogOption struct {
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
}
