package helpers

import (
	"fmt"
	"net/http"
)

func JSONError(w http.ResponseWriter, error []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}
