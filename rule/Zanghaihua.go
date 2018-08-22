package rule

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type Zhh struct {
	Url         string
	NextPageUrl string
	FilePath    string
}

func InitZhh() *Zhh {
	bj := new(Zhh)
	bj.Url = bj.UrlRule()
	return bj
}

func (ahh *Zhh) UrlRule() string {
	return "http://www.zanghaihua.org/645.html"
}

func (ahh *Zhh) NextUrl(document *goquery.Document) (string, bool) {
	nextUrl, succ := document.Find(".linkbtn a:last-child").Attr("href")
	if succ {
		return nextUrl, true
	} else {
		return "", false
	}
}
func (ahh *Zhh) GetText(document *goquery.Document) string {
	var text string
	document.Find("#BookText").Each(func(i int, selection *goquery.Selection) {
		text = selection.Text()
		strings.TrimSpace(text)
		strings.Replace(text, "\n", "", -1)
	})
	return text
}
