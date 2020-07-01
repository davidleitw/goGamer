package gamer

import (
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

	html, _ := getPageBody(target)
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
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

	html, _ := getPageBody(target)
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
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
