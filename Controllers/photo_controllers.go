package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"task-5-vix-fullstack/app"
	"task-5-vix-fullstack/app/auth"
	"task-5-vix-fullstack/helpers/formaterror"
	"task-5-vix-fullstack/models"
)

// Menampilkan semua photo
func GetPhoto(c *gin.Context) {
	//Membuat list object photo untuk menyimpan data photo yang akan ditampilkan
	photos := []models.Photo{}

	//Melakukan inisialisasi list object foto dari db ke photos
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Debug().Model(&models.Photo{}).Limit(100).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": "photo not found", "data": nil})
		return
	}

	//Melakukan inisialisasi author setiap photo
	if len(photos) > 0 {
		for i := range photos {
			user := models.User{}
			err := db.Model(&models.User{}).Where("id = ?", photos[i].UserId).Take(&user).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
				return
			}

			photos[i].Author = app.Author{
				ID: user.ID, Username: user.Username, Email: user.Email,
			}
		}
	}

	//Response success
	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "success", "data": photos})
}

// Membuat photo profile
func CreatePhoto(c *gin.Context) {
	//set database
	db := c.MustGet("db").(*gorm.DB)

	// Mengambil token Bearer yang ada
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		return
	}

	// Mengambil email user login berdasarkan token jwt
	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	//Mengambil data user yang login melalui token jwt
	var user_login models.User

	err = db.Debug().Where("email = ?", email).First(&user_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	// Membaca data body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	//Mengubah json menjadi object Photo
	photo_input := models.Photo{}
	err = json.Unmarshal(body, &photo_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Inisialisasi data photo
	photo_input.Initialize()
	photo_input.UserId = user_login.ID
	photo_input.Author = app.Author{
		ID:       user_login.ID,
		Username: user_login.Username,
		Email:    user_login.Email,
	}

	// Validasi data photo
	err = photo_input.Validate("upload")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan pengecekan apakah user sudah mengupload photo, jika sudah ada maka akan membuat photo
	var old_photo models.Photo
	err = db.Debug().Model(&models.Photo{}).Where("user_id = ?", user_login.ID).Find(&old_photo).Error
	if err != nil {
		if err.Error() == "record not found" {
			//Melakukan create photo ke database
			err = db.Debug().Create(&photo_input).Error
			if err != nil {
				formattedError := formaterror.ErrorMessage(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "T", "message": "success upload photo", "data": photo_input})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan update photo ke database jika user sudah memiliki photo
	photo_input.ID = old_photo.ID
	err = db.Debug().Model(&old_photo).Updates(&photo_input).Error
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	//Response succes
	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "success change photo", "data": photo_input})

}

// Melakukan update photo profile
func UpdatePhoto(c *gin.Context) {
	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Mengambil token Bearer yang ada
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		return
	}

	//Mengambil email user login berdasarkan token jwt
	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	//Mengambil data user yang login melalui token jwt
	var user_login models.User

	err = db.Debug().Where("email = ?", email).First(&user_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	// Membaca data body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	//Mengubah json menjadi object Photo
	photo_input := models.Photo{}
	err = json.Unmarshal(body, &photo_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan validasi data photo
	err = photo_input.Validate("change")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Pengecekan foto berdasarkan id
	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "Photo not found", "data": nil})
		return
	}

	//Validasi user tidak dapat mengupdate photo yang dibuat user lain
	if user_login.ID != photo.UserId {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "no access to change photo", "data": nil})
		return
	}

	//Melakukan update photo ke database
	err = db.Model(&photo).Updates(&photo_input).Error
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	//Custom response data
	photo.Author = app.Author{
		ID:       user_login.ID,
		Username: user_login.Username,
		Email:    user_login.Email,
	}

	//Response success
	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "success change photo", "data": photo})
}

// Menghapus photo profile
func DeletePhoto(c *gin.Context) {

	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Mengambil token Bearer yang ada
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		return
	}

	//Mengambil email user login berdasarkan token jwt
	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	//Mengambil data user yang login melalui token jwt
	var user_login models.User
	if err := db.Debug().Where("email = ?", email).First(&user_login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	//Pengecekan foto berdasarkan id
	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "Photo not found", "data": nil})
		return
	}

	//Validasi user tidak dapat mengupdate photo yang dibuat user lain
	if user_login.ID != photo.UserId {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "no access to delete photo", "data": nil})
		return
	}

	//Menghapus photo dari database
	err = db.Debug().Delete(&photo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Response success
	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "delete photo success", "data": nil})
}




