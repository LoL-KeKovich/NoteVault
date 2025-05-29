package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/LoL-KeKovich/NoteVault/internal/dto"
	"github.com/LoL-KeKovich/NoteVault/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DBClient repository.UserRepo
}

func (srv UserService) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	response := dto.LoginResponse{}
	var loginReq dto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Invalid request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := srv.DBClient.LoginUser(loginReq.Email)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "User not found or wrong password"
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		slog.Error(err.Error())
		response.Error = "User not found or wrong password"
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    user.Email,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	slog.Info("User logged in", slog.String("email", user.Email))
	response.Data = user
	json.NewEncoder(w).Encode(response)
}
