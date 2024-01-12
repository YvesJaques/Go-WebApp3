package models

type Article struct {
	BlogTitle   string
	BlogArticle string
	UserID      int
	ID          int
}

type ArticleList struct {
	ID      []int
	UserID  []int
	Title   []string
	Content []string
}
