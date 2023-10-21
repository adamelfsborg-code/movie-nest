package server

import (
	"net/http"
	"strings"

	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/data"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromRequest(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.Env.SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-UserID", userID)

		next.ServeHTTP(w, r)
	})
}

func CustomAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Authenticate(next)
	}
}

func AccessRoom(next http.Handler, data data.RoomData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.Header.Get("X-UserID")

		userID, err := uuid.Parse(idParam)
		if err != nil {
			http.Error(w, "Permission Not Allowd", http.StatusForbidden)
			return
		}

		idParam = chi.URLParam(r, "room_id")

		roomID, err := uuid.Parse(idParam)
		if err != nil {
			http.Error(w, "Permission Not Allowd", http.StatusForbidden)
			return
		}

		permission, err := data.GetRoomAccess(roomID, userID)

		if !permission {
			http.Error(w, "Permission Not Allowd", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CustomAccessRoomMiddleware(data data.RoomData) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return AccessRoom(next, data)
	}
}

func AccessShelf(next http.Handler, data data.ShelfData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.Header.Get("X-UserID")

		userID, err := uuid.Parse(idParam)
		if err != nil {
			http.Error(w, "Permission Not Allowd", http.StatusForbidden)
			return
		}

		idParam = chi.URLParam(r, "shelf_id")

		shelfID, err := uuid.Parse(idParam)
		if err != nil {
			http.Error(w, "Permission Not Allowd", http.StatusForbidden)
			return
		}

		permission, err := data.GetShelfAccess(shelfID, userID)

		if !permission {
			http.Error(w, "Permission Not Allowd", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CustomAccessShelfMiddleware(data data.ShelfData) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return AccessShelf(next, data)
	}
}

func extractTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	return ""
}
