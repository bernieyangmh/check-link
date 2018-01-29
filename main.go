package main

import (
	"check-link"
	"log"
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	check_link.LanuchCrawl()
	check_link.DailyCheck()
}
