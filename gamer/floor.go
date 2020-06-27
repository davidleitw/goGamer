package gamer

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var count int = 0

type Floor struct {
	num      int    // 樓層數
	userName string // 用戶名稱
	userID   string // 用戶帳號
	content  string // 樓層主體
}

func (f *Floor) SetInfo(num int, name string, id string) {
	f.num = num
	f.userName = name
	f.userID = id
}

func (f *Floor) Setcontent(content string) {
	f.content = content
}

func (f *Floor) GetNum() int {
	return f.num
}

func (f *Floor) GetuserName() string {
	return f.userName
}

func (f *Floor) GetuserID() string {
	return f.userID
}

func (f *Floor) GetContent() string {
	return f.content
}

type FloorSet struct {
	floors []Floor
	total  int
}

func (Fs *FloorSet) SortResult() {
	sort.SliceStable(Fs.floors, func(i, j int) bool {
		return Fs.floors[i].num < Fs.floors[j].num
	})
}

func (Fs *FloorSet) GetTotal() int {
	return Fs.total
}

func (Fs *FloorSet) AddFloors(f []Floor) {
	for _, val := range f {
		Fs.floors = append(Fs.floors, val)
		Fs.total++
	}
}

func (Fs *FloorSet) GetResult() {
	Fs.SortResult()
	for i := 0; i < len(Fs.floors); i++ {
		log.Printf("%d >> %5d樓 ID=%s Name=%s\n", i+1, Fs.floors[i].GetNum(), Fs.floors[i].GetuserID(), Fs.floors[i].GetuserName())
		log.Println(Fs.floors[i].GetContent())
		time.Sleep(1 * time.Second)
	}
	log.Printf("總共%d層樓\n", Fs.GetTotal())
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
			Fs.AddFloors(f)
		}()
		time.Sleep(1950 * time.Microsecond)
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
				f.SetInfo(getFloorNum(s1.Find("a.floor").First().Text()), s1.Find("a.username").First().Text(), s1.Find("a.userid").First().Text())
				found = true
			}
		})
		if found {
			selection.Find("div.c-article__content").Each(func(idx int, s1 *goquery.Selection) {
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
