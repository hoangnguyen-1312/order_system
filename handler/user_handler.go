package handler

import (
	_"fmt"
	"net/http"
	"order_system/auth"
	"order_system/entity"
	"order_system/user_interface"
	"strconv"

	"github.com/gin-gonic/gin"
)

//Users struct defines the dependencies that will be used
type Users struct {
	us user_interface.UserAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

//Users constructor
func NewUsers(us user_interface.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Users {
	return &Users{
		us: us,
		rd: rd,
		tk: tk,
	}
}

func (s *Users) SaveUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}
	//validate the request:
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateErr)
		return
	}
	newUser, err := s.us.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Set("user_id", newUser.ID)
	c.JSON(http.StatusCreated, newUser.PublicUser())
}

func (s *Users) GetUsers(c *gin.Context) {
	var err error
	metadata, err := s.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//lookup the metadata in redis:
	_, err = s.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	users := entity.Users{} //customize user
	
	//us, err = application.UserApp.GetUsers()
	users, err = s.us.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}

func (s *Users) GetUser(c *gin.Context) {
	metadata, err := s.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := s.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userId, err = strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := s.us.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user.PublicUser())
}

// func (s *Users) GetCurrentUserInfor(c *gin.Context) {
// 	userId := c.GetInt64("user_id")
// 	fmt.Println(userId)
// 	// user, err := s.us.GetUser(userId.(uint64))
// 	// if err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, err.Error())
// 	// 	return
// 	// }
// 	// c.JSON(http.StatusOK, user.PublicUser())
// }




