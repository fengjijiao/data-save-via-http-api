package httpLib

import (
	"github.com/fengjijiao/data-save-via-http-api/pkg/coreLib"
	"github.com/fengjijiao/data-save-via-http-api/pkg/sqliteLib"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
	"fmt"
)

const (
	ErrorUnsupportedRequestMethod  = "Unsupported request method"
	ErrorParsingData               = "Error parsing Data"
	ErrorAccountLogin              = "Login Info Error"
	ErrorNoAccessPermission        = "No Access Permission"
	ErrorMissingRequiredParameters = "Missing Required Parameters"
	ErrorContentNoFound 		   = "Content No Found"
)

const (
	StatusSuccess = iota
	StatusFailure
	StatusError
	StatusUnknown
)

type CallbackDataSourceList struct {
	Status int               `json:"status"`
	Data   []sqliteLib.DataSource `json:"data"`
}

type CallbackId struct {
	Status int               `json:"status"`
	Id     int64             `json:"id"`
}

type CallbackDataSetData struct {
	ValueType int    `json:"valueType"`
	Value     interface{} `json:"value"`
	UpdatedTimestamp int64 `json:"updatedTimestamp"`
}

type CallbackDataSet struct {
	Status int               `json:"status"`
	Data   CallbackDataSetData `json:"data"`
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
	ExpireDate   int64    `json:"expireDate"`
}

type CallbackLogin struct {
	Status int               `json:"status"`
	Data   CallbackLoginData `json:"data"`
}

type CallbackTip struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
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

//用户登录
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
				userInfo, err := sqliteLib.GetUserInfoViaUsername(username)
				if err != nil {
					json.NewEncoder(w).Encode(GenError(err.Error()))
					return
				}
				json.NewEncoder(w).Encode(CallbackLogin {
					Status: StatusSuccess,
					Data: CallbackLoginData {
						ExpireDate:   time.Now().Unix(),
						SessionToken: userInfo.Token ,
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

//注册用户
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

//获取数据值(text)
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
	dsid, err := sqliteLib.GetDSIdViaDid(did)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	result, err := sqliteLib.GetDataSet(dsid)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	fmt.Fprintf(w, result.Value)
}

//获取数据值(json)
func GetValJsonHttpHandler(w http.ResponseWriter, r *http.Request) {
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
	if !sqliteLib.ExistDataSource(did) {
		json.NewEncoder(w).Encode(GenError(ErrorContentNoFound))
		return
	}
	dataSource, err := sqliteLib.GetDataSource(did)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	result, err := sqliteLib.GetDataSet(dataSource.Did)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	dataVal, err := coreLib.ValTypeConv(result.Value, result.ValueType)
	if err != nil {
		dataVal = result.Value
	}
	json.NewEncoder(w).Encode(CallbackDataSet {
		Status: StatusSuccess,
		Data: CallbackDataSetData {
			ValueType: result.ValueType,
			Value: dataVal,
			UpdatedTimestamp: dataSource.UpdatedTimestamp,
		},
	})
}

//创建数据源
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
	valtypeStr := r.FormValue("valueType")
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

//更新数据值
func PutValHttpHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := r.ParseMultipartForm(parseMultipartFormMaxMemory); err != nil {
		json.NewEncoder(w).Encode(GenError(ErrorParsingData))
		return
	}
	newValStr := r.FormValue("newValue")
	if newValStr == "" {
		json.NewEncoder(w).Encode(GenError(ErrorMissingRequiredParameters))
		return
	}
	if !sqliteLib.CheckDataSetPermission(uid, did) {
		json.NewEncoder(w).Encode(GenError(ErrorNoAccessPermission))
		return
	}
	if !sqliteLib.ExistDataSource(did) {
		json.NewEncoder(w).Encode(GenError(ErrorContentNoFound))
		return
	}
	dataSource, err := sqliteLib.GetDataSource(did)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	dsid, err := sqliteLib.GetDSIdViaDid(did)
	if err != nil {
		dsid, err = sqliteLib.AddDataSet(did, dataSource.ValueType, newValStr)
		if err != nil {
			json.NewEncoder(w).Encode(GenError(err.Error()))
			return
		}
	}
	if err = sqliteLib.UpdateDataSetValue(dsid, newValStr); err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	if err = sqliteLib.UpdateUpdatedTimestamp(did); err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(CallbackTip {
		Status: StatusSuccess,
		Msg: "Update Successfully",
	})
}

//删除数据源与所有数据
func DelValHttpHandler(w http.ResponseWriter, r *http.Request) {
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
	if !sqliteLib.ExistDataSource(did) {
		json.NewEncoder(w).Encode(GenError(ErrorContentNoFound))
		return
	}
	if err = sqliteLib.DelDataSource(did); err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	if err = sqliteLib.DelDataSetViaDid(did); err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(CallbackTip {
		Status: StatusSuccess,
		Msg: "Delete Successfully",
	})
}

//列出所有数据源
func ListDataSetHttpHandler(w http.ResponseWriter, r *http.Request) {
	uid, err := CheckToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	result, err := sqliteLib.GetDataSourcesViaUid(uid)
	if err != nil {
		json.NewEncoder(w).Encode(GenError(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(CallbackDataSourceList {
		Status: StatusSuccess,
		Data: *result,
	})
}