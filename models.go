package main

import (
	// "fmt"
	// "time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func GetUsers() ([]User, error) {
	var users []User

	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func GetUserByID(id int) (User, error) {
	var user User

	result := db.First(&user, id)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func CreateUser(user User) (User, error) {
	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user.Password = string(hashedPassword)

	result := db.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func UpdateUser(id int, user *User) (User, error) {
	var u User
	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user.Password = string(hashedPassword)

	result := db.First(&u, id)
	if result.Error != nil {
		return User{}, result.Error
	}

	result = db.Model(&u).Updates(User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if result.Error != nil {
		return User{}, result.Error
	}

	return *user, err
}

func DeleteUser(id int) error {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return result.Error
	}

	result = db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
