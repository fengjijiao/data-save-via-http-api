package httpLib

import (
    "fmt"
    "net/http"
	"github.com/gorilla/mux"
)

func CssStaticHttpHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filePath := fmt.Sprintf("css/%s.css",vars["cssName"])
	w.Header().Set("Content-Type", GetFileContentType(filePath))
	w.Write(ReadFile(filePath))
}

func JsStaticHttpHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filePath := fmt.Sprintf("js/%s.js",vars["jsName"])
	w.Header().Set("Content-Type", GetFileContentType(filePath))
	w.Write(ReadFile(filePath))
}

func IndexHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(ReadFile("html/index.html"))
}

func LoginHttpHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method == "GET") {
		w.Write(ReadFile("html/login.html"))
	}else if(r.Method == "POST") {
		
	}
}

func RegisterHttpHandler(w http.ResponseWriter, r *http.Request) {
    w.Write(ReadFile("register.html"))
}

func GetValHttpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func PostValHttpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func CreateDataSetHttpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func PutValHttpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func DelValHttpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}