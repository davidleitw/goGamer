package main

import (
	"fmt"
	"goGamer/gamer"
	"time"
)

const baseurl = "https://forum.gamer.com.tw/C.php?page=2&bsn=60076&snA=3146926"

func main() {
	s := time.Now()
	f, _ := gamer.FindAllFloor("BAHAMUT000", baseurl)
	f.GetResult()
	//gamer.SingleTest(baseurl)
	fmt.Println(time.Since(s))
}
