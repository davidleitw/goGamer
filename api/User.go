package api

import (
	"goGamer/gamer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindUserInfo(ctx *gin.Context) {
	userID := ctx.Query("ID")
	userInfo, err := gamer.FindUserInfo(userID)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   userInfo,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"error":  "伺服器在處理request的時候發生了錯誤, 請稍後再測試",
		})
	}
}
