package user

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/justKody/taskboard-go-api/db/sqlc"
	"github.com/justKody/taskboard-go-api/service/auth"
	"github.com/justKody/taskboard-go-api/utils"
)

func (c *Handler) HandleLogin(w http.ResponseWriter, req *http.Request) {
	var payload LoginRequestDTO
	if err := utils.ParseJSON(req, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	exisitingUser, err := c.store.GetUserByEmail(payload.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.WriteError(w, http.StatusBadRequest, errors.New("Invalid email or password"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, errors.New("Something went wrong"))
		}
		return
	}

	// compare the password
	isPasswordCorrect := auth.ComparePassword(exisitingUser.Password, payload.Password)

	if !isPasswordCorrect {
		utils.WriteError(w, http.StatusBadRequest, errors.New("Invalid email or password"))
		return
	}

	// create token and send to them
	token, err := auth.CreateJWT(exisitingUser.Id)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("Something went wrong while generating token"))
		return
	}

	auth.SetTokenCookie(w, token)
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "logged in successfully",
	})
}

func (c *Handler) HandleLogout(w http.ResponseWriter, req *http.Request) {
	auth.ClearTokenCookie(w)
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
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
