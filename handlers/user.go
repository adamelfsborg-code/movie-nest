package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adamelfsborg-code/movie-nest/data"
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
