package gamer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

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

// 輸入用戶ID與想要爬取的討論串, 就會將所有結果放進FloorSet並且回傳
func FindAllFloor(userid string, baseurl string) (FloorSet, error) {
	var Fs FloorSet
	// 獲得討論串每一頁的連結(一頁總共20層樓)
	urls, err := getUrlSet(baseurl)
	fmt.Printf("total floor number are %d\n", len(urls))
	if err != nil {
		return Fs, err
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(urls))

	// 對於每個頁的連結去get其html, 並且用goquery分析
	for _, url := range urls {
		go func() {
			f := handle(url, userid, wg)
			// 將樓層資訊彙整到Floor set裡面
			if len(f) >= 1 {
				Fs.AddFloors(f)
			}
		}()
		time.Sleep(2500 * time.Microsecond)
	}
	wg.Wait()
	return Fs, nil
}

//只找使用者在討論串的樓(無法獲得實際在討論串樓層數)
func FindAuthorFloor(baseurl, userID string) (FloorSet, error) {
	var Fs FloorSet
	// 獲得討論串每一頁的連結(一頁總共20層樓)
	urls, err := getAuthorUrlSet(baseurl, userID)
	fmt.Printf("total floor number are %d\n", len(urls))
	if err != nil {
		return Fs, err
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(urls))

	// 對於每個頁的連結去get其html, 並且用goquery分析
	for _, url := range urls {
		go func() {
			f := handle(url, userID, wg)
			// 將樓層資訊彙整到Floor set裡面
			if len(f) >= 1 {
				Fs.AddFloors(f)
			}
		}()
		time.Sleep(2500 * time.Microsecond)
	}
	wg.Wait()
	return Fs, nil
}

// 爬蟲主體, 爬完之後把每一層樓的資料放在一個Floor陣列傳回
func handle(url string, userID string, wg *sync.WaitGroup) []Floor {
	var fs []Floor
	// 由url轉換成html文檔
	html, _ := getPageBody(url)
	// 解析html文檔來找尋特定的使用者
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	dom.Find("div.c-section__main").Each(func(idx int, selection *goquery.Selection) {
		var f Floor
		var found bool = false
		selection.Find("div.c-post__header__author").Each(func(idx int, s1 *goquery.Selection) {
			ID := s1.Find("a.userid").First().Text()
			// 如果某一層樓的userID跟目標ID相同的話, 將其記錄下來
			if ID == userID {
				// 設置樓層資訊中的樓層數, 用戶ID, 用戶名稱
				f.SetInfo(getFloorNum(s1.Find("a.floor").First().Text()),
					s1.Find("a.username").First().Text(),
					s1.Find("a.userid").First().Text())
				found = true
			}
		})
		if found {
			selection.Find("div.c-article__content").Each(func(idx int, s1 *goquery.Selection) {
				// 將空格刪掉
				s1.Remove()
				f.Setcontent(s1.Text())
				fs = append(fs, f)
			})
		}
	})
	defer wg.Done()
	return fs
}

func SingleTest(url string) {
	html, _ := getPageBody(url)
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	dom.Find("div.c-section__main").Each(func(idx int, selection *goquery.Selection) {
		selection.Find("div.c-post__header__author").Each(func(idx int, s1 *goquery.Selection) {
			fmt.Println(s1.Text())
		})
		selection.Find("div.c-article__content").Each(func(idx int, s1 *goquery.Selection) {
			s1.Remove()
			fmt.Println(s1.Text())
		})
	})
}
