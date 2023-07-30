package main

import (
	"backend-golang/config"
	"backend-golang/controller"
	"backend-golang/helper"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		helper.Logger("error", "Error Getting Env")
	}

	config.GetDB()

	middleUrl := os.Getenv("MIDDLE_URL")

	router := gin.Default()

	// cors
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "OPTIONS", "GET", "DELETE"},
		AllowHeaders: []string{"Accept", "content-type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	}))

	// AUTH
	router.POST(middleUrl+"/auth/register", controller.RegisterController)

	router.POST(middleUrl+"/auth/login", controller.LoginController)

	// KIOS
	router.POST(middleUrl+"/seller/create", controller.C_CreateSeller)
	router.GET(middleUrl+"/seller/get", controller.C_GetAllSeller)
	router.GET(middleUrl+"/seller/get/:id", controller.C_GetSellerById)
	router.POST(middleUrl+"/seller/update/:id", controller.C_UpdateSeller)
	router.DELETE(middleUrl+"/seller/delete/:id", controller.C_DeleteSeller)

	// FISH
	router.POST(middleUrl+"/fish/create", controller.C_CreateFish)
	router.GET(middleUrl+"/fish/get", controller.C_GetAllFish)
	router.GET(middleUrl+"/fish/get/:id", controller.C_GetFishById)
	router.POST(middleUrl+"/fish/update/:id", controller.C_UpdateFish)
	router.DELETE(middleUrl+"/fish/delete/:id", controller.C_DeleteFish)

	// TEMP
	router.GET(middleUrl+"/temp/get/:seller_id/:fish_id", controller.C_CreateTemporary)

	// RENDER FILE UPLOAD
	router.Static("/uploads", "./uploads")
	router.Static("/fishuploads", "./fishuploads")

	port := os.Getenv("PORT")

	router.Run(":" + port)

}
