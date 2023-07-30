package controller

import (
	"backend-golang/config"
	"backend-golang/helper"
	"backend-golang/model"
	"backend-golang/utils"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Register model.Register
type Login model.Login
type User model.User
type Response map[string]interface{}

func VerifyPassword(hashedPass, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func RegisterController(c *gin.Context) {
	register := Register{}

	bind := c.BindJSON(&register)

	if bind != nil {
		helper.Logger("fail", "In Server: "+bind.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": bind.Error(),
		})
		return
	}

	rune := []rune(register.PhoneNumber)
	if string(rune[0]) == "0" {
		register.PhoneNumber = strings.Replace(register.PhoneNumber, "0", "62", 1)
	}

	hashedPass, err := Hash(register.Password)
	if err != nil {
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": bind.Error(),
		})
		return
	}

	user := User{}
	user.Fullname = register.Fullname
	user.Email = register.Email
	user.Password = string(hashedPass)
	user.PhoneNumber = register.PhoneNumber
	user.Role = register.Role

	tx := config.GetDB().Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": bind.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": bind.Error(),
		})
		return
	}

	response := Response{}
	response["fullname"] = user.Fullname
	response["email"] = user.Email
	response["phone_number"] = user.PhoneNumber
	response["role"] = user.Role

	Logger, _ := json.Marshal(&response)
	helper.Logger("info", "Register Success, Response: "+string(Logger))

	c.JSON(utils.STATUS_CREATED, gin.H{
		"status":  "success",
		"message": "Register Success",
		"data":    response,
	})

}

func LoginController(c *gin.Context) {
	login := Login{}

	auth := c.GetHeader("Authorization")

	// Check if the request body is empty
	if c.Request.ContentLength == 0 {
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": "Empty request body",
		})
		c.Abort()
		return
	}

	if err := c.BindJSON(&login); err != nil {

		helper.Logger("fail", "In Server: "+err.Error())
		c.AbortWithStatusJSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if auth == "" {
		helper.Logger("fail", "In Server: "+"Unauthorized")

		c.AbortWithStatusJSON(utils.STATUS_UNAUTHORIZE, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
		return
	}

	hash := sha256.Sum256([]byte(login.Email + login.Password))

	hexEncrypt := hex.EncodeToString(hash[:])

	encrypt := base64.StdEncoding.EncodeToString([]byte(hexEncrypt))

	authFromEncrypt := "Basic " + encrypt

	if auth != authFromEncrypt {
		helper.Logger("fail", "In Server: "+"Unauthorized")
		c.AbortWithStatusJSON(utils.STATUS_UNAUTHORIZE, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
		return
	}

	user := User{}

	tx := config.GetDB().Begin()

	if err := tx.Where("email = ?", login.Email).Take(&user).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	isValid := VerifyPassword(user.Password, login.Password)

	if isValid != nil && isValid == bcrypt.ErrMismatchedHashAndPassword {
		helper.Logger("error", "In Server: "+isValid.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": "Password is incorrect",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	response := Response{}
	response["fullname"] = user.Fullname
	response["email"] = user.Email
	response["phone_number"] = user.PhoneNumber
	response["role"] = user.Role
	response["token"] = encrypt

	Logger, _ := json.Marshal(&response)
	helper.Logger("info", "Login Success, Response: "+string(Logger))

	c.JSON(utils.STATUS_OK, gin.H{
		"status":  "success",
		"message": "Login Success",
		"data":    response,
	})

}
