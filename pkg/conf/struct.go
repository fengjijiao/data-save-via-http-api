package conf

import (
)

type ConfInfo struct {
	WorkDir string `yaml:"work-dir"`
	HttpServerListen string `yaml:"http-server-listen"`
	BaseUrlPath string `yaml:"base-url-path"`
}

func (ci *ConfInfo) setDefaults() {
	if ci.WorkDir == "" {
		ci.WorkDir = "./"
	}
	if ci.HttpServerListen == "" {
		ci.HttpServerListen = ":9465"
	}
}