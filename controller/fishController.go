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
	"time"

	"github.com/gin-gonic/gin"
)

type Fish model.Fish

func C_CreateFish(c *gin.Context) {
	fishName := c.PostForm("fish_name")
	fishType := c.PostForm("fish_type")
	fishSize := c.PostForm("fish_size")
	fishDescription := c.PostForm("fish_description")
	fishPrice := c.PostForm("fish_price")
	imagePath, err := c.FormFile("image_path")

	if err != nil {
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})

		fmt.Println("ERROR 1")
		return
	}
	ext := filepath.Ext(imagePath.Filename)

	currentTime := time.Now()

	// Format the date in a desired way for the file name
	// For example, "2006-01-02" will give you a file name like "2023-07-24"
	fileNameDate := currentTime.Format("20060102")

	newFilename := fmt.Sprintf("%s%s", "fish"+"_"+fileNameDate+"_"+helper.GenerateRandomString(10), ext)

	// Save the uploaded file to the server
	if err := c.SaveUploadedFile(imagePath, filepath.Join("fishuploads", newFilename)); err != nil {
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		fmt.Println("ERROR 2")
		return
	}

	price, _ := strconv.Atoi(fishPrice)

	fish := Fish{
		FishName:        fishName,
		FishType:        fishType,
		FishSize:        fishSize,
		FishDescription: fishDescription,
		FishPrice:       price,
		ImagePath:       "/fishuploads/" + newFilename,
	}

	tx := config.GetDB().Begin()

	if err := tx.Create(&fish).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		fmt.Println("ERROR 3")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		fmt.Println("ERROR 4")
		return
	}

	Logger, _ := json.Marshal(&fish)
	helper.Logger("info", "Created Kios Success, Response: "+string(Logger))

	c.JSON(utils.STATUS_CREATED, gin.H{
		"status":  "success",
		"message": "Created Fish Success",
		"data":    fish,
	})

}

func C_GetAllFish(c *gin.Context) {

	var data []Fish

	tx := config.GetDB().Begin()

	if err := tx.Preload("Seller").Find(&data).Error; err != nil {
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
		"message": "Fetch Fish Success",
		"data":    data,
	})
}

func C_GetFishById(c *gin.Context) {

	var data Fish

	id := c.Param("id")
	fishId, _ := strconv.Atoi(id)

	tx := config.GetDB().Begin()

	if err := tx.Where("id = ?", fishId).Preload("Seller").First(&data).Error; err != nil {
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
		"message": "Fetch Fish Success",
		"data":    data,
	})
}

func C_UpdateFish(c *gin.Context) {

	id := c.Param("id")
	fishId, _ := strconv.Atoi(id)

	var data Fish

	fishName := c.PostForm("fish_name")
	fishType := c.PostForm("fish_type")
	fishSize := c.PostForm("fish_size")
	fishDescription := c.PostForm("fish_description")
	fishPrice := c.PostForm("fish_price")

	price, _ := strconv.Atoi(fishPrice)

	if fishName != "" {
		data.FishName = fishName
	}

	if fishType != "" {
		data.FishType = fishType
	}

	if fishSize != "" {
		data.FishSize = fishSize
	}

	if fishDescription != "" {
		data.FishDescription = fishDescription
	}
	if fishPrice != "" {
		data.FishPrice = price
	}

	imagePath, err := c.FormFile("image_path")

	if err == nil {
		ext := filepath.Ext(imagePath.Filename)
		currentTime := time.Now()
		fileNameDate := currentTime.Format("20060102")

		newFilename := fmt.Sprintf("%s%s", "fish"+"_"+fileNameDate+"_"+helper.GenerateRandomString(10), ext)

		// Save the uploaded file to the server
		if err := c.SaveUploadedFile(imagePath, filepath.Join("fishuploads", newFilename)); err != nil {
			helper.Logger("fail", "In Server: "+err.Error())
			c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		data.ImagePath = "/fishuploads/" + newFilename
	}

	tx := config.GetDB().Begin()

	if err := tx.Model(Fish{}).Where("id = ?", fishId).Updates(&data).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Where("id = ?", fishId).First(&data).Error; err != nil {
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
	helper.Logger("info", "Update Fish Success, Response: "+string(Logger))

	c.JSON(utils.STATUS_OK, gin.H{
		"status":  "success",
		"message": "Update Fish Success",
		"data":    data,
	})

}

func C_DeleteFish(c *gin.Context) {
	id := c.Param("id")
	fishId, _ := strconv.Atoi(id)

	var data Fish

	tx := config.GetDB().Begin()

	if err := tx.Model(Fish{}).Where("id = ?", fishId).First(&data).Delete(&Fish{}).Error; err != nil {
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
	helper.Logger("info", "Delete Fish Success, Response: "+string(Logger))

	c.JSON(utils.STATUS_OK, gin.H{
		"status":  "success",
		"message": "Delete Fish Success",
		"data":    nil,
	})

}
