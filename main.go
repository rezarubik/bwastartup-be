package main

import (
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/user"
	"fmt"
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
	authService := auth.NewService()

	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozfQ.k98dBbHioAp-jkA4wMzEXlu3Z0U1F93gJdffoY4Ha4o")
	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxN30.L089-fiu9ey1Zy2RTrj-cOnqtHLM8g0M2Qo6aw7HFjk")
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println("ERROR")
		fmt.Println("ERROR")
	}
	if token.Valid {
		fmt.Println("VALID")
		fmt.Println("VALID")
		fmt.Println("VALID")
	} else {
		fmt.Println("INVALID")
		fmt.Println("INVALID")
		fmt.Println("INVALID")
	}

	// fmt.Println(authService.GenerateToken(1001))

	// userService.SaveAvatar(3, "images/1-profile.png")

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)                    //note: Register
	api.POST("/sessions", userHandler.Login)                        //note: Login
	api.POST("/email_checkers", userHandler.CheckEmailAvailability) // note: Check email is available
	api.POST("/avatars", userHandler.UploadAvatar)                  // note: Upload avatar profile

	router.Run()
}
