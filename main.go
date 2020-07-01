package main

import (
	"fmt"
	"goGamer/gamer"
	"time"
)

func main() {
	s:=time.Now()
	f, _ := gamer.FindAllFloor("relaxplay", "https://forum.gamer.com.tw/C.php?page=2&bsn=60076&snA=3146926")
	f.GetResult()
	fmt.Println(time.Since(s))
}
