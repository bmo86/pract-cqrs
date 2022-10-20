package means

import (
	"encoding/json"
	"net/http"
)

func ErrRes(code int, msgError string, w http.ResponseWriter) {
	fields := make(map[string]interface{})

	fields["status"] = "error"
	fields["message"] = msgError

	message, err := json.Marshal(fields)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error ocurred internally"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(message)
}

func SuccessRes(code int, msg string, w http.ResponseWriter) {
	fields := make(map[string]interface{})

	fields["status"] = "success"
	fields["message"] = msg

	message, err := json.Marshal(fields)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("an error ocurred internally"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(message)
}
