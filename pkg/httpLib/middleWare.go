package httpLib

import (
	"../../pkg/sqliteLib"
	"net/http"
)

func CheckToken(w http.ResponseWriter, r *http.Request) (int64, error) {
	token := r.Header.Get("Token")
	uid, err := sqliteLib.GetUidViaToken(token)
	if err != nil {
		return -1, err
	}
	return uid, nil
}