package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-app/models"

	"github.com/golang-jwt/jwt"

	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

type UserController struct {
	session *gocql.Session
}

func NewUserController(session *gocql.Session) *UserController {
	return &UserController{session: session}
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("Error decoding credentials: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := models.GetUserByEmail(c.session, credentials.Email)
	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials (database error)")
		return
	}

	if !user.ValidatePassword(credentials.Password) {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials (password mismatch)")
		return
	}

	token, err := generateJWT(user.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status: "success",
		Data: map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":       user.UserID,
				"email":    user.Email,
				"username": user.Username,
			},
		},
	})
}

func generateJWT(userID gocql.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte("your-secret-key"))
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user := models.NewUser(userData.Username, userData.Email, userData.Password)

	if err := user.Create(c.session); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{
		Status: "success",
		Data:   user,
	})
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := gocql.ParseUUID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := models.GetUserByID(c.session, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status: "success",
		Data:   user,
	})
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameters
	params := mux.Vars(r)
	id, err := gocql.ParseUUID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Delete user from database
	if err := models.DeleteUserByID(c.session, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	// Return success response
	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "User deleted successfully",
	})
}
