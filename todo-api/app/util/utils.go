package util

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetParam(r *http.Request, key string) string {
	vars := mux.Vars(r)
	return vars[key]
}
