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

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := data.NewUser(body.Name)
	err = u.Data.Create(*user)
	if err != nil {
		fmt.Println("Failed to create user: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
