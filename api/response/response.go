package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

var err error
var output []byte

func JSON(w http.ResponseWriter, responseStatus string, httpStatus int, message, data interface{}) {
	//fmt.Println(responseStatus, httpStatus, message, data)
	output, err = json.Marshal(Response{responseStatus, httpStatus, message, data})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		output = []byte(`{"status":"error", "code":500,"message":"failed to marshal JSON in response.JSON()","data":null}`)
		w.Write(output)
		return
	}
	w.WriteHeader(httpStatus)
	w.Write(output)
}

func Error(w http.ResponseWriter, httpStatus int, err error) {
	JSON(w, "error", httpStatus, err.Error(), nil)
}

func Success(w http.ResponseWriter, message string, httpStatus int, data interface{}) {
	JSON(w, "success", httpStatus, message, data)
}
