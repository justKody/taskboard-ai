package user

import (
	"net/http"

	"github.com/justKody/taskboard-go-api/utils"
)

func (c *Handler) HandleLogin(w http.ResponseWriter, req *http.Request) {

}

func (c *Handler) HandleSignup(w http.ResponseWriter, req *http.Request) {
	var payload SignupRequest
	if err := utils.ParseJSON(req, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	
}
