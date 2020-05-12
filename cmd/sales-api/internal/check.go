package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
)

func health(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	if n := rand.Intn(100) % 2; n == 0 {
		return errors.New("test error")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return json.NewEncoder(w).Encode(status)
}
