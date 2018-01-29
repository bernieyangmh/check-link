package check_link

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"log"
	"net/url"
	"strings"
	"time"
	"unicode"
)

//输入一个链接，将状态码放进map，能爬取的链接输进管道

//将url放入管道
func PutChannel(cu CUrl, ch chan<- CUrl) {
	ch <- cu
}

//从管道中取出一个url
func GetChannel(ch chan CUrl) CUrl {

	select {
	case u := <-ch:
		return u
	case <-time.After(time.Second * 10):
		close(ch)
		return CUrl{QueryError: "TimeOutClose", StatusCode: -3}
	}
}

//读取数组内的路径，处理为完整url,如果不在Map里放入ch和map
func ReArrayToUrl(cU CUrl, a [][]string, cH chan<- CUrl, tM map[string]int) {

	var unitCurl CUrl
	for i := 0; i < len(a); i++ {
		ha := a[i][1]

		//引用为路径则拼接为完整url
		if ReHaveSlash(ha) || ReIsLink(ha) {
			unitCurl.Origin = ha
			if ReHaveSlash(ha) {
				unitCurl.CrawlUrl = StitchUrl(cU.Domain, ha)
			} else {
				unitCurl.CrawlUrl = ha
			}
			unitCurl.RefUrl = cU.CrawlUrl
			//如果拼接符合url正则且不在Map内的的放入channel和Map
			if ReIsLink(unitCurl.CrawlUrl) && tM[unitCurl.CrawlUrl] == 0 {
				UrlToChMAP(unitCurl, cH, tM)
			} else {
				unitCurl.Origin = ha
				log.Println("ErrorUrl		" + unitCurl.CrawlUrl)
			}
		} else {
			log.Print("ErrorPath			" + ha)
		}
	}
}

//将连接放入channel和map
func UrlToChMAP(cu CUrl, ch chan<- CUrl, tm map[string]int) {

	tm[cu.CrawlUrl] = -1
	log.Println("put			" + cu.CrawlUrl)
	PutChannel(cu, ch)

}

//从body里拿到href和src的相对路径
func ExtractBody(s string) ([]CUrl, [][]string) {

	hrefArray := GetHerfFromHtml(s)
	srcArray := ReSrcSubMatch(s)
	return hrefArray, srcArray
}

//拼接domain和path
func StitchUrl(DomainString string, PathString string) (UString string) {

	var resUrlBuffer bytes.Buffer

	resUrlBuffer.WriteString(DomainString)
	resUrlBuffer.WriteString(PathString)

	UString = resUrlBuffer.String()

	return UString
}

//将Scheme和Host拼接为domain
func StitchDomain(s string, h string) string {

	var resUrlBuffer bytes.Buffer

	resUrlBuffer.WriteString(s)
	resUrlBuffer.WriteString("://")
	resUrlBuffer.WriteString(h)

	domainString := resUrlBuffer.String()

	return domainString

}

//从链接里提取出domain,host
func GetDomainHost(u string) (string, string, error) {

	if !ReIsLink(u) {
		return "", "", errors.New("不符合链接正则")
	}

	pu, err := url.Parse(u)
	if err != nil {
		log.Println(err)
	}

	domainString := StitchDomain(pu.Scheme, pu.Host)
	return domainString, pu.Host, nil

}

//去掉全部空格
func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

//解析body拿到href链接及文本dom内容
func GetHerfFromHtml(s string) []CUrl {
	hrefArray := make([]CUrl, 0)

	node, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Print(err)
	}
	doc := goquery.NewDocumentFromNode(node)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("html a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		linkText := linkTag.Text()
		linkText = SpaceMap(linkText)
		bbb := CUrl{Origin: link, Context: linkText}
		hrefArray = append(hrefArray, bbb)
	})
	return hrefArray
}


func DomArrayToUrl(cU CUrl, a []CUrl, cH chan<- CUrl, tM map[string]int) {

	for i := 0; i < len(a); i++ {
		ha := a[i].Origin

		//引用为路径则拼接为完整url
		if ReHaveSlash(ha) || ReIsLink(ha) {
			a[i].Origin = ha
			if ReHaveSlash(ha) {
				a[i].CrawlUrl = StitchUrl(cU.Domain, ha)
			} else {
				a[i].CrawlUrl = ha
			}
			a[i].RefUrl = cU.CrawlUrl
			//如果拼接符合url正则且不在Map内的的放入channel和Map
			if ReIsLink(a[i].CrawlUrl) && tM[a[i].CrawlUrl] == 0 {
				UrlToChMAP(a[i], cH, tM)
			} else {
				a[i].Origin = ha
				log.Println("ErrorUrl		" + a[i].CrawlUrl)
			}
		} else {
			log.Print("ErrorPath			" + ha)
		}
	}
}