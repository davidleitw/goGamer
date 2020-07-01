package gamer

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Find開頭的函數都是爬蟲的主要部份

// 獲得單一用戶的帳號資訊
func FindUserInfo(UserID string) (UserInfo, error) {
	var user UserInfo
	// 獲得該用戶的小屋網址
	baseurl := fmt.Sprintf("https://home.gamer.com.tw/homeindex.php?owner=%s", UserID)

	// 由url轉換成html文檔
	html, err := getPageBody(baseurl)
	if err != nil {
		return user, err
	}
	// 解析html文檔來找尋特定的使用者
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return user, err
	}

	userNode := dom.Find("div#BH-slave>div.MSG-list2").First()
	userNode.Find("li").EachWithBreak(func(idx int, Selection *goquery.Selection) bool {
		switch idx {
		case 0:
			// 在巴哈的Html裡面完整形式: 帳號：leichitw  注意! 中間的冒號是全形的..
			userID := strings.Split(Selection.Text(), "：")[1]
			user.UserID = userID
		case 1:
			userName := strings.Split(Selection.Text(), "：")[1]
			user.UserName = userName
		case 2:
			title := strings.Split(Selection.Text(), "：")[1]
			user.Title = title
		case 3:
			info := strings.Split(Selection.Text(), "/")
			// strings.Replace(xx, " ", "", -1) => 把空格過濾掉
			// 等級以整數表示
			level, _ := strconv.Atoi(strings.Replace(info[0][2:], " ", "", -1))
			race := strings.Replace(info[1], " ", "", -1)
			occu := strings.Replace(info[2], " ", "", -1)

			user.Level = level
			user.Race = race
			user.Occupation = occu
		case 4:
			balance, _ := strconv.Atoi(strings.Split(Selection.Text(), "：")[1]) // 以整數表示
			user.Balance = balance
		case 5:
			gp, _ := strconv.Atoi(strings.Split(Selection.Text(), "：")[1])
			user.GP = gp
		default:
			// 讀到gp完就中斷遍歷
			return false
		}
		return true
	})
	return user, nil
}

func FindAllArticleTitle(baseurl string, start, end int) {

}

func FindAllFloorInfo(baseurl string) (FloorSet, error) {
	var Fs FloorSet
	urls, err := getUrlSet(baseurl)
	if err != nil {
		return Fs, err
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(urls))

	for _, url := range urls {
		go func() {
			fs := handleFindAllInfo(url, wg)
			Fs.AddFloors(fs)
		}()
		time.Sleep(25000 * time.Microsecond)
	}
	wg.Wait()
	Fs.SortResult()
	return Fs, nil
}

// 輸入用戶ID與想要爬取的討論串, 就會將所有結果放進FloorSet並且回傳
func FindAllFloor(baseurl, userID string) (FloorSet, error) {
	var Fs FloorSet
	// 獲得討論串每一頁的連結(一頁總共20層樓)
	urls, err := getUrlSet(baseurl)
	if err != nil {
		return Fs, err
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(urls))

	// 對於每個頁的連結去get其html, 並且用goquery分析
	for _, url := range urls {
		go func() {
			f := handleFindUser(url, userID, wg)
			// 將樓層資訊彙整到Floor set裡面
			if len(f) >= 1 {
				Fs.AddFloors(f)
			}
		}()
		// 避免過於頻繁的get, 導致request被擋下來
		time.Sleep(25000 * time.Microsecond)
	}
	wg.Wait()
	Fs.SortResult()
	return Fs, nil
}

// 只找使用者在討論串的樓(無法獲得實際在討論串樓層數)
func FindAuthorFloor(baseurl, userID string) (FloorSet, error) {
	var Fs FloorSet
	// 獲得討論串每一頁的連結(一頁總共20層樓)
	urls, err := getAuthorUrlSet(baseurl, userID)
	if err != nil {
		return Fs, err
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(urls))

	// 對於每個頁的連結去get其html, 並且用goquery分析
	for _, url := range urls {
		go func() {
			f := handleFindUser(url, userID, wg)
			// 將樓層資訊彙整到Floor set裡面
			if len(f) >= 1 {
				Fs.AddFloors(f)
			}
		}()
		time.Sleep(2500 * time.Microsecond)
	}
	wg.Wait()
	Fs.SortResult()
	return Fs, nil
}

// 爬蟲主體, 爬完之後把每一層樓的資料放在一個Floor陣列傳回
func handleFindUser(url string, userID string, wg *sync.WaitGroup) []Floor {
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
					ID)
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

// 不過濾使用者, 爬取每一層樓
func handleFindAllInfo(url string, wg *sync.WaitGroup) []Floor {
	var fs []Floor
	defer wg.Done()
	html, _ := getPageBody(url)
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	dom.Find("div.c-section__main").Each(func(idx int, selection *goquery.Selection) {
		var f Floor
		selection.Find("div.c-post__header__author").Each(func(idx int, s1 *goquery.Selection) {
			// 設置樓層資訊中的樓層數, 用戶ID, 用戶名稱
			f.SetInfo(getFloorNum(s1.Find("a.floor").First().Text()),
				s1.Find("a.username").First().Text(),
				s1.Find("a.userid").First().Text())
		})
		selection.Find("div.c-article__content").Each(func(idx int, s1 *goquery.Selection) {
			// 將空格刪掉
			s1.Remove()
			f.Setcontent(s1.Text())
			fs = append(fs, f)
		})
	})
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
