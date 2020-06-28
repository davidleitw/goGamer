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

type FindInfoer struct {
	BaseUrl string `json:"baseurl"`
}

func FindAllFloorInfo(ctx *gin.Context) {
	var servicer FindInfoer
	// 確認資料有正確綁定
	if err := ctx.BindJSON(&servicer); err == nil {
		result, err := gamer.FindAllFloorInfo(servicer.BaseUrl)
		// 確認服務本身是否有正常運行
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"data":   result.GetFloors(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"error":  "伺服器在處理request的時候發生了錯誤, 請稍後再測試",
			})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  "請確認一下傳入的資料有沒有符合api的格式",
		})
	}
}

func FindAllFloor(ctx *gin.Context) {
	var servicer Finder
	// 資料有正確綁定
	if err := ctx.BindJSON(&servicer); err == nil {
		// fmt.Println("url = ", servicer.BaseUrl, "UserID = ", servicer.UserID)
		result, err := gamer.FindAllFloor(servicer.BaseUrl, servicer.UserID)
		// 確認服務本身是否有正常運行
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"data":   result.GetFloors(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"error":  "伺服器在處理request的時候發生了錯誤, 請稍後再測試",
			})
		}
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
		result, err := gamer.FindAuthorFloor(servicer.BaseUrl, servicer.UserID)
		// 確認服務本身是否有正常運行
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"data":   result.GetFloors(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"error":  "伺服器在處理request的時候發生了錯誤, 請稍後再測試",
			})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  "請確認一下傳入的資料有沒有符合api的格式",
		})
	}
}
