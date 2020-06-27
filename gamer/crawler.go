package gamer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 獲得完整一頁的HTML檔案 以字串表示
func getPageBody(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("錯誤, 請確認您輸入的網址是否正確, 錯誤網址為: %s\n", url)
		return "", nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Error, Status code is ", res.StatusCode)
		return "", errors.New("Status code is not 200!")
	}
	bodyByte, _ := ioutil.ReadAll(res.Body)
	return string(bodyByte), nil
}

func GetPostHeader(body string) {
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(body))
	//div.c-post__header__author>a
	dom.Find("div.c-section__main").Each(func(idx int, selection *goquery.Selection) {
		//fmt.Println(selection.Text())

		selection.Find("div.c-post__header__author>a.floor").Each(func(idx int, selection *goquery.Selection) {
			fmt.Println(selection.Text())
		})

		selection.Find("div.c-post__header__author>a.userid").Each(func(idx int, selection *goquery.Selection) {
			fmt.Println("UserID: ", selection.Text())
		})

		selection.Find("div.c-post__header__author>a.username").Each(func(idx int, selection *goquery.Selection) {
			fmt.Println("UserName: ", selection.Text())
		})

		selection.Find("div.c-article__content>div").Each(func(idx int, selection *goquery.Selection) {
			s := strings.Split(selection.Text(), "\n")
			for _, val := range s {
				fmt.Println(val)
			}
		})
	})
}
