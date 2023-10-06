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
	router.Route("/rooms", a.loadRoomRoutes)
	router.Route("/movies", a.loadMovieRoutes)

	a.router = router
}

func (a *Server) loadUserRoutes(router chi.Router) {
	userHandler := &handlers.UserHandler{
		Data: data.UserData{},
	}

	router.Get("/", userHandler.SelectUsers)
	router.Get("/{user_id}", userHandler.SelectUsers)
	router.Post("/", userHandler.CreateUser)
}

func (a *Server) loadRoomRoutes(router chi.Router) {
	roomHandler := &handlers.RoomHandler{
		Data: data.RoomData{},
	}

	router.Get("/", roomHandler.SelectRooms)
	router.Get("/{room_id}", roomHandler.GetRoomByID)
	router.Get("/withusers", roomHandler.ListRoomsWithUsers)
	router.Get("/withusers/{room_id}", roomHandler.GetRoomWithUsersByID)
	router.Get("/withusers/{user_id}", roomHandler.GetUserRoomsByID)

	router.Post("/", roomHandler.CreateRoom)
	router.Post("/adduser", roomHandler.AddUserToRoom)
}

func (a *Server) loadMovieRoutes(router chi.Router) {
	movieHandler := &handlers.MovieHandler{
		Data: data.MovieData{
			Env: a.config,
		},
	}

	router.Get("/{movie_id}", movieHandler.GetMovie)
}
