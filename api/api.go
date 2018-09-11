package api

import (
	"github.com/gin-gonic/gin"
	"github.com/saisai/gindemo/api/controllers"
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
	//	v1.POST("/login")
	//	v1.GET("/info")
}
