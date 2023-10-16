package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adamelfsborg-code/movie-nest/data"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	Data data.UserData
}

func (u *UserHandler) SelectUsers(w http.ResponseWriter, r *http.Request) {
	users := u.Data.List()

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := data.NewRegisterUser(body.Name, body.Name)
	if err != nil {
		fmt.Println("Failed to create user: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.Data.Register(*user)
	if err != nil {
		fmt.Println("Failed to create user: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := u.Data.Login(body.Name, body.Name)
	if err != nil {
		fmt.Println("Failed to login user: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := map[string]string{"token": token}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (u *UserHandler) GetUserInfoByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.Header.Get("X-UserID")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := u.Data.GetUserInfoByID(userID)
	if err != nil {
		fmt.Println("Failed to login user: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (u *UserHandler) HandleUserAccess(w http.ResponseWriter, r *http.Request) {
	idParam := r.Header.Get("X-UserID")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exists := u.Data.CheckUserExistsByID(userID)

	if !exists {
		fmt.Println("User not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *UserHandler) GetUsersInRoom(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users := u.Data.GetUsersInRoom(roomID)

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
