package health

import (
	"encoding/json"
	"net/http"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		response := map[string]string{"message": "OK"}
		w.WriteHeader(http.StatusOK)
		bytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(json.RawMessage{})
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bytes)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
