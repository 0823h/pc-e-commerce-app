package users

import (
	"github.com/gin-gonic/gin"
)

type ProfileSerializer struct {
	C *gin.Context
	User
}

// Declare your response schema here
type ProfileResponse struct {
	ID             uint64  `json:"-"`
	Email          string  `json:"email"`
	ProfilePicture *string `json:"profile_picture"`
}

func (self *ProfileSerializer) Response() ProfileResponse {
	// myUserModel := self.C.MustGet("my_user_model").(UserModel)
	profile := ProfileResponse{
		ID:             self.ID,
		Email:          self.Email,
		ProfilePicture: &self.ProfilePicture,
	}
	return profile
}

type UserResponse struct {
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
}

type UserSerializer struct {
	c *gin.Context
	User
}

func (self *UserSerializer) Response() UserResponse {
	// myUserModel := self.c.MustGet("my_user_model").(User)
	user := UserResponse{
		Email:          self.Email,
		ProfilePicture: self.ProfilePicture,
	}
	return user
}
