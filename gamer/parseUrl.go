package gamer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// const baseurl = "https://forum.gamer.com.tw/C.php?page=1&bsn=60076&snA=3146926"

// 獲得一個討論串每一層樓的連結
func getUrlSet(baseurl string) ([]string, error) {
	var result []string
	// Query之前的部份 ex: https://forum.gamer.com.tw/C.php
	front := strings.Split(baseurl, "?")[0]
	max, _ := getMaxFloorNumber(baseurl)
	page := int(max/20) + 1

	parameters, _ := getFloorInfo(baseurl)
	for i := 1; i <= page; i++ {
		result = append(result, fmt.Sprintf("%s?page=%d&bsn=%s&snA=%s", front, i, parameters[0], parameters[1]))
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

// 剛開始要獲得一個討論串url中bsn跟snA參數
func getFloorInfo(baseurl string) ([]string, error) {
	parse, err := url.Parse(baseurl)
	if err != nil {
		return []string{}, err
	}
	// 獲得文章的ID
	values, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return []string{}, err
	}
	bsn := values.Get("bsn")
	snA := values.Get("snA")
	return []string{bsn, snA}, nil
}

//指定userID找樓
func getAuthorUrlSet(baseurl, userID string) ([]string, error) {
	var result []string
	fmt.Println(baseurl, userID)
	front := strings.Split(baseurl, "?")[0]
	max, _ := getAuthorMaxFloorNumber(baseurl, userID)
	page := int(max/20) + 1

	bsn, snA, _ := getFloorInfoPlus(baseurl)

	for i := 1; i <= page; i++ {
		result = append(result, fmt.Sprintf("%s?page=%d&bsn=%s&snA=%s&s_author=%s", front, i, bsn, snA, userID))
	}
	return result, nil
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

//info改寫
func getFloorInfoPlus(baseurl string) (string, string, error) {
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
