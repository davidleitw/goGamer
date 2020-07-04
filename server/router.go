package server

import (
	"net/http"

	"goGamer/api"

	"github.com/gin-gonic/gin"
)

func setCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func NewRouter() *gin.Engine {
	// 設置debug模式
	gin.SetMode(gin.DebugMode)
	server := gin.Default()
	server.StaticFS("/searchUser", http.Dir("./frontend/searchUser"))

	server.POST("/FindAllFloor", api.FindAllFloor)
	server.POST("/FindAuthorFloor", api.FindAuthorFloor)
	server.POST("/FindAllFloorInfo", api.FindAllFloorInfo)
	server.POST("/SearchwithTitle", api.SearchwithTitle)
	server.GET("/FindUserInfo", api.FindUserInfo)
	server.GET("/", func(c *gin.Context) { // "/" 自動導向 "/searchUser"
		c.Redirect(http.StatusMovedPermanently, "/searchUser")
	})
	return server
}
