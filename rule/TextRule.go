package rule

import "github.com/PuerkitoBio/goquery"

type TextRule interface {
	NextUrl(document *goquery.Document) (string, bool)
	UrlRule() string
	GetText(document *goquery.Document) string
}
