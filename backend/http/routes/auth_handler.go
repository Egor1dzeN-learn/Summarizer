package routes

import (
	"net/http"
	"summarizer/backend/http/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{svc: authService}
}

func (h *AuthHandler) Bind(r *gin.RouterGroup) {
	r.GET("/me", h.Me)
	r.POST("/login", h.Login)
	r.POST("/logout", h.Logout)
}

func (h *AuthHandler) Me(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, h.svc.GetUser(uid.(uint)))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		LoginToken string `json:"token"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user := h.svc.FindUserByLoginToken(body.LoginToken)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired login token",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("uid", user.ID)
	session.Save()

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {

}
