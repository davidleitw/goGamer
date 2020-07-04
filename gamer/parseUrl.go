package gamer

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// const baseurl = "https://forum.gamer.com.tw/C.php?page=1&bsn=60076&snA=3146926"

// 用Get的方式取得指定網址的html文檔, 並且轉換成goquery用來檢索的strcut
func getDecument(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("錯誤, 請確認您輸入的網址是否正確, 錯誤網址為: %s\n", url)
		return nil, nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Error, Status code is ", res.StatusCode)
		return nil, errors.New("Status code is not 200!")
	}
	bodyByte, _ := ioutil.ReadAll(res.Body)

	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyByte))
	if err != nil {
		return nil, err
	}
	return dom, nil
}

//指定userID找樓
func getAuthorUrlSet(baseurl, userID string) ([]string, error) {
	var result []string
	fmt.Println(baseurl, userID)
	front := strings.Split(baseurl, "?")[0]
	max, _ := getAuthorMaxFloorNumber(baseurl, userID)
	page := int(max/20) + 1

	bsn, snA, _ := getFloorInfo(baseurl)
	for i := 1; i <= page; i++ {
		result = append(result, fmt.Sprintf("%s?page=%d&bsn=%s&snA=%s&s_author=%s", front, i, bsn, snA, userID))
	}
	return result, nil
}

// 獲得一個討論串每一層樓的連結
func getUrlSet(baseurl string) ([]string, error) {
	var result []string
	// Query之前的部份 ex: https://forum.gamer.com.tw/C.php
	front := strings.Split(baseurl, "?")[0]
	max, _ := getMaxFloorNumber(baseurl)
	page := int(max/20) + 1

	bsn, snA, _ := getFloorInfo(baseurl)
	for i := 1; i <= page; i++ {
		result = append(result, fmt.Sprintf("%s?page=%d&bsn=%s&snA=%s", front, i, bsn, snA))
	}

	return result, nil
}

func getSearchUrl(page int, bsn, search string) string {
	return fmt.Sprintf("https://forum.gamer.com.tw/B.php?page=%d&bsn=%s&qt=1&q=%s", page, bsn, search)
}

// 如果要找文章, 需要尋找搜尋結果總共有幾頁
func getMaxPosterNumber(baseurl string) (int, error) {
	dom, err := getDecument(baseurl)
	if err != nil {
		return -1, err
	}
	// 以整數形式表示共有幾頁的搜尋結果
	max, err := strconv.Atoi(dom.Find("p.BH-pagebtnA>a").Last().Text())
	if err == nil {
		return max, nil
	} else {
		return -1, err
	}
}

// 獲得一串討論串最高的樓層數
func getMaxFloorNumber(baseurl string) (int, error) {
	front := strings.Split(baseurl, "?")[0]
	parse, err := url.Parse(baseurl)
	if err != nil {
		return -1, err
	}
	// 獲得文章的ID
	values, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return -1, err
	}

	bsn := values.Get("bsn")
	snA := values.Get("snA")

	target := fmt.Sprintf("%s?page=%d&bsn=%s&snA=%s", front, 999999, bsn, snA)

	dom, _ := getDecument(target)
	s := dom.Find("div.c-post__header__author>a.floor").Last()
	n := getFloorNum(s.Text())
	return n, nil
}

//info改寫
func getFloorInfo(baseurl string) (string, string, error) {
	parse, err := url.Parse(baseurl)
	if err != nil {
		return "", "", err
	}

	values, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return "", "", err
	}
	bsn, snA := values.Get("bsn"), values.Get("snA")
	return bsn, snA, nil
}

func getBsn(baseurl string) (string, error) {
	parse, err := url.Parse(baseurl)
	if err != nil {
		return "", err
	}

	values, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return "", err
	}
	bsn := values.Get("bsn")
	return bsn, nil
}

// 查詢單一query
func getQuery(baseurl string, parameter string) (string, error) {
	parse, err := url.Parse(baseurl)
	if err != nil {
		return "", err
	}

	values, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return "", err
	}

	return values.Get(parameter), nil
}

func getAuthorMaxFloorNumber(baseurl, userID string) (int, error) {
	front := strings.Split(baseurl, "?")[0]
	parse, err := url.Parse(baseurl)
	if err != nil {
		return -1, err
	}
	// 獲得文章的ID
	values, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return -1, err
	}
	bsn := values.Get("bsn")
	snA := values.Get("snA")

	target := fmt.Sprintf("%s?page=%d&bsn=%s&snA=%s&s_author=%s", front, 999999, bsn, snA, userID)

	dom, _ := getDecument(target)
	s := dom.Find("div.c-post__header__author>a.floor").Last()
	n := getFloorNum(s.Text())
	return n, nil
}

// 將爬取到的樓層字串轉換成整數
// Ex 輸入 "20388 樓"  輸出 20388
func getFloorNum(floor string) int {
	if floor == "樓主" {
		return 1
	}
	if strings.ContainsAny(floor, "樓") {
		n, _ := strconv.Atoi(strings.Split(floor, " ")[0])
		return n
	}
	return -1
}
