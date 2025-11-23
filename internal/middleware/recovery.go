package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/pkg/logger"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", logger.WithFields(map[string]interface{}{
					"error": err,
					"stack": string(debug.Stack()),
					"path":  r.URL.Path,
				}))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(model.ErrorResponse("INTERNAL_ERROR", "An unexpected error occurred"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
