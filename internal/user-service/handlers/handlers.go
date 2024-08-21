package handlers

import (
	"encoding/json"
	"net/http"
	"user/pkg/utils"

	"github.com/rs/zerolog/log"
)

type service interface {
	SignUp(string, string) error
}

type Handlers struct {
	S service
}

type SignUpRequestBody struct {
	Email           string
	Password        string
	ConfirmPassword string
}

func (h *Handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpBody SignUpRequestBody
	err := json.NewDecoder(r.Body).Decode(&signUpBody)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	err = utils.SamePasswordVerification(signUpBody.Password, signUpBody.ConfirmPassword)
	if err != nil {
		http.Error(w, "Passwords don`t match", http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(signUpBody.Password)
	if err != nil {
		http.Error(w, "Error creating password", http.StatusInternalServerError)
		return
	}

	err = h.S.SignUp(hash, signUpBody.Email)
	if err != nil {
		http.Error(w, "Error to create user", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
}
