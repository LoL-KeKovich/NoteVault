package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/LoL-KeKovich/NoteVault/internal/dto"
	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"github.com/LoL-KeKovich/NoteVault/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserID string

const (
	userIDKey UserID = "user_id"
)

type UserService struct {
	DBClient repository.UserRepo
}

func (srv UserService) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	response := dto.RegisterResponse{}
	var registerReq dto.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Wrong request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if registerReq.Email == "" || registerReq.Password == "" {
		response.Error = "Email and password are required"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", slog.String("error", err.Error()))
		response.Error = "Internal server error"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	user := model.User{
		Email:        registerReq.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    registerReq.FirstName,
		LastName:     registerReq.LastName,
	}

	res, err := srv.DBClient.RegisterUser(user)
	if err != nil {
		slog.Error("Failed to create user", slog.String("error", err.Error()))
		response.Error = "Failed to create user"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("User registered")
	response.Data = res
	json.NewEncoder(w).Encode(response)
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

	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte("placeholder_secret_key") //В будущем создать нормальный ключ в конфиге

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		slog.Error("Failed to generate JWT", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, //На случай https
		SameSite: http.SameSiteLaxMode,
	})

	slog.Info("User logged in", slog.String("email", user.Email))
	response.Data = user
	json.NewEncoder(w).Encode(response)
}

func (srv UserService) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	response := dto.LoginResponse{}

	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok || userID == "" {
		slog.Error("UserID not found in context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := srv.DBClient.GetProfile(userID)
	if err != nil {
		slog.Error("Failed to get user", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("User checked")
	response.Data = user
	json.NewEncoder(w).Encode(response)
}

func (srv UserService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			slog.Error("Cookie 'auth_token' not found", "error", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte("placeholder_secret_key"), nil //В будущем создать нормальный ключ в конфиге
		})
		if err != nil || !token.Valid {
			slog.Error("Invalid token", "error", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Error("Invalid token claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			slog.Error("user_id is not a string", "type", fmt.Sprintf("%T", claims["user_id"]))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
