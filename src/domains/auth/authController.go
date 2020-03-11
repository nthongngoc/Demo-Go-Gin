package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khoa5773/go-server/src/shared"
)

func ApplyRoutes(r *gin.Engine) {
	authController := r.Group("/auth")
	authController.GET("/login", login)
	authController.GET("/login/callback", loginCallback)
	authController.GET("/signup", signup)
	authController.GET("/signup/callback", signupCallback)
}

func login(c *gin.Context) {
	authenticator, err := shared.NewAuthenticator("login")
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authenticator.Config.AuthCodeURL(""))
}

func signup(c *gin.Context) {
	authenticator, err := shared.NewAuthenticator("signup")
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authenticator.Config.AuthCodeURL(""))
}

func loginCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := handleLoginCallback(c, code)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": token})
}

func signupCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := handleSignupCallback(c, code)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": token})
}
