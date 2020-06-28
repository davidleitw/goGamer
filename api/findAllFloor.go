package api

import (
	"goGamer/gamer"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Finder struct {
	BaseUrl string `json:"baseurl"`
	UserID  string `json:"userID"`
}

func FindAllFloor(ctx *gin.Context) {
	var servicer Finder
	// 資料有正確綁定
	if err := ctx.BindJSON(&servicer); err == nil {
		// fmt.Println("url = ", servicer.BaseUrl, "UserID = ", servicer.UserID)
		result, _ := gamer.FindAllFloor(servicer.BaseUrl, servicer.UserID)

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   result.GetFloors(),
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  "請確認一下傳入的資料有沒有符合api的格式",
		})
	}
}

func FindAuthorFloor(ctx *gin.Context) {
	var servicer Finder

	if err := ctx.BindJSON(&servicer); err == nil {
		result, _ := gamer.FindAuthorFloor(servicer.BaseUrl, servicer.UserID)

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   result.GetFloors(),
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  "請確認一下傳入的資料有沒有符合api的格式",
		})
	}
}
