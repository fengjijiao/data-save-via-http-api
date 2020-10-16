package httpLib

import (
	sqliteLib "../../pkg/sqliteLib"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	ERROR_UNSUPPORTED_REQUEST_METHOD = "Unsupported request method"
	ERROR_PARSING_DATA = "Error parsing Data"
	ERROR_ACCOUNT_LOGIN = "Login Info Error"
)

const (
	StatusSuccess = iota
	StatusFailure
	StatusError
	StatusUnknown
)

type CallbackRegisterData struct {
	Uid int64 `json:"uid"`
}

type CallbackRegister struct {
	Status int               `json:"status"`
	Data   CallbackRegisterData `json:"data"`
}

type CallbackLoginData struct {
	SessionToken string `json:"sessionToken"`
	ExpireDate   int    `json:"expireDate"`
}

type CallbackLogin struct {
	Status int               `json:"status"`
	Data   CallbackLoginData `json:"data"`
}

type CallbackTip struct {
	Status int `json:"status"`
	Msg string `json:"msg"`
}

func GenError(msg string) *CallbackTip {
	return &CallbackTip{
		Status: StatusError,
		Msg:    msg,
	}
}

func IndexHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{}"))
}

func LoginHttpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "POST":
			if err := r.ParseMultipartForm(parseMultipartFormMaxMemory); err != nil {
				json.NewEncoder(w).Encode(GenError(ERROR_PARSING_DATA))
				return
			}
			username := r.FormValue("username")
			password := r.FormValue("password")
			if sqliteLib.AuthUser(username, password) {
				json.NewEncoder(w).Encode(CallbackLogin {
					Status: StatusSuccess,
					Data: CallbackLoginData {
						ExpireDate:   1111111111,
						SessionToken: username + password ,
					},
				})
			}else {
				json.NewEncoder(w).Encode(GenError(ERROR_ACCOUNT_LOGIN))
			}
		default:
			json.NewEncoder(w).Encode(GenError(ERROR_UNSUPPORTED_REQUEST_METHOD))
			return
	}
}

func RegisterHttpHandler(w http.ResponseWriter, r *http.Request) {
    //w.Write(ReadFile("register.html"))
	switch r.Method {
	case "POST":
		if err := r.ParseMultipartForm(parseMultipartFormMaxMemory); err != nil {
			json.NewEncoder(w).Encode(GenError(ERROR_PARSING_DATA))
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		uid, err := sqliteLib.AddUser(username, password)
		if  err != nil {
			json.NewEncoder(w).Encode(GenError(err.Error()))
		}else {
			json.NewEncoder(w).Encode(CallbackRegister {
				Status: StatusSuccess,
				Data: CallbackRegisterData {
					Uid: uid,
				},
			})
		}

	default:
		json.NewEncoder(w).Encode(GenError(ERROR_UNSUPPORTED_REQUEST_METHOD))
		return
	}
}

func GetValHttpHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    fmt.Fprint(w, "Welcome!\n" + vars["did"])
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