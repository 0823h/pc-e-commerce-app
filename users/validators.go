package users

import (
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

type LoginValidator struct {
	User struct {
		Email    string `form:"email" json:"email" binding:"exists,email"`
		Password string `form:"password" json:"password" binding:"exists,min=8,max=255"`
	} `json:"user"`
	userModel UserModel `json:"-"`
}

func (self *LoginValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}

	self.userModel.Email = self.User.Email
	return nil
}

func NewLoginValidator() LoginValidator {
	LoginValidator := LoginValidator{}
	return LoginValidator
}

type UserModelValidator struct {
	User struct {
		Email       string `form:"email" json:"email" binding:"exists,alphanum,min=4,max=255"`
		PhoneNumber string `form:"phone_number" json:"phone_number" binding:"exists"`
		Password    string `form:"password" json:"password" binding:"exists,min=8,max=255"`
	} `json:"user"`
	userModel UserModel `json:"-"`
}

func (self *UserModelValidator) Bind (c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.userModel.
}

func NewUserModelValidator() UserModelValidator {
	userModelValidator := UserModelValidator{}
	return userModelValidator
}
