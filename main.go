package main

import (
	httpLib "data-save-via-http-api/pkg/httpLib"
	sqliteLib "data-save-via-http-api/pkg/sqliteLib"
	"os"
)

func init() {
	sqliteLib.SetDBPath("test.db")
	sqliteLib.Init()
}

func main() {
	Close := make(chan int)
	//sqliteLib.InitTable()
	go httpLib.Run(Close)
	for {
		select {
			case k := <- Close:
				os.Exit(k)
			break
		}
	}
}