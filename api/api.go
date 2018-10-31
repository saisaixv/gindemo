package api

import (
	"github.com/gin-gonic/gin"
	"github.com/saisai/gindemo/api/controllers"

	//	"github.com/dchest/captcha"
)

var (
	engine *gin.Engine
)

func init() {
	engine = gin.Default()
	setupRoutersV1()
}

func Engine() *gin.Engine {
	return engine
}

func setupRoutersV1() {
	v1 := engine.Group("/usersystem/api/v1")
	v1.POST("/register", controllers.Register)
	v1.POST("/login", controllers.Login)
	v1.POST("/logout", controllers.Logout)
	v1.GET("/info", controllers.Info)
	v1.POST("/add_identify_type", controllers.AddIdentifyType)
	v1.POST("/authentication", controllers.Authentication)

	private := engine.Group("/private/api/v1")
	private.POST("/private_register", controllers.Private_Register)

}
