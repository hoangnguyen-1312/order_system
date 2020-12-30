package handler

import (
	"order_system/user_interface"
	"order_system/entity"
	"order_system/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Authenticate struct {
	us user_interface.UserAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewAuthenticate(userInterface user_interface.UserAppInterface, authJWT auth.AuthInterface, token auth.TokenInterface) *Authenticate {
	return &Authenticate{
		us: userInterface,
		rd: authJWT,
		tk: token,
	}
}

func (au *Authenticate) Login(c *gin.Context) {
	var user *entity.User
	var tokenErr = map[string]string{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateUser)
		return
	}
	u, userErr := au.us.GetUserByEmailAndPassword(user)
	if userErr != nil {
		c.JSON(http.StatusInternalServerError, userErr)
		return
	}
	ts, tErr := au.tk.CreateToken(u.ID)
	if tErr != nil {
		tokenErr["token_error"] = tErr.Error()
		c.JSON(http.StatusUnprocessableEntity, tErr.Error())
		return
	}
	saveErr := au.rd.CreateAuth(u.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr.Error())
		return
	}
	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	c.JSON(http.StatusOK, userData)
}

func (au *Authenticate) Logout(c *gin.Context) {
	metadata, err := au.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	deleteErr := au.rd.DeleteTokens(metadata)
	if deleteErr != nil {
		c.JSON(http.StatusUnauthorized, deleteErr.Error())
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}