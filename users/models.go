package users

import (
	"errors"
	"fmt"
	"time"
	"tmdt-backend/common"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                  uint64    `gorm:"primaryKey"`
	Email               string    `gorm:"column:email;not null"`
	EmailVerified       string    `gorm:"column:email_verified;default:false"`
	PhoneNumber         string    `gorm:"column:phone_number"`
	PhoneNumberVerified string    `gorm:"column:phone_number_verified;default:false"`
	Password            string    `gorm:"column:password"`
	Otp                 string    `gorm:"column:otp"`
	FirstName           string    `gorm:"column:first_name;not null"`
	LastName            string    `gorm:"column:last_name;not null"`
	Dob                 time.Time `gorm:"column:dob"`
	Gender              string    `gorm:"column:gender"`
	ProfilePicture      string    `gorm:"column:profile_picture"`
	LastLoggedIn        time.Time `gorm:"column:last_logged_in"`
	IsDeleted           bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt           time.Time ``
	UpdatedAt           time.Time ``
}

func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(&User{})
}

func (u *User) setPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.Password = string(passwordHash)
	return nil
}

func (u *User) checkPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func FindOneUser(condition interface{}) (User, error) {
	db := common.GetDB()
	var model User
	err := db.Where(condition).First(&model).Error
	fmt.Println(model, err)
	return model, err
}

func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}

func (u *User) checkEmailExisted() bool {
	db := common.GetDB()
	var user User
	result := db.Where("email = ?", u.Email).First(&user)
	fmt.Println(user.Email)
	if result.RowsAffected > 0 {
		// c.Set("user_id", user.ID)
		return true
	}
	return false
}

func (u *User) checkPhoneNumberExisted() bool {
	db := common.GetDB()
	var user User
	result := db.Where("phone_number = ?", u.PhoneNumber).First(&user)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

func NewUser() User {
	var user User
	return user
}
