package controllers

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"

	"../../api/auth"
	message "../../api/constants"
	"../../api/models"
	"../../api/responses"
	"../../api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {
	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	token, _ := auth.CreateToken(user.ID)
	return token, nil
}

func (server *Server) SignUp(w http.ResponseWriter, r *http.Request) {

	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = u.Validate("sign-up"); err != nil {
		http.Error(w, message.InvalidParams, http.StatusBadRequest)
		return
	}

	err = server.
		DB.
		Debug().
		Model(models.User{}).Where("email = ?", u.Email).Find(&u).Error

	if !gorm.IsRecordNotFoundError(err) {
		http.Error(w, message.UserExist, http.StatusConflict)
		return
	}

	u.Prepare()

	if err = server.DB.Debug().Model(&models.User{}).Create(&u).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responses.JSON(w, http.StatusOK, u)
	return

}
