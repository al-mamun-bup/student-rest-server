package auth

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// Map storing valid username-password pairs (for demo purposes)
var validUsers = map[string]string{
	"admin": "password123",
	"user1": "pass1",
	"user2": "pass2",
}

// BasicAuthMiddleware ensures only authenticated users access certain routes
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Format: "Basic base64(username:password)"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username, password := credentials[0], credentials[1]

		// Validate credentials
		if validPassword, exists := validUsers[username]; !exists || validPassword != password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If authentication succeeds, pass request to next handler
		next.ServeHTTP(w, r)
	})
}
