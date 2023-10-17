package server

import (
	"net/http"

	"github.com/adamelfsborg-code/movie-nest/data"
	"github.com/adamelfsborg-code/movie-nest/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (a *Server) loadRoutes() {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Replace with your allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/users", a.loadUserRoutes)

	router.Group(func(r chi.Router) {
		r.Use(CustomAuthMiddleware())
		r.Route("/rooms", a.loadRoomRoutes)
		r.Route("/movies", a.loadMovieRoutes)
		r.Route("/shelves", a.loadShelfRoutes)
	})

	a.router = router
}

func (a *Server) loadUserRoutes(router chi.Router) {
	userHandler := &handlers.UserHandler{
		Data: data.UserData{
			Env: a.config,
			DB:  a.datbase,
		},
	}

	router.Group(func(r chi.Router) {
		r.Use(CustomAuthMiddleware())
		r.Get("/", userHandler.SelectUsers)
		r.Get("/user", userHandler.GetUserInfoByID)
		r.Get("/access", userHandler.HandleUserAccess)
		r.Get("/rooms/{room_id}", userHandler.GetUsersInRoom)
	})

	router.Post("/register", userHandler.Register)
	router.Post("/login", userHandler.Login)
}

func (a *Server) loadRoomRoutes(router chi.Router) {
	roomHandler := &handlers.RoomHandler{
		Data: data.RoomData{
			DB: a.datbase,
		},
	}

	router.Get("/", roomHandler.SelectRooms)
	router.Get("/{room_id}", roomHandler.GetRoomByID)
	router.Get("/{room_id}/info", roomHandler.GetRoomInfoByID)
	router.Get("/{room_id}/available-users", roomHandler.GetAvailableUsers)
	router.Get("/withusers", roomHandler.ListRoomsWithUsers)
	router.Get("/withusers/{room_id}", roomHandler.GetRoomWithUsersByID)
	router.Get("/users", roomHandler.GetUserRoomsByID)

	router.Post("/", roomHandler.CreateRoom)
	router.Post("/users", roomHandler.AddUserToRoom)
}

func (a *Server) loadMovieRoutes(router chi.Router) {
	movieHandler := &handlers.MovieHandler{
		Data: data.MovieData{
			Env: a.config,
			DB:  a.datbase,
		},
	}

	router.Get("/{movie_id}", movieHandler.GetMovie)
	router.Get("/{movie_id}/details", movieHandler.GetMovieDetails)

	router.Post("/", movieHandler.CreateMovie)
	router.Post("/{movie_id}/ratings", movieHandler.RateMovie)
}

func (a *Server) loadShelfRoutes(router chi.Router) {
	shelfHandler := &handlers.ShelfHandler{
		Data: data.ShelfData{
			Env: a.config,
			DB:  a.datbase,
		},
	}

	router.Get("/{shelf_id}/movies", shelfHandler.GetShelfMoviesByID)
	router.Get("/rooms/{room_id}", shelfHandler.GetShelvesByRoomID)
	router.Get("/{shelf_id}/info", shelfHandler.GetShelfInfoByID)
	router.Get("/{shelf_id}/available-movies", shelfHandler.GetAvailableMovies)

	router.Post("/", shelfHandler.CreateShelf)
}
