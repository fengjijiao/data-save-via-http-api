package main

import (
	"github.com/fengjijiao/data-save-via-http-api/pkg/commonio"
	"github.com/fengjijiao/data-save-via-http-api/pkg/conf"
	"github.com/fengjijiao/data-save-via-http-api/pkg/logio"
	"github.com/fengjijiao/data-save-via-http-api/pkg/coreLib"
	"github.com/fengjijiao/data-save-via-http-api/pkg/httpLib"
	"github.com/fengjijiao/data-save-via-http-api/pkg/sqliteLib"
	"flag"
	"os"
	"go.uber.org/zap"
)

var (
	configFilePath string
	flagQuiet bool
	Close chan int
)

func init() {
	Close = make(chan int)
	
	logio.Init()
	
	flag.StringVar(&configFilePath, "c", "config.yaml", "the path of configuration file")
	flag.BoolVar(&flagQuiet, "quiet", false, "quiet for log print.")
	
	flag.Parse()
	
	if(flagQuiet) {
		logio.Cfg.Level.SetLevel(zap.ErrorLevel)
	}
	
	if !coreLib.IsFile(configFilePath) || !coreLib.ExistsFile(configFilePath) {
		logio.Logger.Fatal("configure file not found!")
	}
	
	dat, err := commonio.ReadFile(configFilePath)
	if err != nil {
		logio.Logger.Fatal("read configure file fail!", zap.Error(err))
	}
	
	err = conf.Load(dat)
	if err != nil {
		logio.Logger.Fatal("parse configure file fail!", zap.Error(err))
	}
}

func main() {
	sqliteLib.Init()
	sqliteLib.InitTable()
	go httpLib.Run(Close)
	for {
		select {
			case k := <- Close:
				os.Exit(k)
		}
	}
}