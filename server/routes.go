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
	roomData := data.RoomData{
		DB: a.datbase,
	}
	userHandler := &handlers.UserHandler{
		Data: data.UserData{
			Env:  a.config,
			DB:   a.datbase,
			Nats: a.nats,
		},
	}

	router.Group(func(r chi.Router) {
		r.Use(CustomAuthMiddleware())
		r.Get("/", userHandler.SelectUsers)
		r.Get("/user", userHandler.GetUserInfoByID)
		r.Get("/access", userHandler.HandleUserAccess)

		r.Group(func(r chi.Router) {
			r.Use(CustomAccessRoomMiddleware(roomData))
			r.Get("/rooms/{room_id}", userHandler.GetUsersInRoom)
		})
	})

	router.Post("/register", userHandler.Register)
	router.Post("/login", userHandler.Login)
}

func (a *Server) loadRoomRoutes(router chi.Router) {
	data := data.RoomData{
		DB:   a.datbase,
		Nats: a.nats,
	}
	roomHandler := &handlers.RoomHandler{
		Data: data,
	}

	router.Get("/", roomHandler.SelectRooms)

	router.Group(func(r chi.Router) {
		r.Use(CustomAccessRoomMiddleware(data))
		r.Get("/{room_id}", roomHandler.GetRoomByID)
		r.Get("/{room_id}/info", roomHandler.GetRoomInfoByID)
		r.Get("/{room_id}/access", roomHandler.GetRoomAccess)
		r.Get("/{room_id}/available-users", roomHandler.GetAvailableUsers)
		r.Get("/withusers/{room_id}", roomHandler.GetRoomWithUsersByID)

		r.Get("/{room_id}/access-allowed", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	router.Get("/withusers", roomHandler.ListRoomsWithUsers)
	router.Get("/users", roomHandler.GetUserRoomsByID)

	router.Post("/", roomHandler.CreateRoom)
	router.Post("/users", roomHandler.AddUserToRoom)

}

func (a *Server) loadMovieRoutes(router chi.Router) {
	movieHandler := &handlers.MovieHandler{
		Data: data.MovieData{
			Env:  a.config,
			DB:   a.datbase,
			Nats: a.nats,
		},
	}

	router.Get("/{movie_id}", movieHandler.GetMovie)
	router.Get("/{movie_id}/details", movieHandler.GetMovieDetails)

	router.Post("/", movieHandler.CreateMovie)
	router.Post("/{movie_id}/ratings", movieHandler.RateMovie)
}

func (a *Server) loadShelfRoutes(router chi.Router) {
	shelfData := data.ShelfData{
		Env:  a.config,
		DB:   a.datbase,
		Nats: a.nats,
	}

	roomData := data.RoomData{
		DB: a.datbase,
	}

	shelfHandler := &handlers.ShelfHandler{
		Data: shelfData,
	}

	router.Group(func(r chi.Router) {
		r.Use(CustomAccessShelfMiddleware(shelfData))
		r.Get("/{shelf_id}/movies", shelfHandler.GetShelfMoviesByID)
		r.Get("/{shelf_id}/info", shelfHandler.GetShelfInfoByID)
		r.Get("/{shelf_id}/available-movies", shelfHandler.GetAvailableMovies)
	})

	router.Group(func(r chi.Router) {
		r.Use(CustomAccessRoomMiddleware(roomData))
		r.Get("/rooms/{room_id}", shelfHandler.GetShelvesByRoomID)
	})

	router.Post("/", shelfHandler.CreateShelf)
}
