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
)

const (
	PATTERN_SRC   = `src=\"(.*?)\"`
	PATTERN_HERF  = `href=\"(.*?)\"`
	PATTERN_HTTP  = `http(.*?)`
	PATTERN_LINK  = `https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	PATTERN_SLASH = `\/(.*?)`
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var ROOT_DOMAIN = [1]string{"http://www.qiniu.com"}

	var executeChannel = make(chan string, 2000)
	var trailMap = make(map[string]int)

	//将根域名放入channel
	PutChannel(ROOT_DOMAIN[0], executeChannel)

	for aimUrl := range executeChannel {
		IterCraw(aimUrl, trailMap, executeChannel)
	}

	for k, v := range trailMap {
		fmt.Println(k, v)
	}

}

//输入一个链接，将状态码放进map，能爬取的链接输进管道
func IterCraw(surl string, tM map[string]int, cH chan<- string) {

	s_domain, _, err := GetDomainHost(surl)
	if err != nil {
		log.Println(err)
	}

	respBody, StatusCode, ContentType := Crawling(surl)

	//爬过的链接放入trailMap
	if tM[surl] == 0 {
		tM[surl] = StatusCode

	}

	//如果链接的Content-Type为html，进入读取且不在trailMap内
	if (ContentType == "text/html; charset=utf-8") && (tM[surl] != 0) {
		log.Println("aimUrl		" + surl)

		hrefArray, srcArray := ExtractBody(respBody)

		ArrayToUrl(s_domain, hrefArray, cH, tM)
		ArrayToUrl(s_domain, srcArray, cH, tM)
	}
}

//将url放入管道
func PutChannel(u string, ch chan<- string) {
	ch <- u
}

//从管道中取出一个url
func GetChannel(ch <-chan string) string {
	url := <-ch
	return url
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
func ArrayToUrl(d string, a [][]string, cH chan<- string, tM map[string]int) {
	var unitUrl string
	for i := 0; i < len(a); i++ {
		ha := a[i][1]

		//引用为路径则拼接为完整url
		if !ReIsLink(ha) && ReHaveSlash(ha) {
			unitUrl = StitchUrl(d, ha)
			//如果拼接符合url正则且不在Map内的的放入channel和Map
			if ReIsLink(unitUrl) && tM[unitUrl] == 0{
				UrlToChMAP(unitUrl, cH, tM)
			}else {
				unitUrl = ha
				log.Println("ErrorUrl		"+unitUrl)
			}
		}
	}
}

//将连接放入channel和map
func UrlToChMAP(d string, ch chan<- string, tm map[string]int) {
	tm[d] = -1
	log.Println("put			" + d)
	PutChannel(d, ch)

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
