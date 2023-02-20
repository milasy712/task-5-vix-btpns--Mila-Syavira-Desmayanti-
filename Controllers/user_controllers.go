package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"task-5-vix-fullstack/app"
	"task-5-vix-fullstack/app/auth"
	"task-5-vix-fullstack/helpers/formaterror"
	"task-5-vix-fullstack/helpers/hash"
	"task-5-vix-fullstack/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Register User
func CreateUser(c *gin.Context) {
	//set database
	db := c.MustGet("db").(*gorm.DB)

	// Baca data body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	// ubah json menjadi object User
	user_input := models.User{}
	err = json.Unmarshal(body, &user_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// inisialisasi data user
	user_input.Initialize()

	//Melakukan validasi data user
	err = user_input.Validate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan Hash password
	err = user_input.HashPassword()
	if err != nil {
		log.Fatal(err)
	}

	//Create data user ke database
	err = db.Debug().Create(&user_input).Error
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	//Custom response data
	data := app.UserRegister{
		ID:        user_input.ID,
		Username:  user_input.Username,
		Email:     user_input.Email,
		CreatedAt: user_input.CreatedAt,
		UpdatedAt: user_input.UpdatedAt,
	}

	//Response success
	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "register user success", "data": data})
}

// Melakukan login
func Login(c *gin.Context) {
	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Membaca data dari body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Mengubah json ke objek user
	user_input := models.User{}
	err = json.Unmarshal(body, &user_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan perisapan inisialisi dan validasi
	user_input.Initialize()
	err = user_input.Validate("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan pengecekan user di database berdasarkan email
	var user_login app.UserLogin

	err = db.Debug().Table("users").Select("*").Joins("left join photos on photos.user_id = users.id").
		Where("users.email = ?", user_input.Email).Find(&user_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	//Melakukan verifikasi password db dan user input
	err = hash.VerifyPassword(user_login.Password, user_input.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	//Ketika berhasil login akan membuat token jwt
	token, err := auth.GenerateJWT(user_login.Email, user_login.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan custom response untuk data
	data := app.DataUser{
		ID: user_login.ID, Username: user_login.Username, Email: user_login.Email, Token: token,
		Photos: app.Photo{Title: user_login.Title, Caption: user_login.Caption, PhotoUrl: user_login.PhotoUrl},
	}

	//Ketika berhasil login memeberikan response success
	c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "T", "message": "login success", "data": data})
}

// Update data user
func UpdateUser(c *gin.Context) {

	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Melakukan cek apakah data yang ingin diubah ada berdasarkan user id dari param
	var user models.User
	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	//Membaca data body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	//Mengubah json menjadi object User
	user_input := models.User{}
	user_input.ID = user.ID
	err = json.Unmarshal(body, &user_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan validasi data user
	err = user_input.Validate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan Hash password
	err = user_input.HashPassword()
	if err != nil {
		log.Fatal(err)
	}

	//Melakukan update data user ke database
	err = db.Debug().Model(&user).Updates(&user_input).Error
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	//Custom response data
	data := app.UserRegister{
		ID:        user_input.ID,
		Username:  user_input.Username,
		Email:     user_input.Email,
		CreatedAt: user_input.CreatedAt,
		UpdatedAt: user_input.UpdatedAt,
	}

	//Response success
	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "update user success", "data": data})
}

// Menghapus data user
func DeleteUser(c *gin.Context) {

	//Set databse
	db := c.MustGet("db").(*gorm.DB)

	//Melakukan cek apakah data yang ingin dihapus ada berdasarkan user id dari param
	var user models.User

	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	//Menghapus data user dari database
	err = db.Debug().Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Response succes
	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "delete user success", "data": nil})
}











