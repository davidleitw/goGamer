package main

import (
	"github.com/davidleitw/goGamer/gamer"
)

func main() {
	// urlPtr := flag.String("url", "", "想要搜尋的討論串網址(哪一層樓的網址都可以)")
	// idPtr := flag.String("userID", "", "想要搜尋的使用者ID")
	// flag.Parse()
	// s := time.Now()
	// // 比起原先的方法更加的快速
	// f, _ := gamer.FindAuthorFloor(*urlPtr, *idPtr)
	// f.GetResult()
	// fmt.Println(time.Since(s))
	Ps, _ := gamer.SearchSpecifideTitle("https://forum.gamer.com.tw/C.php?bsn=30861&snA=24342&tnum=1&subbsn=1", "推薦")
	Ps.Result()
}
