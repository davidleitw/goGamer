package main

import (
	"flag"
	"fmt"
	"goGamer/gamer"
	"time"
)

func main() {
	urlPtr := flag.String("url", "", "想要搜尋的討論串網址(哪一層樓的網址都可以)")
	idPtr := flag.String("userID", "", "想要搜尋的使用者ID")
	flag.Parse()
	s := time.Now()
	// 比起原先的方法更加的快速
	f, _ := gamer.FindAuthorFloor(*urlPtr, *idPtr)
	f.GetResult()
	fmt.Println(time.Since(s))
	// f, _ := gamer.FindAllFloorInfo("https://forum.gamer.com.tw/C.php?page=70&bsn=30861&snA=18013&tnum=698")
	// f.GetResult()
}
