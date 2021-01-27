package main

import (
	"github.com/fengjijiao/data-save-via-http-api/pkg/coreLib"
	"github.com/fengjijiao/data-save-via-http-api/pkg/httpLib"
	"github.com/fengjijiao/data-save-via-http-api/pkg/sqliteLib"
	"flag"
	"fmt"
	"github.com/jinzhu/configor"
	"os"
)

var ConfigPATH string
var Config = struct {
	DB struct {
		FilePath string `default:"test.db"`
	}
	HTTP struct {
		Listen string `default:":8090"`
	}
}{}

func init() {
	flag.StringVar(&ConfigPATH, "c", "config.yaml", "the path of configuration file")
	flag.Parse()
	if !coreLib.IsFile(ConfigPATH) || !coreLib.ExistsFile(ConfigPATH) {
		fmt.Printf("Failed to find configuration from %s\n", ConfigPATH)
		os.Exit(3)
		return
	}
	if err := configor.Load(&Config, ConfigPATH); err != nil {
		fmt.Println(err)
		os.Exit(3)
		return
	}
	sqliteLib.SetDBPath(Config.DB.FilePath)
	httpLib.SetListenAddress(Config.HTTP.Listen)
}

func main() {
	Close := make(chan int)
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