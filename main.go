package main

import (
	"check-link"
	"log"
	"os"
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var rootLinkArray []string
	var logPath string
	var resultPath string



	for i:=0;i<len(os.Args);i++{
		switch os.Args[i] {

		case "-o":
			for os.Args[i+1][0] != '-'{
				rootLinkArray = append(rootLinkArray, os.Args[i+1])
				if i < len(os.Args)-2{
					i++;
				} else {
					break
				}
			}
		case "-l":
			logPath = os.Args[i+1]


		case "-r":
			resultPath = os.Args[i+1]

		}
	}

	check_link.LanuchCrawl(rootLinkArray, logPath, resultPath)
}
