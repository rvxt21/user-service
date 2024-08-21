package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"user/internal/user-service/storage"
	"user/pkg/utils"

	"github.com/rs/zerolog/log"
)

type service interface {
	SignUp(string, string, string) error
	SignIn(email string) (string, error)
}

type Handlers struct {
	S service
}

type SignUpRequestBody struct {
	Email           string
	Password        string
	ConfirmPassword string
	Login           string
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

	err = h.S.SignUp(hash, signUpBody.Email, signUpBody.Login)
	if err != nil {
		if errors.Is(err, storage.ErrEmailOrLoginAlreadyExists) {
			http.Error(w, "Email or login already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type LoginReqBody struct {
	Email    string
	Password string
}

func (h *Handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	var login LoginReqBody
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	hash, err := h.S.SignIn(login.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Debug().Err(err).Msgf("%s", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	ok := utils.VerifyPassword(login.Password, hash)
	if !ok {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
