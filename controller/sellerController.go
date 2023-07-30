package controller

import (
	"backend-golang/config"
	"backend-golang/helper"
	"backend-golang/model"
	"backend-golang/utils"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Seller model.Seller

func C_CreateSeller(c *gin.Context) {
	sellerName := c.PostForm("seller_name")
	sellerAddress := c.PostForm("seller_address")
	sellerPhone := c.PostForm("seller_phone")
	sellerDescription := c.PostForm("seller_description")
	imagePath, err := c.FormFile("image_path")

	if err != nil {
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	ext := filepath.Ext(imagePath.Filename)

	currentTime := time.Now()

	// Format the date in a desired way for the file name
	// For example, "2006-01-02" will give you a file name like "2023-07-24"
	fileNameDate := currentTime.Format("20060102")

	newFilename := fmt.Sprintf("%s%s", "seller"+"_"+fileNameDate+"_"+helper.GenerateRandomString(10), ext)

	// Save the uploaded file to the server
	if err := c.SaveUploadedFile(imagePath, filepath.Join("uploads", newFilename)); err != nil {
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	seller := Seller{
		SellerName:        sellerName,
		SellerAddress:     sellerAddress,
		SellerPhoneNumber: sellerPhone,
		SelerDescription:  sellerDescription,
		ImagePath:         "/uploads/" + newFilename,
	}

	rune := []rune(seller.SellerPhoneNumber)
	if len(rune) > 0 { // Check if the string has at least one character
		if string(rune[0]) == "0" {
			seller.SellerPhoneNumber = strings.Replace(seller.SellerPhoneNumber, "0", "62", 1)
		}
	}

	tx := config.GetDB().Begin()

	if err := tx.Create(&seller).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	Logger, _ := json.Marshal(&seller)
	helper.Logger("info", "Created Seller Success, Response: "+string(Logger))

	c.JSON(utils.STATUS_CREATED, gin.H{
		"status":  "success",
		"message": "Created Seller Success",
		"data":    seller,
	})

}

func C_GetAllSeller(c *gin.Context) {

	var data []Seller

	tx := config.GetDB().Begin()

	if err := tx.Preload("Fish").Find(&data).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(utils.STATUS_OK, gin.H{
		"status":  "success",
		"message": "Fetch Seller Success",
		"data":    data,
	})
}

func C_GetSellerById(c *gin.Context) {

	var data Seller

	id := c.Param("id")
	sellerId, _ := strconv.Atoi(id)

	tx := config.GetDB().Begin()

	if err := tx.Where("id = ?", sellerId).Preload("Fish").First(&data).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(utils.STATUS_OK, gin.H{
		"status":  "success",
		"message": "Fetch Seller Success",
		"data":    data,
	})
}

func C_UpdateSeller(c *gin.Context) {

	id := c.Param("id")
	sellerId, _ := strconv.Atoi(id)

	var data Seller

	sellerName := c.PostForm("seller_name")
	sellerAddress := c.PostForm("seller_address")
	sellerPhone := c.PostForm("seller_phone")
	sellerDesc := c.PostForm("seller_description")

	if sellerName != "" {
		data.SellerName = sellerName
	}

	if sellerDesc != "" {
		data.SelerDescription = sellerDesc
	}

	if sellerAddress != "" {
		data.SellerAddress = sellerAddress
	}

	if sellerPhone != "" {
		rune := []rune(data.SellerPhoneNumber)
		if len(rune) > 0 { // Check if the string has at least one character
			if string(rune[0]) == "0" {
				data.SellerPhoneNumber = strings.Replace(data.SellerPhoneNumber, "0", "62", 1)
				return
			}
		}
		data.SellerPhoneNumber = sellerPhone
	}

	imagePath, err := c.FormFile("image_path")

	if err == nil {
		ext := filepath.Ext(imagePath.Filename)
		currentTime := time.Now()
		fileNameDate := currentTime.Format("20060102")

		newFilename := fmt.Sprintf("%s%s", "seller"+"_"+fileNameDate+"_"+helper.GenerateRandomString(10), ext)

		// Save the uploaded file to the server
		if err := c.SaveUploadedFile(imagePath, filepath.Join("uploads", newFilename)); err != nil {
			helper.Logger("fail", "In Server: "+err.Error())
			c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		data.ImagePath = "/uploads/" + newFilename
	}

	tx := config.GetDB().Begin()

	if err := tx.Model(Seller{}).Where("id = ?", sellerId).Updates(&data).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Where("id = ?", sellerId).First(&data).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	Logger, _ := json.Marshal(&data)
	helper.Logger("info", "Update Seller Success, Response: "+string(Logger))

	c.JSON(utils.STATUS_OK, gin.H{
		"status":  "success",
		"message": "Update Seller Success",
		"data":    data,
	})

}

func C_DeleteSeller(c *gin.Context) {
	id := c.Param("id")
	sellerId, _ := strconv.Atoi(id)

	var data Seller

	tx := config.GetDB().Begin()

	if err := tx.Model(Seller{}).Where("id = ?", sellerId).First(&data).Delete(&Seller{}).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	Logger, _ := json.Marshal(&data)
	helper.Logger("info", "Delete Seller Success, Response: "+string(Logger))

	c.JSON(utils.STATUS_OK, gin.H{
		"status":  "success",
		"message": "Delete Seller Success",
		"data":    nil,
	})

}
