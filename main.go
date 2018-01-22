package main

import (
	"check-link"
	"log"
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var ROOT_DOMAIN = [2]string{"https://www.qiniu.com", "https://developer.qiniu.com"}

	var executeChannel = make(chan check_link.CUrl, 5000)
	var trailMap = make(map[string]int)
	var finishArray = make([]check_link.CUrl, 10000)
	var errorArryay = make([]check_link.CUrl, 500)

	firCrawl := check_link.CUrl{CrawlUrl: ROOT_DOMAIN[0]}
	secCrawl := check_link.CUrl{CrawlUrl: ROOT_DOMAIN[1]}
	//将根域名放入channel
	check_link.PutChannel(firCrawl, executeChannel)
	check_link.PutChannel(secCrawl, executeChannel)


	check_link.LanuchCrawl(executeChannel, trailMap, finishArray, errorArryay)

}
