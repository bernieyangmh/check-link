package main

import (
	"bytes"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const (
	PATTERN_SRC   = `src=\"(.*?)\"`
	PATTERN_HERF  = `href=\"(.*?)\"`
	PATTERN_HTTP  = `http(.*?)`
	PATTERN_LINK  = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	PATTERN_SLASH = `\/(.*?)`
)

func main() {

	waitUrlMap := make(map[string]int)
	finishUrlMap := make(map[string]int)


	respBody, StatusCode, ContentType := Crawling("http://www.qiniu.com")

	finishUrlMap["http://www.baidu.com"] = StatusCode

	if ContentType != "text/html; charset=utf-8" {
		
	}
		

	hrefArray, srcArray := ExtractBody(respBody)


	for i := 0; i < len(srcArray); i++ {
		srcOriginLink := srcArray[i][1]

		if ReHaveSlash(srcOriginLink) != true {
			continue
		}

		if ReIsHttp(srcOriginLink) != true {

			resLink := StitchUrl("http://www.qiniu.com", srcOriginLink)

			res := ReIsLink(resLink)
			if res == true {

				waitUrlMap[resLink] = 0

			}
		}
	}

	for i := 0; i < len(hrefArray); i++ {

		hrefOriginLink := hrefArray[i][1]

		if ReHaveSlash(hrefOriginLink) != true {
			continue
		}

		if ReIsHttp(hrefOriginLink) != true {

			resLink := StitchUrl("http://www.qiniu.com", hrefOriginLink)
			res := ReIsLink(resLink)

			if res == true {

				waitUrlMap[resLink] = 0
			}

		}
	}

	for k, v := range waitUrlMap {
		fmt.Println(k, v)
	}
	fmt.Println(len(waitUrlMap))


}

type Spider struct {
	Rooturl string
}

func StitchUrl(DomainString string, PathString string) (UString string) {
	var resUrlBuffer bytes.Buffer

	resUrlBuffer.WriteString(DomainString)
	resUrlBuffer.WriteString(PathString)

	UString = resUrlBuffer.String()

	return UString
}

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

func ReHrefSubMatch(s string) [][]string {
	reHref, _ := regexp.Compile(PATTERN_HERF)
	hrefArray := reHref.FindAllStringSubmatch(s, 10000)

	return hrefArray

}

func ReSrcSubMatch(s string) [][]string {
	reHref, _ := regexp.Compile(PATTERN_SRC)
	srcArray := reHref.FindAllStringSubmatch(s, 10000)

	return srcArray

}

func ReIsHttp(s string) bool {
	reHttp, _ := regexp.Compile(PATTERN_HTTP)
	return reHttp.MatchString(s)

}

func ReIsLink(s string) bool  {
	reLink, _ := regexp.Compile(PATTERN_LINK)
	return reLink.MatchString(s)
}

func ReHaveSlash(s string) bool  {
	reSlash, _ := regexp.Compile(PATTERN_SLASH)
	return reSlash.MatchString(s)

}

func ExtractBody(s string) ([][]string, [][]string) {
	hrefArray := ReHrefSubMatch(s)
	srcArray := ReSrcSubMatch(s)

	return hrefArray, srcArray
}



func IterCraw(surl string, wum map[string]int, fum map[string]int){
	respBody, StatusCode, ContentType := Crawling(surl)

	fum[surl] = StatusCode

	if ContentType == "text/html; charset=utf-8"{
		hrefArray, srcArray := ExtractBody(respBody)

		SaToMap(srcArray, wum)
		HaToMap(hrefArray, wum)
	}
}

func SaToMap(sa [][]string,  wum map[string]int)  {
	for i := 0; i < len(sa); i++ {
		srcOriginLink := sa[i][1]

		if ReHaveSlash(srcOriginLink) != true {
			continue
		}

		if ReIsHttp(srcOriginLink) != true {

			resLink := StitchUrl("http://www.qiniu.com", srcOriginLink)

			res := ReIsLink(resLink)
			if res == true {

				wum[resLink] = 0

			}
		}
	}

}

func HaToMap(ha [][]string,  wum map[string]int)  {
	for i := 0; i < len(ha); i++ {

		hrefOriginLink := ha[i][1]

		if ReHaveSlash(hrefOriginLink) != true {
			continue
		}

		if ReIsHttp(hrefOriginLink) != true {

			resLink := StitchUrl("http://www.qiniu.com", hrefOriginLink)
			res := ReIsLink(resLink)

			if res == true {

				wum[resLink] = 0
			}
		}
	}

}