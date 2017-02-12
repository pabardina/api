package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func httpError(w http.ResponseWriter, code int, msg, description string) {
	writeJSON(w, struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}{msg, description}, code)
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) error {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)

	return nil
}

func GetQueryParamToStr(param string, r *http.Request) int {
	if val := r.FormValue(param); val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		return intVal
	}
	return 0

}
