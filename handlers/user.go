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

type UserHandler struct {
	Data data.UserData
}

func (u *UserHandler) SelectUsers(w http.ResponseWriter, r *http.Request) {
	users := u.Data.List()

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
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
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	user, err := data.NewRegisterUser(body.Name, body.Password)
	if err != nil {
		fmt.Println("Failed to create user: ", err)
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		return
	}

	err = u.Data.Register(*user)
	if err != nil {
		fmt.Println("Failed to create user: ", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(map[string]string{"message": "Registerd"})
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	token, err := u.Data.Login(body.Name, body.Password)
	if err != nil {
		fmt.Println("Failed to login user: ", err)
		http.Error(w, "Failed to login user", http.StatusBadRequest)
		return
	}

	response := map[string]string{"token": token}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Failed to encode token: ", err)
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
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	user := u.Data.GetUserInfoByID(userID)
	if err != nil {
		fmt.Println("Failed to login user: ", err)
		http.Error(w, "Failed to login user", http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Failed to encode user: ", err)
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
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}
	exists := u.Data.CheckUserExistsByID(userID)

	if !exists {
		fmt.Println("User not found")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(map[string]string{"message": "Access Allowed"})
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *UserHandler) GetUsersInRoom(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
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

	excludeSelfParm := r.URL.Query().Get("excludeSelf")

	excludeSelf, err := strconv.ParseBool(excludeSelfParm)
	if err != nil {
		fmt.Println("Failed to parse excludeSelf: ", err)
		fmt.Println("Using excludeSelf as default: true ")
		excludeSelf = true
	}

	users := u.Data.GetUsersInRoom(roomID, userID, excludeSelf)

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
