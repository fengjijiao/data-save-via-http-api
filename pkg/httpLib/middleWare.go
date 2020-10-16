package httpLib

import (
	"../../pkg/sqliteLib"
	"errors"
	"net/http"
)

func CheckToken(w http.ResponseWriter, r *http.Request) (int64, error) {
	token := r.Header.Get("Token")
	if token == "" {
		return -1, errors.New(ErrorMissingRequiredParameters)
	}
	uid, err := sqliteLib.GetUidViaToken(token)
	if err != nil {
		return -1, err
	}
	return uid, nil
}