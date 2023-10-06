package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adamelfsborg-code/movie-nest/data"
	"github.com/go-chi/chi/v5"
)

type MovieHandler struct {
	Data data.MovieData
}

func (m *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "movie_id")

	movieID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		fmt.Println("Failed to convert param to uint: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movie, err := m.Data.GetMovie(uint(movieID))

	jsonBytes, err := json.Marshal(movie)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
