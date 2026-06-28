package httpserver

import (
	"encoding/json"
	"net/http"

	pkgerrors "github.com/pkg/errors"
)

const maxBodyBytes = 1 << 20 // 1 MiB

// DecodeJSON decodes the JSON body into dst with a 1 MiB limit and unknown-field rejection.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return pkgerrors.WithStack(err)
	}
	return nil
}

// WriteJSON encodes v as JSON and writes it with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return pkgerrors.WithStack(json.NewEncoder(w).Encode(v))
}
