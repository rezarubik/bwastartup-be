package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	campaignRepository := campaign.NewRepository(db)

	// todo: Test repository
	// campaigns, err := campaignRepository.FindByUserID(3)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()

	// todo: Testing create campaign directly (example)
	// inputCampaign := campaign.CreateCampaignInput{}
	// inputCampaign.Name = "Pengadaan Campaign"
	// inputCampaign.ShortDescription = "Short Descriptin"
	// inputCampaign.Description = "Long Description"
	// inputCampaign.Perks = "Perks satu, dua, tiga"
	// inputCampaign.GoalAmount = 1000000

	// inputUser, _ := userService.GetUserByID(3)
	// inputCampaign.User = inputUser

	// _, err = campaignService.CreateCampaign(inputCampaign)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	log.Fatal()
	// }

	// fmt.Println(authService.GenerateToken(1001))

	// userService.SaveAvatar(3, "images/1-profile.png")

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)                                             //note: Register
	api.POST("/sessions", userHandler.Login)                                                 //note: Login
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)                          // note: Check email is available
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar) // note: Upload avatar profile

	/* --------------------------- // start: Campaign --------------------------- */
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)
	/* ---------------------------- // end: Campaign ---------------------------- */

	router.Run()
}

// note: Auth Middleware
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo: Ambil nilai header Authorization: Bearer tokentokentoken
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// todo: Dari header Authorization, kita ambil tokennya
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		// todo: Validasi Token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// todo: Ambil user dari db berdasarkan user_id lewat service
		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// todo: kita set context isinya user
		c.Set("currentUser", user)
	}
}
