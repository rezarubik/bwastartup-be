package main

import (
	"bwastartup/handler"
	"bwastartup/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	userService.SaveAvatar(3, "images/1-profile.png")

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)                    //note: Register
	api.POST("/sessions", userHandler.Login)                        //note: Login
	api.POST("/email_checkers", userHandler.CheckEmailAvailability) // note: Check email is available
	api.POST("/avatars", userHandler.UploadAvatar)                  // note: Upload avatar profile

	router.Run()
}
