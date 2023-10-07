package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adamelfsborg-code/movie-nest/data"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shelf := data.NewShelf(body.Name, body.RoomID)
	err = s.Data.CreateShelf(*shelf)
	if err != nil {
		fmt.Println("Failed to create shelf: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
