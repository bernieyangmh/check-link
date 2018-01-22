package main

import (
	"fmt"
	"log"
	"check-link"
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

	for len(executeChannel) > 0 {
		aimUrl := check_link.GetChannel(executeChannel)
		if aimUrl.CrawlUrl != "close" {
			check_link.IterCrawl(aimUrl, trailMap, executeChannel, &finishArray, &errorArryay)
			fmt.Println(len(executeChannel))
		}
	}

	for i := 0; i < len(finishArray); i++ {
		if finishArray[i].StatusCode != 0 {
			fmt.Println(finishArray[i])
		}
	}

	log.Println("/n url num is %d/n", len(finishArray))

	for i := 0; i < len(errorArryay); i++ {
		if errorArryay[i].StatusCode != 0 {
			fmt.Println(errorArryay[i].CrawlUrl)
			fmt.Println(errorArryay[i].RefUrl)
			fmt.Println(errorArryay[i].StatusCode)
			fmt.Println(errorArryay[i].QueryError)
			fmt.Println("\n")
		}
	}
}
