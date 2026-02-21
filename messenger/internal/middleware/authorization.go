package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/tousart/messenger/internal/usecase"
)

const (
	ContextKeyAuthMetadata = "metadata"
)

func Authorization(usersService usecase.UsersService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: под заголовок

			cookie, err := r.Cookie("session_id")
			if err != nil {
				log.Printf("authorization error: get cookie: %v\n", err)
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}

			userPayload, err := usersService.ValidateSessionID(context.Background(), cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				log.Printf("authorization error: validate session id: %v\n", err)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyAuthMetadata, userPayload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
