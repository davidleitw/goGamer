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
* [一般使用者](#一般使用者)
* [API格式](#API格式)
    * [獲得所有樓層資訊](#獲得此討論串所有文章的訊息)
    * [根據使用者ID獲得該使用者在特定討論串的發文紀錄](#找到某個User在某個討論串中的所有文章)
    * [根據使用者ID獲得該使用者在特定討論串的發文紀錄,速度較慢,因為會保留原始樓層](#找到某個User在某個討論串中的所有文章並保留原始樓層)
    * [獲得指定用戶的帳號訊息](#查詢單一用戶資料)
--- 

## 巴哈姆特討論串網址分析
![](https://imgur.com/qah05AL.png)
上圖是我從資工討論串中第2041頁擷取下來的網址(我自己的設定是每頁有十層樓)
來簡單說一下各個參數代表的意義
* page &ensp; => 目前在這個討論串中的第幾頁
* bsn &ensp;&ensp; => 代表哈拉版的號碼(場外＝60076, 公主連結=30861..等等) 每一個bsn都可以決定唯一的哈拉版
* snA&ensp;&ensp;&ensp;=> 代表著這篇文章在該版的文章編號

所以說範例網址就是在60076號哈拉版中的3146926號文章
不過實際代表意義我們不得而知, 有可能是採用不重複的唯一值代表文章編號, 不一定有排序性

---

## 討論串樓層Html架構分析
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

我們爬蟲通常要爬的就是`<article>`的區塊了, 因為如果要爬gp, bp的話header的部份也可以獲得

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

## 一般使用者
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
--- 

## API格式
### 獲得此討論串所有文章的訊息
#### 此方法可以獲得所有樓層的基本資訊
傳遞json參數意義
<br>"baseurl"欄位擺放的是想要查詢的討論串欄位, 值得一提的是不管貼的連結所在的page在哪頁, 都可以藉由api找到整串的資料

#### Request

- Method: **POST**
- Url: ```https://go-gamer.herokuapp.com/FindAllFloorInfo```
- Headers: Content-Type:application/json
- Body:

```json
{
    "baseurl": "https://forum.gamer.com.tw/C.php?bsn=60030&snA=536651&tnum=39"
}
```
#### Response

回傳的data是以樓層(從1樓開始, 失去原始樓層資料)為排列的資料

```json
{
    "status": 200,
    "data":[
         {
            "Num": 1,
            "UserName": "aviry",
            "UserID": "zzzz89755",
            "Content": "\n華碩蘆竹廠維修品質真的太扯了..事情是這樣，小弟在10/29號前往高雄皇家送修主機板（型號：ROG Crosshair Vi Extreme),過幾天來電告知已經送去華碩蘆竹廠進行維修，禮拜三中午，我上網去看進度發現寫著無料，工程師會盡快與您聯繫，因為我禮拜四禮拜五要上課，可能沒時間接電話，所以當下立馬進線客服，請客服轉告工程師聯繫我，後來被放了4次鴿子我就不計較了。昨天工程師告知無法修復，要幫我換一塊同型號的主機板，我今天收到檢查，發現了不少問題（可能有些人會覺得我是在吹毛求疵，可是...一塊好好的去，換一塊這東西回來，真的很幹）1,AM4扣具損毀2,鎖M.2散熱片的螺絲不同3,USB3.0針腳歪斜4，風扇針腳外斜5.Power Butten針腳歪斜這真的是良品嗎？這到底是什麼維修品質？華碩啊華碩，不要太超過，我以後再也不會買，也不敢買了..後續：今天下午13：30一收到發現問題就馬上致電華碩他們卻說針腳是包裝去影響到，USB3.0的則是測試去動到之類的...總之我是覺得有點太扯了...然後我就要求換其他板子，但工程師就一直說他沒有權限，沒有辦法做更換等等...真的是心累啊對了，工程師說要跟高層回報，報到現在都還沒有回覆...針腳歪斜部分是線插不進去才發現的\n"
        },
        {
            "Num": 2,
            "UserName": "Bang你全家v",
            "UserID": "marty11440",
            "Content": "\n幫高調 這太離譜了\n"
        },
        {
            ...
        },
    ]
}
```

錯誤時候的回傳

```json
{
    "status": 400,
    "error":  "請確認一下傳入的資料有沒有符合api的格式",
}
```
```json 
{
    "status": 500, 
    "error": "伺服器在處理request的時候發生了錯誤, 請稍後再測試" 
}
```


---
### 找到某個User在某個討論串中的所有文章
#### 此方法不會獲得原始的樓層, 不過相較於FindAllFloor來說速度更快, 也幾乎不會出錯
#### 感謝FizzyEit好夥伴丟的PR
 
<br>"userID"欄位擺放使用者想要查詢對象的巴哈ID
#### Request

- Method: **POST**
- Url: ```https://go-gamer.herokuapp.com/FindAuthorFloor```
- Headers: Content-Type:application/json
- Body:

```json
{
    "baseurl": "https://forum.gamer.com.tw/C.php?page=2&bsn=60076&snA=3146926",
    "userID": "leichitw"
}
```

#### Response

回傳的data是以樓層(從1樓開始, 失去原始樓層資料)為排列的資料

```json
{
    "status": 200,
    "data":[
        {
            "Num": 1,
            "UserName": "驥哥",
            "UserID": "leichitw",
            "Content": "\n資工大一路過 請問c++學完指標目前也只會用dev c++寫程式 如果要物件導向類是否要換環境？\n"
        },
        {
            "Num": 2,
            "UserName": "驥哥",
            "UserID": "leichitw",
            "Content": "\n  各位大大好,我是今年入學資訊工程一年級的學生,目前學完一學期,c++大概學到指標,教授覺得我能力還不錯就打算培養我大二可以陪教授的團隊接一下簡單的專案,比較著重的是臉部辨識,還有一些深度學習的應用.  所以近期開始在學習python的時候有幾個關於ide的問題想問問大家.一開始我是聽了教授的建議安裝了anaconda這個懶人包,之後用了一段時間她內建的ide(spyder)之後在教授的推薦下我改用了pycharm這個ide,目前有幾個不懂的是:(1):在spyder可以藉由Anaconda prompt 來安裝套件,如果我要在pycharm安裝套件的話,要如何安裝呢?(2):Anaconda裡面不是有很多內建的套件嗎?請問要怎麼跟pycharm做連結呢?(3):所謂的pip是什麼?我知道他能管理與安裝軟體包,但具體不太會用,也不知道要怎麼按出來呢?(4):在現在的市場裏面,資安以及多媒體領域影像處理相關,哪個比較有前途呢?我在想之後選課的時候想要兩種 課程都選下來,業界有人是這兩方面都很精通的嗎?因為教授給了我很多關於深度學習跟臉部辨識的資料,(英文真的苦手)我覺得真的有點不太擅長,所以也沒辦法立刻了解現在的志向是什麼..請版上大大多多指教.(5):所謂的深度學習我的概念還是很模糊,覺得從一般的程式跳入深度的世界是一個門檻,感覺會不知道從何下手呢,深度學習是為了達到目地的一種方法嗎?還是為了加強某些速度的手段呢?目前大概問題就這些,請各位幫我解惑 謝謝各位大大 如果還有疑問我會提出的 謝謝你們!!\n"
        },
        {
            "Num": 3,
            "UserName": "驥哥",
            "UserID": "leichitw",
            "Content": "\n想問此串的前輩們,微積分這門科目,以後會很常用到嗎?還有準備研究所考試的時候佔的比例高嗎?因為我們教授的講法我不太習慣...所以都自己讀偏多,考試也考得不錯.但這最近意識到我觀念不是說很清楚..,感覺有點像只會做題目,如果有人問我觀念,或者某個公式的含意,我都回答不出來.這樣子之後如果研究所複習會不會很吃力呢?要改善這樣的狀況有什麼推薦的用書嗎?還是繼續要學校的參考書就好? 謝謝回答!\n"
        },
        {
            .....
        }
    ]
}
```

- 錯誤時候的回傳

```json
{
    "status": 400,
    "error":  "請確認一下傳入的資料有沒有符合api的格式"
}
```
```json 
{
    "status": 500, 
    "error": "伺服器在處理request的時候發生了錯誤, 請稍後再測試" 
}
```
--- 
### 找到某個User在某個討論串中的所有文章並保留原始樓層
#### 此方法會獲得文章的原始樓層, 不過因為開了大量的goroutine, 所以有機會某些樓層的request會被擋下來

傳遞json參數意義: 
<br>"userID"欄位擺放使用者想要查詢對象的巴哈ID
<br>"baseurl"欄位擺放的是想要查詢的討論串欄位, 值得一提的是不管貼的連結所在的page在哪頁, 都可以藉由api找到整串的資料

#### Request

- Method: **POST**
- Url: ```https://go-gamer.herokuapp.com/FindAllFloor```
- Headers: Content-Type:application/json
- Body:

```json
{
    "baseurl": "https://forum.gamer.com.tw/C.php?page=2&bsn=60076&snA=3146926",
    "userID": "leichitw"
}
```

#### Response

回傳的data是以樓層為排列的資料

```json
{
    "status": 200,
    "data":[
        {
            "Num": 5230,
            "UserName": "驥哥",
            "UserID": "leichitw",
            "Content": "\n如果學校先學c++,還有必要回頭學c嗎?我看網路很多範例都用c寫,沒學過表示有時候看不太懂qq\n"
        },
        {
            "Num": 5208,
            "UserName": "驥哥",
            "UserID": "leichitw",
            "Content": "\n想問此串的前輩們,微積分這門科目,以後會很常用到嗎?還有準備研究所考試的時候佔的比例高嗎?因為我們教授的講法我不太習慣...所以都自己讀偏多,考試也考得不錯.但這最近意識到我觀念不是說很清楚..,感覺有點像只會做題目,如果有人問我觀念,或者某個公式的含意,我都回答不出來.這樣子之後如果研究所複習會不會很吃力呢?要改善這樣的狀況有什麼推薦的用書嗎?還是繼續要學校的參考書就好? 謝謝回答!\n"
        },
        {
            ....
        },
    ]
}
```

- 錯誤時候的回傳

```json
{
    "status": 400,
    "error":  "請確認一下傳入的資料有沒有符合api的格式"
}
```
```json 
{
    "status": 500, 
    "error": "伺服器在處理request的時候發生了錯誤, 請稍後再測試" 
}
```

### 查詢單一用戶資料

以單一用戶的ID獲得其帳號的個人資料


#### Request 

- Method: **GET**
- URL: ```https://go-gamer.herokuapp.com/FindUserInfo?ID={UserID}```
- Example: ```https://go-gamer.herokuapp.com/FindUserInfo?ID=leichitw```
- Headers
- Body 
```


```
#### Response 
- Body 

```json 
{
    "status": 200, 
    "data": {
        "UserID": "leichitw",          
        "UserName": "驥哥",            
        "Title": "只知kuso的小平民",     
        "Lever": 24,                  
        "Race": "人類",                
        "Occupation": "劍士",          
        "Balance": 1285,              
        "GP": 77                       
    }
}
```

- 錯誤時候的回傳

```json
{
    "status": 400,
    "error":  "請確認一下傳入的資料有沒有符合api的格式",
}
```















