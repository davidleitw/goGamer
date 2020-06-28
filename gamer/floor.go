package gamer

import (
	"log"
	"sort"
)

var count int = 0

type Floor struct {
	Num      int    // 樓層數
	UserName string // 用戶名稱
	UserID   string // 用戶帳號
	Content  string // 樓層主體
}

func (f *Floor) SetInfo(num int, name string, id string) {
	f.Num = num
	f.UserName = name
	f.UserID = id
}

func (f *Floor) Setcontent(content string) {
	f.Content = content
}

func (f *Floor) GetNum() int {
	return f.Num
}

func (f *Floor) GetuserName() string {
	return f.UserName
}

func (f *Floor) GetuserID() string {
	return f.UserID
}

func (f *Floor) GetContent() string {
	return f.Content
}

type FloorSet struct {
	floors []Floor
	total  int
}

func (Fs *FloorSet) GetOneFloor(index int) Floor {
	return Fs.floors[index]
}

func (Fs *FloorSet) GetFloors() []Floor {
	return Fs.floors
}

func (Fs *FloorSet) SortResult() {
	// 依照樓層排序
	sort.SliceStable(Fs.floors, func(i, j int) bool {
		return Fs.floors[i].Num < Fs.floors[j].Num
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
