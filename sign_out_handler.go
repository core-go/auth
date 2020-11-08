package auth

import (
	"net/http"
	"strings"
	"time"
)

type SignOutHandler struct {
	TokenVerifier         TokenVerifier
	Resource              string
	Action                string
	Secret                string
	TokenBlacklistService TokenBlacklistService
	LogWriter             AuthActivityLogWriter
}
func NewDefaultSignOutHandler(tokenVerifier TokenVerifier, secret string, tokenBlacklistService TokenBlacklistService, logWriter AuthActivityLogWriter) *SignOutHandler {
	return NewSignOutHandler(tokenVerifier, "", "", secret, tokenBlacklistService, logWriter)
}
func NewSignOutHandler(tokenVerifier TokenVerifier, resource string, action string, secret string, tokenBlacklistService TokenBlacklistService, logWriter AuthActivityLogWriter) *SignOutHandler {
	if len(resource) == 0 {
		resource = "authentication"
	}
	if len(action) == 0 {
		action = "signout"
	}
	return &SignOutHandler{TokenVerifier: tokenVerifier, Resource: resource, Action: action, Secret: secret, TokenBlacklistService: tokenBlacklistService, LogWriter: logWriter}
}
func (h *SignOutHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	data := r.Header["Authorization"]
	token := parseToken(data)

	if len(token) == 0 {
		RespondString(w, r, http.StatusBadRequest, "Invalid token")
		return
	}

	_, _, expiresAt, er1 := h.TokenVerifier.VerifyToken(token, h.Secret)

	if er1 != nil {
		RespondString(w, r, http.StatusBadRequest, "Invalid token")
		return
	}

	if h.TokenBlacklistService == nil {
		RespondString(w, r, http.StatusInternalServerError, "No service to sign out")
		return
	}

	expiresTime := time.Unix(expiresAt, 0)

	er2 := h.TokenBlacklistService.Revoke(token, "The token has signed out.", expiresTime)
	if er2 != nil {
		if h.LogWriter != nil {
			h.LogWriter.Write(r.Context(), h.Resource, h.Action, false, er2.Error())
		}
		RespondString(w, r, http.StatusInternalServerError, er2.Error())
		return
	}
	if h.LogWriter != nil {
		h.LogWriter.Write(r.Context(), h.Resource, h.Action, true, "")
	}
	RespondString(w, r, http.StatusOK, "true")
}

func parseToken(data []string) string {
	if len(data) == 0 {
		return ""
	}
	authorization := data[0]
	if strings.HasPrefix(authorization, "Bearer ") != true {
		return ""
	}
	return authorization[7:]
}