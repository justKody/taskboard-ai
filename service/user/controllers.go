package user

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/justKody/taskboard-go-api/db/sqlc"
	"github.com/justKody/taskboard-go-api/middleware"
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

func (c *Handler) HandleGetMe(w http.ResponseWriter, req *http.Request) {

	userID, ok := middleware.GetUserID(req.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	user, err := c.store.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("user not found"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (c *Handler) HandleUpdateUser(w http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserID(req.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	user, err := c.store.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	if user == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("user not found"))
		return
	}

	var payload UpdateUserRequestDTO
	if err := utils.ParseJSON(req, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// update the user in database
	params := sqlc.UpdateUserParams{
		ID:   userID,
		Name: payload.Name,
	}

	updatedUser, err := c.store.UpdateUser(params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedUser)
}

func (c *Handler) HandleChangePassword(w http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserID(req.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	user, err := c.store.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("user not found"))
		return
	}

	var payload ChangePasswordRequestDTO
	if err := utils.ParseJSON(req, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !auth.ComparePassword(user.Password, payload.OldPassword) {
		utils.WriteError(w, http.StatusBadRequest, errors.New("incorrect old password"))
		return
	}

	hashedNewPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	params := sqlc.ChangePasswordParams{
		ID:       userID,
		Password: hashedNewPassword,
	}

	updatedUser, err := c.store.ChangePassword(params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedUser)
}

func (c *Handler) HandleDeleteUser(w http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserID(req.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	err := c.store.DeleteUser(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user deleted successfully",
	})
}

func (c *Handler) HandleGetUsersList(w http.ResponseWriter, req *http.Request) {
	users, err := c.store.GetUsersList()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}
