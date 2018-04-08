package check_link

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
	DisableKeepAlives:   true,
}

var client = http.Client{
	Timeout:   time.Duration(10 * time.Second),
	Transport: netTransport,
}

func IterCrawl(cu CUrl, tM map[string]int, cH chan<- CUrl, fA *[]CUrl, eA *[]CUrl, rdl []string) {
	log.Println("test")
	s_domain, _, err := GetDomainHost(cu.CrawlUrl)
	if err != nil {
		log.Println(err)
		cu.QueryError = err.Error()
	}

	respBody, StatusCode, ContentType := Crawling(cu.CrawlUrl)

	//爬过的链接放入trailMap,避免重复检查
	tM[cu.CrawlUrl] = StatusCode

	cu.StatusCode = StatusCode
	cu.ContentType = ContentType
	cu.Domain = s_domain

	*fA = append(*fA, cu)

	//如果访问异常,QueryError为相关响应记录下来
	if cu.StatusCode == -2 {
		cu.QueryError = respBody
	}

	//错误链接放入errorArryay
	if cu.StatusCode != 200 {
		*eA = append(*eA, cu)
	}

	//链接主域名只有在爬取列表内，Content-Type为html且不在trailMap内，才会进入读取
	if (strings.Contains("text/html; charset=utf-8", ContentType)) && (tM[cu.CrawlUrl] != 0) && stringInStringList(cu.CrawlUrl, rdl) {
		log.Println("aimUrl		" + cu.CrawlUrl)
		hrefArray, srcArray := ExtractBody(respBody)
		DomArrayToUrl(cu, hrefArray, cH, tM)
		ReArrayToUrl(cu, srcArray, cH, tM)
	} else {
		log.Println("Not need read body		" + cu.CrawlUrl)
	}
}

//获取链接的body，状态码，contentType
func Crawling(surl string) (ResponseBodyString string, StatusCode int, ContentType string) {
	log.Println("test")
	var respBody string
	var resp *http.Response
	var err error
	var respstatusCode int
	var respContentType string
	var body []byte
	//Todo switch
	log.Println("Head		" + surl)
	resp, err = client.Head(surl)
	if err != nil {
		log.Print(err)
	}

	//链接不允许HEAD方法或直接关闭链接，换用Get
	if resp == nil || resp.StatusCode == 405 {
		log.Println("GetForNoHead		" + surl)
		resp, err = client.Get(surl)
		if err != nil {
			log.Println(err)
		} else {
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
		}

	}
	if resp == nil {
		return err.Error(), -2, "error"
	}

	respstatusCode = resp.StatusCode
	respContentType = resp.Header.Get("Content-Type")

	//如果3xx跳转，检查跳转是否正常
	if 301 == resp.StatusCode || resp.StatusCode == 302 {

		lurl := GetUrlFromLocation(*resp)

		return GetFromRedirectUrl(lurl, 1)
	}

	//如果响应类型为html文件，获取其body
	if strings.Contains("text/html; charset=utf-8", ContentType) {

		if len(body) == 0 {
			log.Println("GetForBoby		" + surl)
			resp, err = client.Get(surl)
			if err != nil {
				log.Println(err)
			} else {
				if resp == nil {
					return err.Error(), -2, "error"
				} else {
					log.Println("test")
					body, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Println(err)
					}
					defer resp.Body.Close()
				}
			}
		}
		respBody = string(body)
	} else {
		respBody = "NoHtml"
	}

	return respBody, respstatusCode, respContentType
}

//检查重定向是否正确
func GetFromRedirectUrl(lu string, rn int) (string, int, string) {
	log.Println("test")
	resp, err := http.Head(lu)
	if err != nil {
		log.Println(err)
		return err.Error(), -2, "error"
	}
	if resp == nil {
		return err.Error(), -2, "error"
	}

	if resp.StatusCode > 200 && resp.StatusCode < 299 {
		return "CorrectlyRedict", resp.StatusCode, resp.Header.Get("Content-Type")
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
	log.Println("test")
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
