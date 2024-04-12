package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
	contentType       = "application/json"
)

type errorReponse struct {
	ErrorMsg string `json:"error"`
}

func JSONError(w http.ResponseWriter, code int, msg string, log *slog.Logger) {
	w.WriteHeader(code)
	w.Header().Set(contentTypeHeader, contentType)

	b, err := json.Marshal(errorReponse{ErrorMsg: msg})
	if err != nil {
		log.Error(err.Error())
	}

	if _, err := w.Write(b); err != nil {
		log.Error(err.Error())
	}
}

func JSONWithMarshal(w http.ResponseWriter, code int, body any, log *slog.Logger) {
	w.WriteHeader(code)
	w.Header().Set(contentTypeHeader, contentType)

	b, err := json.Marshal(body)
	if err != nil {
		log.Error(err.Error())
	}
	if _, err := w.Write([]byte(b)); err != nil {
		log.Error(err.Error())
	}
}

func JSON(w http.ResponseWriter, code int, body string, log *slog.Logger) {
	w.WriteHeader(code)
	w.Header().Set(contentTypeHeader, contentType)

	if _, err := w.Write([]byte(body)); err != nil {
		log.Error(err.Error())
	}
}
