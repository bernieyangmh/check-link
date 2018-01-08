package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"regexp"
	"bytes"
	"github.com/golang-collections/go-datastructures/set"
	//"log"
)



const (
	PATTERN_SRC = `src=\"(.*?)\"`
	PATTERN_HERF = `href=\"(.*?)\"`
	PATTERN_HTTP = `http(.*?)`
	PATTERN_LINK = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	PATTERN_SLASH = `\/(.*?)`
)


func main() {
	test_num := 0
	resp, err := http.Get("http://www.qiniu.com")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}


	reHref, _ := regexp.Compile(PATTERN_HERF)
	reSrc, _ := regexp.Compile(PATTERN_SRC)
	reHttp, _ := regexp.Compile(PATTERN_HTTP)
	reLink, _ := regexp.Compile(PATTERN_LINK)
	reSlash, _ := regexp.Compile(PATTERN_SLASH)

	hrefArray := reHref.FindAllStringSubmatch(string(body), 10000)
	srcArray := reSrc.FindAllStringSubmatch(string(body), 10000)

	waitUrlSet := set.New()
	//finishSet := set.New()


	for i := 0; i < len(srcArray); i++ {
		var srcLogBuffer bytes.Buffer
		srcOriginLink := srcArray[i][1]

		srcLogBuffer.WriteString("srcOriginLink is ")
		srcLogBuffer.WriteString(srcOriginLink)

		if reSlash.MatchString(srcOriginLink) != true {
			continue
		}

		if reHttp.MatchString(srcOriginLink) != true {

			var resUrlBuffer bytes.Buffer
			resUrlBuffer.WriteString("http://www.qiniu.com")
			resUrlBuffer.WriteString(srcOriginLink)
			resLink := resUrlBuffer.String()

			srcLogBuffer.WriteString(", resLink is")
			srcLogBuffer.WriteString(resLink)
			//logString := srcLogBuffer.String()
			//log.Println(logString)

			res := reLink.MatchString(resLink)
			if res == true {



				waitUrlSet.Add(resLink)



			}
		}
	}


	for i := 0; i < len(hrefArray); i++ {
		var hrefLogBuffer bytes.Buffer

		hrefOriginLink := hrefArray[i][1]

		hrefLogBuffer.WriteString("hrefOriginLink is ")
		hrefLogBuffer.WriteString(hrefOriginLink)

		if reSlash.MatchString(hrefOriginLink) != true {
			continue
		}

		if reHttp.MatchString(hrefOriginLink) != true {

			var resUrlBuffer bytes.Buffer
			resUrlBuffer.WriteString("http://www.qiniu.com")
			resUrlBuffer.WriteString(hrefOriginLink)
			resLink := resUrlBuffer.String() // 拼接结果
			res := reLink.MatchString(resLink)


			hrefLogBuffer.WriteString(", resLink is")
			hrefLogBuffer.WriteString(resLink)
			//logString := hrefLogBuffer.String()
			//log.Println(logString)

			if res == true {



				waitUrlSet.Add(resLink)
			}

		}
	}
	fmt.Println("now show urls")
	
	fmt.Println(waitUrlSet.Flatten())
	fmt.Println(waitUrlSet.Len())

	waitUrlArray := waitUrlSet.Flatten()
	for i :=0; i< len(waitUrlArray);i++ {
		verifyUrl := waitUrlArray[i]
		fmt.Println(verifyUrl)
		test_num += 1
		fmt.Println(test_num)


	}
	
	
	


}
