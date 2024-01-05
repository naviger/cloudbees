package main

import (
	"cloudbees-api/controllers"
	"cloudbees-api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Res401Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"401"`
	Message  string `json:"message" example:"authorisation failed"`
}

// claims component of jwt contains mainy fields , we need only roles of DemoServiceClient
// "DemoServiceClient":{"DemoServiceClient":{"roles":["pets-admin","pet-details","pets-search"]}},
type Claims struct {
	ResourceAccess client `json:"resource_access,omitempty"`
	JTI            string `json:"jti,omitempty"`
}

type client struct {
	DemoServiceClient clientRoles `json:"DemoServiceClient,omitempty"`
}

type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

var RealmConfigURL string = "https://isperience.web:4443/auth/realms/iSperience"
var clientID string = "isperience_client"

func main() {

	//API AND ROUTING SETUP >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	api := gin.Default()

	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "content-type", "accept", "authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api.GET("/train", middleware.SetAuthorizedJWT(false), controllers.GetTrains())
	api.GET("/train/:trainId", middleware.SetAuthorizedJWT(false), controllers.GetTrain())
	api.POST("/train/:trainId/seat", middleware.SetAuthorizedJWT(false), controllers.RequestSeat())
	api.GET("/train/:trainId/seat/:seatId/receipt", middleware.SetAuthorizedJWT(true), controllers.RequestReceipt())
	api.PATCH("/train/:trainId/seat/:seatId/change", middleware.SetAuthorizedJWT(true), controllers.ChangeSeat())
	api.PATCH("/train/:trainId/seat/:seatId/cancel", middleware.SetAuthorizedJWT(true), controllers.CancelSeat())

	api.RunTLS(":5007", "./certs/localhost.crt", "./certs/localhost.key")
}
