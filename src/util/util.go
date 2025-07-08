package util

import (
	"encoding/json"
	"net/http"
	"time"
)

type StandardResponse struct {
	ResponseStatus  int         `json:"responseStatus"`
	TimeDate        string      `json:"timedate"`
	ResponseMessage string      `json:"responseMessage"`
	ResponseData    interface{} `json:"responseData,omitempty"`
}

func ResponseWithJSON(w http.ResponseWriter, data interface{}, message string) {
	res := StandardResponse{
		ResponseStatus:  http.StatusOK,
		TimeDate:        time.Now().Format("2006-01-02 15:04:05"),
		ResponseMessage: message,
		ResponseData:    data,
	}
	responseWithJSON(w, http.StatusOK, res)
}

func ResponseWithError(w http.ResponseWriter, status int, message string) {
	res := StandardResponse{
		ResponseStatus:  status,
		TimeDate:        time.Now().Format("2006-01-02 15:04:05"),
		ResponseMessage: message,
		ResponseData:    nil,
	}
	responseWithJSON(w, status, res)
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
