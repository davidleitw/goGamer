package gamer

type Article struct {
	Index     int      // 文章編號
	SummaryGP int      // 文章總獲得GP數目
	Href      string   // 文章超連結
	Title     string   // 文章標題
	Subbsn    string   // 文章Query編號
	Author    UserInfo // 作者資料
}
