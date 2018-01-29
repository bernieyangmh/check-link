package check_link

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

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

	//err = cu.Insert()
	//if err != nil{
	//	log.Println("Insert		" + err.Error())
	//}

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
	resp, err := client.Head(surl)
	if err != nil {
		log.Print(err)
	}

	//链接不允许HEAD方法或直接关闭链接，换用Get
	if resp == nil || resp.StatusCode == 405 {
		log.Println("GetForNoHead		" + surl)
		resp, err = client.Get(surl)
		if err != nil {
			log.Println(err)
		}

	}
	if resp == nil {
		return err.Error(), -2, "error"
	}

	respstatusCode := resp.StatusCode
	respContentType := resp.Header.Get("Content-Type")

	if respContentType == "text/html; charset=utf-8" {
		log.Println("GetForBoby		" + surl)
		resp, err = client.Get(surl)
		if err != nil {
			log.Print(err)
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

func LanuchCrawl() {

	var ROOT_DOMAIN = [2]string{"https://www.qiniu.com", "https://developer.qiniu.com"}

	var executeChannel = make(chan CUrl, 5000)
	var trailMap = make(map[string]int)
	var finishArray = make([]CUrl, 0, 10000)
	var errorArryay = make([]CUrl, 0, 1000)

	firCrawl := CUrl{CrawlUrl: ROOT_DOMAIN[0]}
	secCrawl := CUrl{CrawlUrl: ROOT_DOMAIN[1]}
	//将根域名放入channel
	PutChannel(firCrawl, executeChannel)
	PutChannel(secCrawl, executeChannel)

	for len(executeChannel) > 0 {
		aimUrl := GetChannel(executeChannel)
		if aimUrl.CrawlUrl != "close" {
			IterCrawl(aimUrl, trailMap, executeChannel, &finishArray, &errorArryay)
			fmt.Println(len(executeChannel))
		}
	}

	for i := 0; i < len(finishArray); i++ {
		fmt.Println(finishArray[i])
		err := finishArray[i].Insert()
		if err != nil {
			log.Println(err)

		}
	}

	log.Println("/n url num is %d/n", len(finishArray))

	for i := 0; i < len(errorArryay); i++ {
		if errorArryay[i].StatusCode != 0 {
			fmt.Println(errorArryay[i].CrawlUrl)
			fmt.Println(errorArryay[i].RefUrl)
			fmt.Println(errorArryay[i].StatusCode)
			fmt.Println(errorArryay[i].Context)
			fmt.Println(errorArryay[i].QueryError)
			fmt.Println("\n")
		}
	}
}

func DailyCheck() {
	type Item struct {
		CrawlUrl string `bson:"crawl_url"`
		RefUrl   string `json:"RefUrl" bson:"ref_url"`
	}
	item := Item{}
	items := GetIterUrl()
	for items.Next(&item) {
		url := item.CrawlUrl
		ResponseBodyString, StatusCode, _ := Crawling(url)

		fmt.Println("\n")
		fmt.Println(url)
		fmt.Println(item.RefUrl)
		fmt.Println(StatusCode)
		if StatusCode == -2 {
			fmt.Println(ResponseBodyString)
		}
		fmt.Println("\n\n----------------------------------------------")

	}
}
