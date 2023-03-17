package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
	"github.com/zhenghaoz/gorse/client"
)

func UsersLogin(c *gin.Context) {
	loginValidator := NewLoginValidator()
	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	userModel, err := FindOneUser(&User{Email: loginValidator.userModel.Email})

	if err != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}

	if userModel.checkPassword(loginValidator.User.Password) != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}

	// UpdateContextUserModel(c, userModel.ID)
	// serializer := UserSerializer{c}
	// c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
	jwtString, err := common.GenerateToken(userModel.ID, loginValidator.userModel.Email)

	if err != nil {
		common.SendResponse(c, http.StatusConflict, err.Error(), nil)
		return
	}

	common.SendResponse(c, http.StatusOK, "Success", gin.H{"token": jwtString})
	return
}

func UsersRegistration(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	fmt.Println("Email " + userModelValidator.userModel.Email)
	// Check if email has been used
	if userModelValidator.userModel.checkEmailExisted() {
		c.JSON(http.StatusBadRequest, common.NewError("register", errors.New("Email has been used")))
		return
	}

	// Check if phone number has been used
	if userModelValidator.userModel.checkPhoneNumberExisted() {
		c.JSON(http.StatusBadRequest, common.NewError("register", errors.New("Phone number has been used")))
		return
	}

	db := common.GetDB()
	err := db.Create(&userModelValidator.userModel).Error
	if err != nil {
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	// TO DO : add joins address for user
	createdUser := NewUser()
	db.Where("users.id = ?", userModelValidator.userModel.ID).First(&createdUser)

	gorse := common.GetGorse()
	gorse.InsertUser(context.Background(), client.User{UserId: strconv.FormatUint(createdUser.ID, 10), Labels: []string{"User labels"}})

	serializer := UserSerializer{c, createdUser}
	c.JSON(http.StatusCreated, gin.H{"user": serializer.Response()})
}
