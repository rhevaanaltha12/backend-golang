package controller

import (
	"backend-golang/config"
	"backend-golang/helper"
	"backend-golang/model"
	"backend-golang/utils"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Temporary model.Temporary

func C_CreateTemporary(c *gin.Context) {

	var sellerData Seller
	var fishData Fish

	idOne := c.Param("seller_id")
	sellerId, _ := strconv.Atoi(idOne)

	idTwo := c.Param("fish_id")
	fishId, _ := strconv.Atoi(idTwo)

	var tempData Temporary

	tx := config.GetDB().Begin()

	if err := tx.Where("id = ?", sellerId).First(&sellerData).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Where("id = ?", fishId).First(&fishData).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Where("seller_id = ? AND fish_id = ?", sellerId, fishId).Find(&tempData).Error; err != nil {
		tx.Rollback()
		helper.Logger("fail", "In Server: "+err.Error())
		c.JSON(utils.STATUS_BAD_REQUEST, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	tempData.SellerId = uint(sellerId)
	tempData.FishId = uint(fishId)

	if err := tx.Table("temporary").Create(&tempData).Error; err != nil {
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

	Logger, _ := json.Marshal(&tempData)
	helper.Logger("info", "Created Kios Success, Response: "+string(Logger))

	c.JSON(utils.STATUS_CREATED, gin.H{
		"status":  "success",
		"message": "Created Temp Success",
		"data":    tempData,
	})

}
