package server

import (
	"net/http"

	"github.com/adamelfsborg-code/movie-nest/data"
	"github.com/adamelfsborg-code/movie-nest/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *Server) loadRoutes() {
	router := chi.NewRouter()

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
		},
	}

	router.Group(func(r chi.Router) {
		r.Use(CustomAuthMiddleware())
		r.Get("/", userHandler.SelectUsers)
		r.Get("/{user_id}", userHandler.SelectUsers)
	})

	router.Post("/register", userHandler.Register)
	router.Post("/login", userHandler.Login)
}

func (a *Server) loadRoomRoutes(router chi.Router) {
	roomHandler := &handlers.RoomHandler{
		Data: data.RoomData{},
	}

	router.Get("/", roomHandler.SelectRooms)
	router.Get("/{room_id}", roomHandler.GetRoomByID)
	router.Get("/withusers", roomHandler.ListRoomsWithUsers)
	router.Get("/withusers/{room_id}", roomHandler.GetRoomWithUsersByID)
	router.Get("/users/{user_id}", roomHandler.GetUserRoomsByID)

	router.Post("/", roomHandler.CreateRoom)
	router.Post("/users", roomHandler.AddUserToRoom)
}

func (a *Server) loadMovieRoutes(router chi.Router) {
	movieHandler := &handlers.MovieHandler{
		Data: data.MovieData{
			Env: a.config,
		},
	}

	router.Get("/{movie_id}", movieHandler.GetMovie)

	router.Post("/", movieHandler.CreateMovie)
}

func (a *Server) loadShelfRoutes(router chi.Router) {
	shelfHandler := &handlers.ShelfHandler{
		Data: data.ShelfData{},
	}

	router.Post("/", shelfHandler.CreateShelf)
}
