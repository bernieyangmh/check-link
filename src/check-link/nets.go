package check_link

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func IterCrawl(cu CUrl, tM map[string]int, cH chan<- CUrl, fA *[]CUrl, eA *[]CUrl) {

	s_domain, _, err := GetDomainHost(cu.CrawlUrl)
	if err != nil {
		log.Println(err)
		cu.QueryError = err.Error()
	}

	respBody, StatusCode, ContentType := Crawling(cu.CrawlUrl)

	//爬过的链接放入trailMap
	tM[cu.CrawlUrl] = StatusCode

	cu.StatusCode = StatusCode
	cu.ContentType = ContentType
	cu.Domain = s_domain

	*fA = append(*fA, cu)
	if cu.StatusCode == -2 {
		cu.QueryError = respBody
	}

	if cu.StatusCode != 200 {
		*eA = append(*eA, cu)
	}

	//如果链接主域名在爬取列表内，Content-Type为html且不在trailMap内，进入读取
	if (ContentType == "text/html; charset=utf-8") && (tM[cu.CrawlUrl] != 0) && ReDomainMatch(cu.CrawlUrl) {
		log.Println("aimUrl		" + cu.CrawlUrl)
		hrefArray, srcArray := ExtractBody(respBody)
		DomArrayToUrl(cu, hrefArray, cH, tM)
		ReArrayToUrl(cu, srcArray, cH, tM)
	}
}

//获取链接的body，状态码，contentType
func Crawling(surl string) (ResponseBodyString string, StatusCode int, ContentType string) {

	var respBody string

	log.Println("Head		" + surl)
	resp, err := http.Head(surl)
	if err != nil {
		log.Print(err)
	}
	if resp == nil {
		return err.Error(), -2, "error"
	}

	//链接不允许HEAD方法或直接关闭链接，换用Get
	if resp == nil || resp.StatusCode == 405 {
		log.Println("GetForNoHead		" + surl)
		resp, err = http.Get(surl)
		if err != nil {
			log.Println(err)
		}

	}
	if resp == nil {
		return err.Error(), -2, "error"
	}

	respstatusCode := resp.StatusCode
	respContentType := resp.Header.Get("Content-Type")

	if 301 == resp.StatusCode || resp.StatusCode == 302 {

		lurl := GetUrlFromLocation(*resp)

		respBody, respstatusCode, respContentType = GetFromRedirectUrl(lurl, 1)
	}

	if respContentType == "text/html; charset=utf-8" {
		log.Println("GetForBoby		" + surl)
		resp, err = http.Get(surl)
		if err != nil {
			log.Print(err)
		}
		if resp == nil {
			return err.Error(), -2, "error"
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		respBody = string(body)
	} else {
		respBody = "nohtml"
	}

	defer resp.Body.Close()

	return respBody, respstatusCode, respContentType
}

//检查重定向是否正确
func GetFromRedirectUrl(lu string, rn int) (string, int, string) {

	resp, err := http.Head(lu)
	if err != nil {
		log.Println(err)
		return err.Error(), -2, "error"
	}
	if resp == nil {
		return err.Error(), -2, "error"
	}

	if resp.StatusCode == 200 {
		return "redict200nohtml", resp.StatusCode, resp.Header.Get("Content-Type")
	}

	if resp.StatusCode == 301 || resp.StatusCode == 302 {
		if rn < 10 {
			rn += 1
			lurl := GetUrlFromLocation(*resp)
			return GetFromRedirectUrl(lurl, rn)
		} else {
			return "redirect too much times", -2, "error"
		}

	}
	return "xxxnohtml", resp.StatusCode, resp.Header.Get("Content-Type")
}

func GetUrlFromLocation(resp http.Response) string {
	var lurl string
	if ReIsLink(resp.Header.Get("Location")) {
		lurl = resp.Header.Get("Location")
	} else {
		lurl = resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + "/" + resp.Header.Get("Location")

		var locationUrlBuffer bytes.Buffer

		locationUrlBuffer.WriteString(resp.Request.URL.Scheme)
		locationUrlBuffer.WriteString("://")
		locationUrlBuffer.WriteString(resp.Request.URL.Host)
		locationUrlBuffer.WriteString("/")
		locationUrlBuffer.WriteString(resp.Header.Get("Location"))

		lurl = locationUrlBuffer.String()
	}
	return lurl
}
