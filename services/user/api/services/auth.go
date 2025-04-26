package services

import (
	"errors"

	"github.com/Creative-genius001/Stacklo/services/user/db"
	"github.com/Creative-genius001/Stacklo/services/user/models"
	"github.com/Creative-genius001/Stacklo/services/user/types"
	"github.com/Creative-genius001/Stacklo/services/user/utils"
	"github.com/google/uuid"
)

func LoginService() {

}

func RegisterService(RegForm types.RegisterType) error {
	user := models.User{
		Id:           uuid.New().String(),
		Phone:        RegForm.Phone,
		PasswordHash: RegForm.Password,
		Email:        RegForm.Email,
		FirstName:    RegForm.FirstName,
		LastName:     RegForm.LastName,
	}

	//check if user already exist in db
	var count int64
	db.DB.Model(&user).Where("email = ?", RegForm.Email).Count(&count)
	if count > 0 {
		return errors.New("email already exists")
	} else {
		//hash password
		passwordHash, err := utils.HashPassword(RegForm.Password)
		if err != nil {
			return errors.New("hashing error")
		}

		user.PasswordHash = passwordHash

		//add user to database
		res := db.DB.Create(&user)
		if res.Error != nil {
			return errors.New("signup failed! server error")
		}

	}
	return nil
}
