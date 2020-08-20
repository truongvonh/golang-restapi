package controllers

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"time"

	"../../api/auth"
	message "../../api/constants"
	"../../api/models"
	"../../api/responses"
	"../../api/services"
	"../../api/utils"
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
		Model(models.User{}).
		Where("email = ?", u.Email).
		Find(&u).
		Error

	if !gorm.IsRecordNotFoundError(err) {
		http.Error(w, message.UserExist, http.StatusConflict)
		return
	}

	u.Prepare()
	u.VerifyCode = models.GenerateEmailToken()

	if err = server.DB.Debug().Model(&models.User{}).Create(&u).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mailData := services.MailData{
		UserName: u.Email,
		UserMail: u.Email,
		Content: "This is your activation code <strong>" + u.VerifyCode +
			"</strong>.It will be expired after 2 hours",
	}

	if _, err = services.SendMail(mailData); err != nil {
		http.Error(w, message.ValidationEmailFailed, http.StatusBadRequest)
		return
	}

	responses.JSON(w, http.StatusOK, u)
	return

}

func (server *Server) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	body := utils.GetBodyFromRequest(w, r)

	if err := json.Unmarshal(body, &user); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err := server.
		DB.
		Debug().
		Model(models.User{}).
		Where("verifyCode = ?", user.VerifyCode).
		Error

	if gorm.IsRecordNotFoundError(err) {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	now := time.Now()
	updateTime := user.UpdatedAt

	if updateTime.Sub(now).Hours() > 1 {
		http.Error(w, message.ExpiredVerifyCode, http.StatusBadRequest)
		return
	}

	user.VerifyCode = ""
	user.IsActive = true
	err = server.DB.Debug().Model(models.User{}).Save(&user).Error
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, message.VerifyAccountSuccess)
	return
}
