package httpLib

import (
	"../../pkg/sqliteLib"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	ErrorUnsupportedRequestMethod  = "Unsupported request method"
	ErrorParsingData               = "Error parsing Data"
	ErrorAccountLogin              = "Login Info Error"
	ErrorNoAccessPermission        = "No Access Permission"
	ErrorMissingRequiredParameters = "Missing Required Parameters"
)

const (
	StatusSuccess = iota
	StatusFailure
	StatusError
	StatusUnknown
)

type CallbackId struct {
	Status int               `json:"status"`
	Id   int64 `json:"id"`
}

type CallbackDataSet struct {
	Status int               `json:"status"`
	Data   sqliteLib.DataSet `json:"data"`
}

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
	return &CallbackTip {
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
				json.NewEncoder(w).Encode(GenError(ErrorParsingData))
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
				json.NewEncoder(w).Encode(GenError(ErrorAccountLogin))
			}
		default:
			json.NewEncoder(w).Encode(GenError(ErrorUnsupportedRequestMethod))
			return
	}
}

func RegisterHttpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseMultipartForm(parseMultipartFormMaxMemory); err != nil {
			json.NewEncoder(w).Encode(GenError(ErrorParsingData))
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
		json.NewEncoder(w).Encode(GenError(ErrorUnsupportedRequestMethod))
		return
	}
}

func GetValHttpHandler(w http.ResponseWriter, r *http.Request) {
	uid, err := CheckToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	vars := mux.Vars(r)
	did, err := strconv.ParseInt(vars["did"], 10, 64)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	if !sqliteLib.CheckDataSetPermission(uid, did) {
		json.NewEncoder(w).Encode(GenError(ErrorNoAccessPermission))
		return
	}
	result, err := sqliteLib.GetDataSet(did)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(CallbackDataSet {
		Status: StatusSuccess,
		Data: *result,
	})
}

func CreateDataSetHttpHandler(w http.ResponseWriter, r *http.Request) {
	uid, err := CheckToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}

	if err := r.ParseMultipartForm(parseMultipartFormMaxMemory); err != nil {
		json.NewEncoder(w).Encode(GenError(ErrorParsingData))
		return
	}
	valtypeStr := r.FormValue("valtype")

	if valtypeStr == "" {
		json.NewEncoder(w).Encode(GenError(ErrorMissingRequiredParameters))
		return
	}
	valtype, err := strconv.Atoi(valtypeStr)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}

	did, err := sqliteLib.AddDataSource(uid, valtype)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(CallbackId {
		Status: StatusSuccess,
		Id: did,
	})
}

func PutValHttpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func DelValHttpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func ListDataSetHttpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}