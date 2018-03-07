package check_link

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"
	"strconv"
	"github.com/qiniu/rpc.v1"
	"qbox.us/api/message"
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
		if ReHaveSinlgeSlash(ha) || ReIsLink(ha) || ReHaveMoreSlash(ha) {
			unitCurl.Origin = ha
			if ReHaveSinlgeSlash(ha) {
				unitCurl.CrawlUrl = StitchUrl(cU.Domain, ha)
			}

			//引用为绝对路径,直接赋值
			if ReIsLink(ha) {
				unitCurl.CrawlUrl = ha
			}

			//拿到链接所属链接的协议，与//形式的相对链接合成新链接
			if ReHaveMoreSlash(ha) {
				pu, err := url.Parse(cU.Domain)
				if err != nil {
					log.Println(err)
				}

				var resUrlBuffer bytes.Buffer
				resUrlBuffer.WriteString(pu.Scheme)
				resUrlBuffer.WriteString("://")
				resUrlBuffer.WriteString(ha[2:])

				unitCurl.CrawlUrl = resUrlBuffer.String()
			}

			if cU.CrawlUrl != "" {
				unitCurl.RefUrl = cU.CrawlUrl
			} else {
				log.Print("Nil CrawlUrl!")
			}

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
		if ReHaveSinlgeSlash(ha) || ReIsLink(ha) || ReHaveMoreSlash(ha) {
			//单个/,合成绝对路径
			if ReHaveSinlgeSlash(ha) {
				a[i].CrawlUrl = StitchUrl(cU.Domain, ha)
			}
			//引用为绝对路径,直接赋值
			if ReIsLink(ha) {
				a[i].CrawlUrl = ha
			}

			//拿到链接所属链接的协议，与//形式的相对链接合成新链接
			if ReHaveMoreSlash(ha) {
				pu, err := url.Parse(cU.Domain)
				if err != nil {
					log.Println(err)
				}

				var resUrlBuffer bytes.Buffer
				resUrlBuffer.WriteString(pu.Scheme)
				resUrlBuffer.WriteString("://")
				resUrlBuffer.WriteString(ha[2:])

				a[i].CrawlUrl = resUrlBuffer.String()
			}
			if cU.CrawlUrl != "" {
				a[i].RefUrl = cU.CrawlUrl
			} else {
				log.Print("Nil CrawlUrl!")
			}

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

type ConfigJson struct {
	WhiteLink []string `json:"WhiteLink"`
}

//从配置文件中读取配置项并配置
func ReadJsonConfig(tm map[string]int) {

	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c ConfigJson
	json.Unmarshal(raw, &c)

	for i := 0; i < len(c.WhiteLink); i++ {
		tm[c.WhiteLink[i]] = 1
	}

}

type Transport struct {
	Transport http.RoundTripper
}

func NewTransport(transport http.RoundTripper) *Transport {

	if transport == nil {
		transport = http.DefaultTransport
	}

	return &Transport{Transport: transport}
}

//发送邮件
func SendMail(uid uint32, to []string, subject string, content string) (oid string, err error) {

	cli := &rpc.Client{&http.Client{Transport: NewTransport(nil)}}
	nf := message.NewHandleNotification("https://morse.qiniu.io", cli)

	mailIn := message.SendMailIn{Uid: uid, To: to, Subject: subject, Content: content}
	mailOut, err := nf.SendMail(nil, mailIn)
	oid = mailOut.Oid

	return
}

//爬取,检查,更新
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

	//读取配置文件
	ReadJsonConfig(trailMap)

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
			if err.Error()[:6] == `E11000` {
				err := finishArray[i].Update()
				if err != nil {
					fmt.Println(err)
				}
			}
			log.Println(err)
		}
	}

	log.Println("/n url num is %d/n", len(finishArray))

	for i := 0; i < len(errorArryay); i++ {
		if errorArryay[i].StatusCode != 0 {
			fmt.Println("错误链接		" + errorArryay[i].CrawlUrl)

			fmt.Println("引用链接		" + errorArryay[i].RefUrl)

			fmt.Println(errorArryay[i].StatusCode)

			fmt.Println("链接内容		" + errorArryay[i].Context)

			fmt.Println("访问报错		" + errorArryay[i].QueryError)

			fmt.Println("\n")
		}
	}

	var emailTable bytes.Buffer

	emailTable.WriteString("<html><body><table><tr><td>问题链接<td><td>状态码<td><td>原引用链接<td><td>链接内容<td><td>报错信息<td></tr>")
	for _, ic :=range errorArryay {

		emailTable.WriteString(`<tr><td>`)
		emailTable.WriteString(ic.CrawlUrl)
		emailTable.WriteString(`<td><td>`)
		emailTable.WriteString(strconv.Itoa(ic.StatusCode))
		emailTable.WriteString(`<td><td>`)
		emailTable.WriteString(ic.RefUrl)
		emailTable.WriteString(`<td><td>`)
		emailTable.WriteString(ic.Context)
		emailTable.WriteString(`<td><td>`)
		emailTable.WriteString(ic.QueryError)
		emailTable.WriteString(`<td></tr>`)
	}
	emailTable.WriteString(`</table></body></html>`)

	fmt.Println(emailTable.String())

	to := []string{"yangminghui@qiniu.com"}
	oid, err := SendMail(123, to, "七牛网站扫描结果", emailTable.String())
	if err != nil{
		fmt.Println(oid)
	}


}
