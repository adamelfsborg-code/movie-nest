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
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	room, err := u.Data.GetRoomByID(roomID)
	if err != nil {
		fmt.Println("Failed to get user: ", err)
		http.Error(w, "Failed to get user", http.StatusBadRequest)
		return

	}

	jsonBytes, err := json.Marshal(room)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return

	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *RoomHandler) GetRoomInfoByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return

	}

	room, err := u.Data.GetRoomInfoByID(roomID)
	if err != nil {
		fmt.Println("Failed to get user: ", err)
		http.Error(w, "Failed to get user", http.StatusBadRequest)
		return

	}

	jsonBytes, err := json.Marshal(room)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *RoomHandler) GetAvailableUsers(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return

	}

	idParam = r.Header.Get("X-UserID")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return

	}

	searchTerm := r.URL.Query().Get("searchTerm")
	excludeSelfParm := r.URL.Query().Get("excludeSelf")
	excludeExistingParam := r.URL.Query().Get("excludeExisting")

	excludeSelf, err := strconv.ParseBool(excludeSelfParm)
	if err != nil {
		fmt.Println("Failed to parse excludeSelf: ", err)
		fmt.Println("Using excludeSelf as default: true ")
		excludeSelf = true
	}

	excludeExisting, err := strconv.ParseBool(excludeExistingParam)
	if err != nil {
		fmt.Println("Failed to parse excludeExisting: ", err)
		fmt.Println("Using excludeExisting as default: true ")
		excludeExisting = true
	}

	availableUsers := u.Data.GetAvailableUsers(roomID, userID, searchTerm, excludeSelf, excludeExisting)

	jsonBytes, err := json.Marshal(availableUsers)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	idParam := r.Header.Get("X-UserID")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	var body struct {
		Name string `json:"name"`
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	room := data.NewRoom(body.Name)
	err = u.Data.CreateRoom(*room, userID)
	if err != nil {
		fmt.Println("Failed to create room: ", err)
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(map[string]string{"message": "Room created"})
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (u *RoomHandler) AddUserToRoom(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RoomID uuid.UUID `json:"room_id"`
		UserID uuid.UUID `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	room := data.NewRoomUser(body.RoomID, body.UserID)
	err = u.Data.AddUserToRoom(*room)
	if err != nil {
		fmt.Println("Failed to add user to room: ", err)
		http.Error(w, "Failed to add user to room", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(map[string]string{"message": "User added to room"})
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (u *RoomHandler) ListRoomsWithUsers(w http.ResponseWriter, r *http.Request) {
	rooms := u.Data.ListRoomsWithUsers()

	jsonBytes, err := json.Marshal(rooms)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
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
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	room := u.Data.GetRoomWithUsersByID(roomID)

	jsonBytes, err := json.Marshal(room)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *RoomHandler) GetUserRoomsByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.Header.Get("X-UserID")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	rooms := u.Data.GetUserRoomsByID(userID)

	jsonBytes, err := json.Marshal(rooms)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *RoomHandler) GetRoomAccess(w http.ResponseWriter, r *http.Request) {
	idParam := r.Header.Get("X-UserID")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	idParam = chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	access, err := u.Data.GetRoomAccess(roomID, userID)
	if err != nil {
		fmt.Println("Failed to get access: ", err)
		http.Error(w, "Failed to get access", http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(access)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
