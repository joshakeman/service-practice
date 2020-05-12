package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

func health(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return json.NewEncoder(w).Encode(status)
}
