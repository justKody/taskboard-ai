package user

import (
	"errors"
	"net/http"

	"github.com/justKody/taskboard-go-api/db/sqlc"
	"github.com/justKody/taskboard-go-api/service/auth"
	"github.com/justKody/taskboard-go-api/utils"
)

func (c *Handler) HandleLogin(w http.ResponseWriter, req *http.Request) {

}

func (c *Handler) HandleSignup(w http.ResponseWriter, req *http.Request) {
	var payload SignupRequestDTO
	if err := utils.ParseJSON(req, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	existedUser, err := c.store.GetUserByEmail(payload.Email)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if existedUser != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("user already exists"))
		return
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// create the user in database

	params := sqlc.CreateUserParams{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	user, err := c.store.CreateUser(params)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, user)

}
