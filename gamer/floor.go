package gamer

import (
	"log"
	"sort"
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
	// 依照樓層排序
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
	}
	log.Printf("總共%d層樓\n", Fs.GetTotal())
}
