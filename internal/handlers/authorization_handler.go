package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/manabie-com/togo/internal/authorizers"
	"github.com/manabie-com/togo/internal/services"
	"github.com/manabie-com/togo/internal/utils"
)

type AuthorizationHandler struct {
	Authorizer  authorizers.Authorizer
	UserService services.UserService
}

func (h *AuthorizationHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)
	userID := req.FormValue("user_id")
	pwd := req.FormValue("password")
	token, err := h.UserService.VerifyUser(req.Context(), userID, pwd)
	if err != nil {
		if services.IncorrectUserIDPassword == err {
			utils.RespondWithError(resp, http.StatusUnauthorized, err.Error())
			return
		}
		utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
		return
	}
	data := struct {
		Token string `json:"token"`
	}{
		token,
	}
	utils.RespondWithJSON(resp, http.StatusOK, data)
	return
}

// HandleAuthorization validate the auth token from the request header
func (h *AuthorizationHandler) HandleAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		log.Println("Executing authorization")
		if req.URL.Path != "/login" {
			// Grab the raw Authorization header
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Authorization header is missing")
				utils.RespondWithError(resp, http.StatusUnauthorized, "Invalid authorization token")
				return
			}

			// Confirm the request is sending Bearer token.
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Println("Authorization requires Bearer scheme")
				utils.RespondWithError(resp, http.StatusUnauthorized, "Invalid authorization token")
				return
			}

			// Get the token from the request header
			// The first seven characters are skipped - "Bearer ".
			token := authHeader[7:]
			id, ok := h.Authorizer.VerifyToken(token)
			if !ok {
				log.Println("Invalid authorization token")
				utils.RespondWithError(resp, http.StatusUnauthorized, "Invalid authorization token")
				return
			}
			req = req.WithContext(context.WithValue(req.Context(), utils.ContextUserKey, id))
		}
		next.ServeHTTP(resp, req)
	})
}
