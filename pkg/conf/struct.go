package conf

import (
)

type ConfInfo struct {
	WorkDir string `yaml:"work-dir"`
	HttpServerListen string `yaml:"http-server-listen"`
	BaseUrlPath string `yaml:"base-url-path"`
	AllowRegister bool `yaml:"allow-register"`
	EnableCors bool `yaml:"enable-cors"`
	CorsAllowOriginUrl string `yaml:"cors-allow-origin-url"`
}

func (ci *ConfInfo) setDefaults() {
	if ci.WorkDir == "" {
		ci.WorkDir = "./"
	}
	if ci.HttpServerListen == "" {
		ci.HttpServerListen = ":9465"
	}
	if ci.CorsAllowOriginUrl == "" {
		ci.CorsAllowOriginUrl = "*"
	}
}