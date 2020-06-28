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
API的部份是以gin當作伺服器框架, 到時候可以會部屬到Heroku上面

下面寫的文檔基本上都是初步的用法, 一些進階的用法或者更加人性化的用法會慢慢更新
如果有人有什麼更好的想法或者架構也歡迎討論喔

開頭的部份會介紹一下巴哈的哈拉版架構以及網址的規則
如果是一般的使用者請直接跳至下方


### 目錄
* [巴哈姆特討論串網址分析](#巴哈姆特討論串網址分析)
* [討論串樓層Html解析](#討論串樓層Html架構分析)
* [爬蟲思路](#整體程式架構概述)
* [一般使用者](#直接clone下來的作法)
* [提供API](#方便程式開發者運用爬蟲去爬取特定的資料)

--- 

### 巴哈姆特討論串網址分析
![](https://imgur.com/qah05AL.png)
上圖是我從資工討論串中第2041頁擷取下來的網址(我自己的設定是每頁有十層樓)
來簡單說一下各個參數代表的意義
* page &ensp; => 目前在這個討論串中的第幾頁
* bsn &ensp;&ensp; => 代表哈拉版的號碼(場外＝60076, 公主連結=30861..等等) 每一個bsn都可以決定唯一的哈拉版
* snA&ensp;&ensp;&ensp;=> 代表著這篇文章在該版的文章編號

所以說範例網址就是在60076號哈拉版中的3146926號文章
不過實際代表意義我們不得而知, 有可能是採用不重複的唯一值代表文章編號, 不一定有排序性

---

### 討論串樓層Html架構分析
- **整體架構**

爬蟲的第一步往往都是按下F12觀察網頁的架構
這邊先來討論每層樓的樓層數, 用戶名稱以及用戶帳號放在html的哪個部份
![](https://imgur.com/pU5VAMS.png)
若仔細觀察會發現討論串中的每層樓都是以
```html
<div class="c-section__main c-post ">
```
為主體, 把整層樓包含留言包起來

其中每層樓又包含著三個部份
```html
<div class="c-post__header">         
<div class="c-post__body">
<div class="c-post__footer c-reply">
```
所以每一層樓的架構是這個樣子
```html
<div class="c-section__main c-post ">               // 每一層樓都用c-section__main包起來
    <div class="c-post__header"></div>              // 主要放樓層數, 以及樓主的資料, 還有文章的基本訊息
    <div class="c-post__body"></div>                // 文章內容
    <div class="c-post__footer c-reply"></div>      // 留言區
</div>
```
我們初步的目的是要用來找樓以及內容, 所以一開始只要著重在header跟body的處理即可

#### **c-post__header(存放樓層數以及樓主資料)**
```html
<div class="c-post__header"> 
    <div class="c-post__header__tools"></div>      放置開圖工具
    <div class="c-post__header__author">...</div>  樓層數, 樓主帳戶名稱等資料
    <div class="c-post__header__info">...</div>    發文時間, 發送ip等等
</div>
```

- **c-post_header__tools** 

![](https://imgur.com/1Q9W4xd.png)
c-post__header__tools這個區域就是在每頁的第一樓存放討論串標題跟開圖工具的區域(如上圖所示)
除了第一樓之外通常這個區域都是空的
<br>

- **c-post_header__author**

![](https://imgur.com/QzwFkzL.png)
這部份就是我們爬蟲的第一個重點, 這個區域能獲取的資訊有樓層數, 用戶ID以及名稱, 最後還有推噓的數目

因為這邊的資訊相較一下欄位比較沒有這麼複雜, 所以說爬起來相對的簡單
<br>

- **c-post_header__info**

![](https://imgur.com/ygZ5ghd.png)
info區塊就是存放一些發文時的資訊(發文時間, ip位置等等..)

#### **c-post__body**(存放文章內容)

```html
c-post__body架構
<div class="c-post_body"> 
    <article class="c-article FM-P2" id="xxxxxxxx">... </article> // 存放文章內容
    <div class="c-post__body__buttonbae">                         // 工具列(贊, 噓, 回覆等等..)
</div>
```

下圖資工串的第一樓, 拿阿條的文章來做範例
![](https://imgur.com/cvkiRGD.png)

我們爬蟲通常要爬的就是`<article>`的區塊了, 因為如果要爬gp,bp的話header的部份也可以獲得

展開`<article>`區塊來觀察看看, 如下圖所示

![](https://imgur.com/igPLMKc.png)

如果這篇文章沒有插入圖片影片或者超連結, 就會是個單純的div把文章內容都包起來

```html
<div class="c-article__content">
    文章內容
</div>
```

我們爬蟲就會針對上述區塊去獲取我們想要的文章內容


---

### 一般使用者
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
    .
    .
    . 
    略
```


