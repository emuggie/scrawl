package route

import (
	"fmt"
	"net/http"
)

func setHeaders(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
}

func onError(res http.ResponseWriter) {
	err := recover()
	if err == nil {
		return
	}
	fmt.Println(err)
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte("500"))
	return
}
