package httpLib

import (
	"github.com/fengjijiao/data-save-via-http-api/pkg/logio"
	"github.com/fengjijiao/data-save-via-http-api/pkg/conf"
    "net/http"
	"go.uber.org/zap"
)

func Run(Close chan int) {
	router := NewRouter()
	err := http.ListenAndServe(conf.Config.HttpServerListen, router)
	if (err != nil) {
		logio.Logger.Error("SqliteLib.Run Error: ", zap.Error(err))
		Close <- 1
	}
}