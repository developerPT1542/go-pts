package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pts/mdes/service/auth"
	"github.com/pts/mdes/types"
	"github.com/pts/mdes/utils"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService types.UserService
}

func NewHandler(service types.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func (handler *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", handler.HandleLogin).Methods("POST")
	router.HandleFunc("/register", handler.HandleRegister).Methods("POST")
}

func (handler *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Getting JSON Payload
	// Check if the user exists
	// If it doesn't we create the new User
}
func (handler *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// Get JSON Payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validating JSON Payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload &v", errors))
		return
	}

	// Check if the user exists
	_, err := handler.userService.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists..", payload.Email))
		return
	}

	hashedPwd, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// If it doesn't we create the new user
	err = handler.userService.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPwd,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
