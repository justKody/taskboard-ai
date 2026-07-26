package user

import (
	"fmt"
	"net/http"
)

func (c *Handler) HandleLogin(w http.ResponseWriter, req *http.Request) {

}

func (c *Handler) HandleSignup(w http.ResponseWriter, req *http.Request) {
	fmt.Println("HandleSignup")
}
