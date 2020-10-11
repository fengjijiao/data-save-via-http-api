package httpLib

import (
    "fmt"
    "net/http"
)

func Run(Close chan int) {
	router := NewRouter()
	err := http.ListenAndServe(listenAddress, router)
	if (err != nil) {
		fmt.Printf("SqliteLib.Run Error: %s\n", err)
		Close <- 1
	}
}