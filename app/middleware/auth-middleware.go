package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"zen-test/app/auth"
	"zen-test/app/exceptions"
)

type contextKey string

const userContextKey contextKey = "user"

func isPublicRoute(r *http.Request) bool {
	return (r.URL.Path == "/users/login" || r.URL.Path == "/users/signup") && r.Method == "POST"
}

func RedirectSwagger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger") && r.URL.Path == "/swagger" {
			http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicRoute(r) || strings.HasPrefix(r.URL.Path, "/swagger/") {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// if _, err := r.Cookie(helpers.UserSession); err != nil {
		// 	http.Error(w, "Invalid user cookie", http.StatusBadRequest)
		// 	return
		// }

		ctx := context.WithValue(r.Context(), userContextKey, claims["id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				exceptions.ErrorHandler(w, r, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Fungsi untuk mendapatkan user ID dari context
func GetUserID(r *http.Request) string {
	if userID, ok := r.Context().Value(userContextKey).(string); ok {
		return userID
	}
	return ""
}
