package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adamelfsborg-code/movie-nest/data"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type RoomHandler struct {
	Data data.RoomData
}

func (u *RoomHandler) SelectRooms(w http.ResponseWriter, r *http.Request) {
	rooms := u.Data.ListRooms()

	jsonBytes, err := json.Marshal(rooms)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *RoomHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room, err := u.Data.GetRoomByID(roomID)
	if err != nil {
		fmt.Println("Failed to get user: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(room)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room := data.NewRoom(body.Name)
	err = u.Data.CreateRoom(*room)
	if err != nil {
		fmt.Println("Failed to create room: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (u *RoomHandler) AddUserToRoom(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RoomID uuid.UUID `json:"room_id"`
		UserID uuid.UUID `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room := data.NewRoomUser(body.RoomID, body.UserID)
	err = u.Data.AddUserToUser(*room)
	if err != nil {
		fmt.Println("Failed to add user to room: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *RoomHandler) ListRoomsWithUsers(w http.ResponseWriter, r *http.Request) {
	rooms := u.Data.ListRoomsWithUsers()

	jsonBytes, err := json.Marshal(rooms)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *RoomHandler) GetRoomWithUsersByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room := u.Data.GetRoomWithUsersByID(roomID)

	jsonBytes, err := json.Marshal(room)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *RoomHandler) GetUserRoomsByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "user_id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rooms := u.Data.GetUserRoomsByID(userID)

	jsonBytes, err := json.Marshal(rooms)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
