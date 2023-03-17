package users

import (
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

type LoginValidator struct {
	User struct {
		Email    string `form:"email" json:"email" binding:"required,email"`
		Password string `form:"password" json:"password" binding:"required,min=8,max=255"`
	} `json:"user"`
	userModel User `json:"-"`
}

func (self *LoginValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}

	self.userModel.Email = self.User.Email
	self.userModel.Password = self.User.Password
	return nil
}

func NewLoginValidator() LoginValidator {
	LoginValidator := LoginValidator{}
	return LoginValidator
}

// Register validator struct
type UserModelValidator struct {
	Email       string `form:"email" json:"email" binding:"required,min=4,max=255"`
	PhoneNumber string `form:"phone_number" json:"phone_number" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required,min=8,max=255"`
	userModel   User   `json:"-"`
}

// Bind register body to validator model
func (self *UserModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.userModel.Email = self.Email
	self.userModel.PhoneNumber = self.PhoneNumber
	self.userModel.setPassword(self.Password)
	return nil
}

func NewUserModelValidator() UserModelValidator {
	userModelValidator := UserModelValidator{}
	return userModelValidator
}
