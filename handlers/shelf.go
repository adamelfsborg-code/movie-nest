package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adamelfsborg-code/movie-nest/data"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ShelfHandler struct {
	Data data.ShelfData
}

func (s *ShelfHandler) CreateShelf(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name   string    `json:"name"`
		RoomID uuid.UUID `json:"room_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	shelf := data.NewShelf(body.Name, body.RoomID)
	err = s.Data.CreateShelf(*shelf)
	if err != nil {
		fmt.Println("Failed to create shelf: ", err)
		http.Error(w, "Failed to create shelf", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(map[string]string{"message": "Shelf created"})
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (s *ShelfHandler) GetShelvesByRoomID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	shelf := s.Data.GetShelvesByRoomID(roomID)
	if err != nil {
		fmt.Println("Failed to get shelf: ", err)
		http.Error(w, "Failed to get shelf", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(shelf)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (s *ShelfHandler) GetShelfMoviesByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "shelf_id")

	shelfID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	shelf := s.Data.GetShelfMoviesByID(shelfID)
	if err != nil {
		fmt.Println("Failed to get shelf: ", err)
		http.Error(w, "Failed to get shelf", http.StatusInternalServerError)
		return

	}

	jsonBytes, err := json.Marshal(shelf)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (s *ShelfHandler) GetShelfInfoByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "shelf_id")

	shelfID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	shelf := s.Data.GetShelfInfoByID(shelfID)
	if err != nil {
		fmt.Println("Failed to get shelf: ", err)
		http.Error(w, "Failed to get shelf", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(shelf)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (s *ShelfHandler) GetAvailableMovies(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "shelf_id")

	shelfID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	searchTerm := r.URL.Query().Get("searchTerm")
	excludeExistingParam := r.URL.Query().Get("excludeExisting")

	excludeExisting, err := strconv.ParseBool(excludeExistingParam)
	if err != nil {
		fmt.Println("Failed to parse excludeExisting: ", err)
		fmt.Println("Using excludeExisting as default: true ")
		excludeExisting = true
	}

	availableMovies, err := s.Data.GetAvailableMovies(shelfID, searchTerm, excludeExisting)
	if err != nil {
		fmt.Println("Failed to search movies: ", err)
		http.Error(w, "Failed to search movies", http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(availableMovies)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
