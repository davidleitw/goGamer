package api

import (
	"goGamer/gamer"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchWithTitle struct {
	BaseUrl string `json:"baseurl"`
	KeyWord string `json:"search_title"`
}

func SearchwithTitle(ctx *gin.Context) {
	var servicer SearchWithTitle
	if err := ctx.BindJSON(&servicer); err == nil {
		// 多條件查找
		if strings.Contains(servicer.KeyWord, "&&") {
			sp := strings.Split(servicer.KeyWord, "&&")
			key0 := sp[0]
			key1 := sp[1]
			p1, err1 := gamer.SearchSpecifideTitle(servicer.BaseUrl, key0)
			p2, err2 := gamer.SearchSpecifideTitle(servicer.BaseUrl, key1)
			// 對兩次查詢的結果取交集
			result := gamer.Intersection(p1, p2)
			if err1 == nil && err2 == nil {
				ctx.JSON(http.StatusOK, gin.H{
					"status": http.StatusOK,
					"data":   result.GetResult(),
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status": 500,
					"error":  "伺服器在處理request的時候發生了錯誤, 請稍後再測試",
				})
			}

		} else {
			// 單條件查詢
			p, err := gamer.SearchSpecifideTitle(servicer.BaseUrl, servicer.KeyWord)
			if err == nil {
				ctx.JSON(http.StatusOK, gin.H{
					"status": http.StatusOK,
					"data":   p.GetResult(),
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status": 500,
					"error":  "伺服器在處理request的時候發生了錯誤, 請稍後再測試",
				})
			}
		}
	} else {
		// json綁定失敗
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  "請確認一下傳入的資料有沒有符合api的格式",
		})
	}
}
