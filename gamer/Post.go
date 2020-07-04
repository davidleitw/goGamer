package gamer

import (
	"log"
	"sync"
)

type Post struct {
	SubBsn    string   // 子版分類
	SummaryGP int      // 文章總獲得GP數目
	Href      string   // 文章超連結
	Title     string   // 文章標題
	Author    UserInfo // 作者資料
}

type Posts struct {
	list  []Post
	total int
}

// 將post的子集合新增到Posts裡面
func (Ps *Posts) AppendPostSet(p []Post) {
	for _, val := range p {
		Ps.list = append(Ps.list, val)
		Ps.total++
	}
}

func (Ps Posts) Result() {
	for index, article := range Ps.list {
		log.Printf("%4d => %s\n", index+1, article.Title)
	}
}

func (Ps Posts) GetResult() []Post {
	return Ps.list
}

// 多條件查詢的時候取交集
func Intersection(p1, p2 Posts) Posts {
	var intersection []Post
	var total int = 0
	wg := new(sync.WaitGroup)
	wg.Add(len(p1.list))
	for _, s1 := range p1.list {
		go func(s1 Post, _p2 Posts) {
			defer wg.Done()
			for _, s2 := range p2.list {
				if s1.Title == s2.Title {
					intersection = append(intersection, s1)
					break
				}
			}
		}(s1, p2)
	}
	wg.Wait()
	return Posts{
		list:  intersection,
		total: total,
	}
}
