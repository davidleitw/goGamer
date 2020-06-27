<h1 align="center">goGamer</h1>

<p align="center">
<a href="https://www.gnu.org/licenses/"> <img src="https://img.shields.io/github/license/davidleitw/goGamer.svg" alt="License"></a>
 <a href="http://hits.dwyl.com/davidleitw/goGame">
<img src=http://hits.dwyl.com/davidleitw/goGame.svg alt="HitCount">
</a>
<a href="https://github.com/davidleitw/goGamer/stargazers"> <img src="https://img.shields.io/github/stars/davidleitw/goGamer" alt="GitHub stars"></a>
</p>

---
巴哈姆特找樓工具, 使用Go語言的goroutine加速, 兩萬樓大約7秒左右即可完成
為了避免造成巴哈姆特伺服器的負擔, 請不要重複的查找

本工具預計會提供給一般使用者使用的執行程式(.exe)以及給程式開發者所使用的API
若之後完成API的部份也希望有熱心人士可以寫一個簡單的前端接起來更加方便使用者使用XD

下面寫的文檔基本上都是初步的用法, 一些進階的用法或者更加人性化的用法會慢慢更新
如果有人有什麼更好的想法或者架構也歡迎討論喔

### 一般使用者
---

若是有Go語言的環境
可以將程式碼clone下來
```shell
user@user:~$ git clone https://github.com/davidleitw/goGamer
```
並且將main.go打包成執行檔以便後續的操作
```shell
user@user:~$ go build main.go
```
執行exe檔案, 指令中的雙引號也要加喔
```shell
user@user:~$ ./main -url="想要查詢的討論串網址" -userID="想要查詢的使用者ID"
```

範例 找薯條串中ID為"leichitw"的發文樓層以及內容
```shell 
user@user:~$ git clone https://github.com/davidleitw/goGamer
user@user:~$ go build main.go
user@user:~$ ./main -url="https://forum.gamer.com.tw/C.php?page=2&bsn=60076&snA=3146926" -userID="leichitw"

2020/06/28 03:31:32 81 >>  8337樓 ID=leichitw Name=驥哥
2020/06/28 03:31:32 
大家有聽過一些線上課程的經驗嗎? 大家都是怎麼聽得呢? 這幾天在看台大ㄉ機器學習的課，每次都要暫停理解老師在說甚麼然後抄筆記QQ，深刻的體會了自己的不足，不知道在現場上課的時候那些學生是不是都可以直接聽懂R

2020/06/28 03:31:32 82 >>  8420樓 ID=leichitw Name=驥哥
2020/06/28 03:31:32 
上完台大ㄉ機器學習之後..我感覺已經回不去今天上了自己學校的機器學習就在想，如果沒看過台大的課我這鬼東西怎麼聽的懂

2020/06/28 03:31:32 83 >>  8464樓 ID=leichitw Name=驥哥
2020/06/28 03:31:32 
大家第一次都怎麼讀資料結構那本厚厚的聖經呢?慢慢看我感覺看了好多不必要的東西，跳著看又感覺會漏看很多QQ

2020/06/28 03:31:32 84 >>  8469樓 ID=leichitw Name=驥哥
2020/06/28 03:31:32 
大家github都放些什麼R，昨天我上傳了第一個檔案，目前只會這個還不會分支或版本控制，想說之後來認真經營一下XD
```


