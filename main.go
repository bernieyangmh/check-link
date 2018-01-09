package main

import (
	"bytes"
	"fmt"
	"net/url"
	"io/ioutil"
	"log"
	"net/http"

	"regexp"
	"github.com/Workiva/go-datastructures/queue"
	"time"
)

const (
	PATTERN_SRC   = `src=\"(.*?)\"`
	PATTERN_HERF  = `href=\"(.*?)\"`
	PATTERN_HTTP  = `http(.*?)`
	PATTERN_LINK  = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	PATTERN_SLASH = `\/(.*?)`
)

func main() {


	ROOT_DOMAIN := `http://www.qiniu.com`
	waitUrlMap := make(map[string]int)
	finishUrlMap := make(map[string]int)

	waitQueue := queue.New(2000)

	waitQueue.Put(ROOT_DOMAIN)

	//todo 从queue里开始
	//respBody, StatusCode, ContentType := Crawling(ROOT_DOMAIN)
	//
	//finishUrlMap[ROOT_DOMAIN] = StatusCode
	//waitUrlMap[ROOT_DOMAIN] = -1

	//if ContentType != "text/html; charset=utf-8" {
	//
	//}
	//
	//hrefArray, srcArray := ExtractBody(respBody)
	//
	//SaToMapQueue(ROOT_DOMAIN, srcArray, waitUrlMap, *waitQueue)
	//
	//HaToMapQueue(ROOT_DOMAIN, hrefArray, waitUrlMap, *waitQueue)

	tmpUrl, err := GetUrlFromQueue(*waitQueue)
	if err != nil {
		if err ==queue.ErrTimeout {
			log.Println( "队列读取完毕", err)
		} else {
			log.Println(err)
		}
	}




	IterCraw(tmpUrl, waitUrlMap, finishUrlMap, *waitQueue)

	for k, v := range waitUrlMap {
		fmt.Println(k, v)
	}
	fmt.Println(len(waitUrlMap))

	for k, v := range finishUrlMap {
		fmt.Println(k, v)
	}
	fmt.Println(len(finishUrlMap))

}




























//拼接域名和路径
func StitchUrl(DomainString string, PathString string) (UString string) {
	var resUrlBuffer bytes.Buffer

	resUrlBuffer.WriteString(DomainString)
	resUrlBuffer.WriteString(PathString)

	UString = resUrlBuffer.String()

	return UString
}

//爬取指定链接，返回响应
func Crawling(UrlString string) (ResponseBodyString string, StatusCode int, ContentType string) {
	resp, err := http.Get(UrlString)
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	return string(body), resp.StatusCode, resp.Header.Get("Content-Type")
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

//从响应body里提取出hre=和src=的相对路径
func ExtractBody(s string) ([][]string, [][]string) {
	hrefArray := ReHrefSubMatch(s)
	srcArray := ReSrcSubMatch(s)
	return hrefArray, srcArray
}

//从指定链接内爬取链接，放入waitUrlMap, waitUrlQueue
func IterCraw(surl string, wum map[string]int, fum map[string]int, wQ queue.Queue) {
	respBody, StatusCode, ContentType := Crawling(surl)

	fum[surl] = StatusCode

	if ContentType == "text/html; charset=utf-8" {
		hrefArray, srcArray := ExtractBody(respBody)

		host := ParseUrlHost(surl)

		SaToMapQueue(host,srcArray, wum, wQ)
		HaToMapQueue(host, hrefArray, wum, wQ)
	}
}

//处理SreArry内的链接写入到waitUrlMap,waitUrlQueue
func SaToMapQueue(d string, sa [][]string, wum map[string]int, wQ queue.Queue) {
	for i := 0; i < len(sa); i++ {
		srcOriginLink := sa[i][1]

		if ReHaveSlash(srcOriginLink) != true {
			continue
		}

		if ReIsHttp(srcOriginLink) != true {

			resLink := StitchUrl(d, srcOriginLink)

			res := ReIsLink(resLink)
			if res == true {

				wum[resLink] = -1
				wQ.Put(resLink)

			}
		}
	}

}

//处理HrefArry内的链接写入到waitUrlMap,waitUrlQueue
func HaToMapQueue(d string, ha [][]string, wum map[string]int, wQ queue.Queue) {
	for i := 0; i < len(ha); i++ {

		hrefOriginLink := ha[i][1]

		if ReHaveSlash(hrefOriginLink) != true {
			continue
		}

		if ReIsHttp(hrefOriginLink) != true {

			resLink := StitchUrl(d, hrefOriginLink)
			res := ReIsLink(resLink)

			if res == true {

				wum[resLink] = -1
				wQ.Put(resLink)

			}
		}
	}
}

func ParseUrlHost(u string) string{
	//解析这个 URL 并确保解析没有出错。
	pU, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return pU.Host

}

func GetUrlFromQueue(wq queue.Queue) (turl string,err error) {
	tmpUrlArray, err := wq.Poll(1, time.Millisecond)
	if (err != nil) || err == queue.ErrTimeout {
		log.Print(err)
	}

	tmpUrl := fmt.Sprintf("%s", tmpUrlArray[0])
	return tmpUrl, err

}