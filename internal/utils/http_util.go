package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	ContextUserKey string = "userID"
)

func RespondWithError(resp http.ResponseWriter, code int, message string) {
	error := struct {
		ErrorMessage string `json:"error_message"`
	}{
		message,
	}
	RespondWithJSON(resp, code, error)
}

func RespondWithJSON(resp http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(code)
	resp.Write(response)
}

func UserIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(ContextUserKey)
	id, ok := v.(string)
	return id, ok
}
