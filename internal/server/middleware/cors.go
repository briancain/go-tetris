package middleware

import (
	"net/http"
	"strings"
)

// CORS adds Cross-Origin Resource Sharing headers to allow browser requests
func CORS(allowedOrigins []string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is in allowed list
			allowed := false
			if origin == "" {
				// Allow requests without Origin header (like curl)
				allowed = true
			} else {
				for _, allowedOrigin := range allowedOrigins {
					if origin == allowedOrigin {
						allowed = true
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				if allowed {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}

			next(w, r)
		}
	}
}

// ParseCORSOrigins parses a comma-separated string of origins
func ParseCORSOrigins(originsStr string) []string {
	if originsStr == "" {
		return []string{}
	}

	origins := strings.Split(originsStr, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return origins
}
