package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	PATTERN_SRC   = `src=\"(.*?)\"`
	PATTERN_HERF  = `href=\"(.*?)\"`
	PATTERN_HTTP  = `http(.*?)`
	PATTERN_LINK  = `https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	PATTERN_SLASH = `^/(.*?)`
	ALLOW_DOMAIN = `(qiniu.com)|(qiniu.com.cn)`
)

type CUrl struct {
	CrawlUrl    string      `json:"Url" bson:"url"`
	StatusCode	int      	`json:"StatusCode" bson:"status_code"`
	Origin		string      `json:"Origin" bson:"origin"`
	Domain		string      `json:"Domain" bson:"domain"`
	RefUrl		string		`json:"RefUrl" bson:"ref_url"`
	ContentType string		`json:"ContentType" bson:"content_type"`
}

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var ROOT_DOMAIN = [2]string{"https://www.qiniu.com", "https://developer.qiniu.com"}

	var executeChannel = make(chan CUrl, 2000)
	var trailMap = make(map[string]int)
	var finishArray = make([]CUrl, 3000)

	firCrawl := CUrl{CrawlUrl:ROOT_DOMAIN[0]}
	secCrawl := CUrl{CrawlUrl:ROOT_DOMAIN[0]}
	//将根域名放入channel
	PutChannel(firCrawl, executeChannel)
	PutChannel(secCrawl, executeChannel)

	//for aimUrl := range executeChannel {
	//	IterCraw(aimUrl, trailMap, executeChannel)
	//}

	for len(executeChannel) > 0 {
		aimUrl := GetChannel(executeChannel)
		if aimUrl.CrawlUrl != "close" {
			IterCrawl(aimUrl, trailMap, executeChannel, finishArray)
		}
	}

	for i := range finishArray {
		fmt.Println(i)
	}

}

//输入一个链接，将状态码放进map，能爬取的链接输进管道
func IterCrawl(cu CUrl, tM map[string]int, cH chan<- CUrl, fA []CUrl) {


	s_domain, _, err := GetDomainHost(cu.CrawlUrl)
	if err != nil {
		log.Println(err)
	}

	log.Println("Crawl		" + cu.CrawlUrl)
	respBody, StatusCode, ContentType := Crawling(cu.CrawlUrl)

	//爬过的链接放入trailMap
	tM[cu.CrawlUrl] = StatusCode

	cu.StatusCode = StatusCode
	cu.ContentType = ContentType
	cu.Domain = s_domain

	fA = append(fA, cu)



	//如果链接主域名在爬取列表内，Content-Type为html且不在trailMap内，进入读取
	if (ContentType == "text/html; charset=utf-8") && (tM[cu.CrawlUrl] != 0) && ReDomainMatch(cu.CrawlUrl){
		log.Println("aimUrl		" + cu.CrawlUrl)

		hrefArray, srcArray := ExtractBody(respBody)

		ArrayToUrl(cu, hrefArray, cH, tM)
		ArrayToUrl(cu, srcArray, cH, tM)
	}
}

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
		return CUrl{CrawlUrl:"close"}
	}
}

//返回匹配href=的相对路径数组
func ReDomainMatch(s string) bool {
	reDomain, _ := regexp.Compile(ALLOW_DOMAIN)
	isAllow := reDomain.MatchString(s)
	return isAllow

}

//返回匹配href=的相对路径数组
func ReHrefSubMatch(s string) [][]string {
	reHref, _ := regexp.Compile(PATTERN_HERF)
	hrefArray := reHref.FindAllStringSubmatch(s, 10000)

	return hrefArray

}

//返回匹配src=的相对路径数组
func ReSrcSubMatch(s string) [][]string {
	reHref, _ := regexp.Compile(PATTERN_SRC)
	srcArray := reHref.FindAllStringSubmatch(s, 10000)

	return srcArray

}

func ReLinkSubMatch(s string) [][]string {
	reLink, _ := regexp.Compile(PATTERN_LINK)
	srcArray := reLink.FindAllStringSubmatch(s, 10000)

	return srcArray

}

//re匹配http链接
func ReIsHttp(s string) bool {
	reHttp, _ := regexp.Compile(PATTERN_HTTP)
	return reHttp.MatchString(s)

}

//re匹配链接
func ReIsLink(s string) bool {
	reLink, _ := regexp.Compile(PATTERN_LINK)
	return reLink.MatchString(s)
}

//re匹配/
func ReHaveSlash(s string) bool {
	reSlash, _ := regexp.Compile(PATTERN_SLASH)
	return reSlash.MatchString(s)

}

//读取数组内的路径，处理为完整url,如果不在Map里放入ch和map
func ArrayToUrl(cU CUrl, a [][]string, cH chan<- CUrl, tM map[string]int) {
	var unitCurl CUrl
	for i := 0; i < len(a); i++ {
		ha := a[i][1]

		//引用为路径则拼接为完整url
		if ReHaveSlash(ha) {
			unitCurl.Origin = ha
			unitCurl.CrawlUrl = StitchUrl(cU.Domain, ha)
			unitCurl.RefUrl = cU.CrawlUrl
			//如果拼接符合url正则且不在Map内的的放入channel和Map todo "http://url/a.jpg"
			if ReIsLink(unitCurl.CrawlUrl) && tM[unitCurl.CrawlUrl] == 0{
				UrlToChMAP(unitCurl, cH, tM)
			}else {
				unitCurl.Origin = ha
				log.Println("ErrorUrl		"+unitCurl.CrawlUrl)
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
func ExtractBody(s string) ([][]string, [][]string) {
	hrefArray := ReHrefSubMatch(s)
	srcArray := ReSrcSubMatch(s)
	return hrefArray, srcArray
}

//获取链接的body，状态码，contentType
func Crawling(surl string) (ResponseBodyString string, StatusCode int, ContentType string) {
	resp, err := http.Get(surl)
	if err != nil {
		log.Print(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	respstatusCode := resp.StatusCode
	respContentType := resp.Header.Get("Content-Type")
	respBody := string(body)

	defer resp.Body.Close()

	return respBody, respstatusCode, respContentType
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
