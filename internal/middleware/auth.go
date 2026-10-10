package middleware

import (
	"context"
	"database/sql"
	"net/http"
)

type contextKey string

const MerchantIDKey contextKey = "merchant_id"

// AuthMiddleware validates X-Api-Key and X-Secret-Key headers against the database
func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-Api-Key")
			secretKey := r.Header.Get("X-Secret-Key")

			// 1. Ensure required security headers are present
			if apiKey == "" || secretKey == "" {
				http.Error(w, `{"error": "X-Api-Key and X-Secret-Key headers are required"}`, http.StatusUnauthorized)
				return
			}

			// 2. Validate API Key and Secret Key against the database
			query := `SELECT id FROM merchants WHERE api_key = $1 AND secret_key = $2 LIMIT 1`
			var merchantID string

			err := db.QueryRowContext(r.Context(), query, apiKey, secretKey).Scan(&merchantID)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, `{"error": "Invalid API Key or Secret Key"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error": "Failed to process authentication"}`, http.StatusInternalServerError)
				return
			}

			// 3. Inject authenticated merchant_id into request context
			ctx := context.WithValue(r.Context(), MerchantIDKey, merchantID)

			// 4. Forward request to the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
